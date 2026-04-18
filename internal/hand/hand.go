// Package hand evaluates the value of a hand of Flesh and Blood cards played in isolation.
package hand

import (
	"fmt"
	"sort"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/hero"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/weapon"
)

// Role is what a card does on a given turn cycle.
type Role uint8

const (
	Pitch Role = iota
	Attack
	Defend
	Held
	// Arsenal is assigned to exactly one hand card at most, representing "placed into the
	// arsenal at end of turn". The card contributes nothing to this turn's Value (true value
	// accrues when it's played from arsenal on a future turn), and it carries across the turn
	// boundary as Play.ArsenalCard. Never produced by the recurse itself — Best upgrades one
	// Held card to Arsenal post-hoc when the arsenal slot is empty.
	Arsenal
)

// CardAssignment is a single card + the role it took this turn. Hand cards produce one per
// card; the previous-turn arsenal card also contributes one CardAssignment with FromArsenal set,
// so the whole turn fits in one slice rather than separate hand + role + arsenal structures.
// Contribution is the per-card credit toward TurnSummary.Value (damage dealt, block share, or
// pitch resource depending on Role), populated by fillContributions once the winner is picked.
type CardAssignment struct {
	Card         card.Card
	Role         Role
	FromArsenal  bool
	Contribution float64
}

// TurnSummary is the result of running Best on a hand: the winning "best line" of card-role
// assignments plus aggregate metadata about the turn.
type TurnSummary struct {
	// BestLine is the winning partition as a sequence of CardAssignment values. Hand cards come
	// first in canonical (post-sort) order; the previous-turn arsenal card, if any, follows as
	// the last entry with FromArsenal=true. Never mutate — memoized results alias this slice.
	BestLine []CardAssignment
	// AttackChain is the winning chain in play order — attack-role cards from BestLine mixed
	// with any swung weapons at the positions the solver chose to swing them. Each entry
	// carries its Play-time damage (including hero-trigger damage for cards) so callers can
	// attribute contribution even to weapons (which have no BestLine entry). Swung weapons can
	// be recovered by filtering AttackChain for AttackChainEntry.Card of type weapon.Weapon.
	// Empty when no attacks were played.
	AttackChain []AttackChainEntry
	// Value is the turn's total score (damage dealt + damage prevented). Equals the sum of
	// BestLine[].Contribution plus any weapon-swing damage (weapons aren't in BestLine).
	Value int
	// LeftoverRunechants is the number of Runechant tokens in play at the end of the chosen
	// chain, which carry into the next turn's Best call. For partitions with no attacks, this
	// equals the carryover the caller passed in.
	LeftoverRunechants int
	// ArsenalCard is the card occupying the arsenal slot at the end of this turn — either a new
	// hand card just arsenaled (role=Arsenal) or a previous-turn arsenal card that stayed. Nil
	// when the slot is empty. The caller feeds this back as the next turn's arsenalCardIn.
	ArsenalCard card.Card
}

// AttackChainEntry is a single played attack — either a card with role=Attack or a swung
// weapon — carrying the damage it contributed at the moment it resolved in the winning chain.
// Damage is the Play() return; TriggerDamage is the hero's OnCardPlayed contribution (e.g.
// Viserai creating a Runechant) so callers can surface hero attribution on its own line.
// For BestLine Attack entries Damage + TriggerDamage mirrors CardAssignment.Contribution;
// weapons live only here.
type AttackChainEntry struct {
	Card          card.Card
	Damage        float64
	TriggerDamage float64
}

// ArsenalIn returns the assignment for the card that started the turn in the arsenal, if any.
// Handy for callers (display, per-card stats) that want to treat the arsenal-in card differently
// from hand cards without scanning BestLine themselves.
func (t TurnSummary) ArsenalIn() (CardAssignment, bool) {
	for _, a := range t.BestLine {
		if a.FromArsenal {
			return a, true
		}
	}
	return CardAssignment{}, false
}


// String returns a human-readable role name ("PITCH", "ATTACK", "DEFEND", "HELD", "ARSENAL").
func (r Role) String() string {
	switch r {
	case Pitch:
		return "PITCH"
	case Attack:
		return "ATTACK"
	case Defend:
		return "DEFEND"
	case Held:
		return "HELD"
	case Arsenal:
		return "ARSENAL"
	}
	return "UNKNOWN"
}

// formatContribution renders a contribution/damage value for the best-turn printout. Integer
// values (the common case: pitch amounts, attack power, hero-trigger damage) render without a
// decimal; fractional values (proportional defense share when multiple blockers split an
// incoming attack) show one decimal place.
func formatContribution(v float64) string {
	if v == float64(int(v)) {
		return fmt.Sprintf("%d", int(v))
	}
	return fmt.Sprintf("%.1f", v)
}

// FormatBestLine pairs each card in BestLine with its assigned role for debug output, e.g.
// "Hocus Pocus (Blue): PITCH, Runic Reaping (Red): ATTACK". Cards that came in from arsenal
// get a " (from arsenal)" tag on their name so the reader can see why that card isn't in the
// hand list the optimiser reports alongside. This is the compact one-line form — use
// FormatBestTurn for the chronological play-order presentation.
func FormatBestLine(line []CardAssignment) string {
	parts := make([]string, len(line))
	for i, a := range line {
		name := a.Card.Name()
		if a.FromArsenal {
			name += " (from arsenal)"
		}
		parts[i] = name + ": " + a.Role.String()
	}
	return strings.Join(parts, ", ")
}

// FormatBestTurn renders a TurnSummary as a numbered play-order list, one card per line. The
// order mirrors the actual sequence of a FaB turn:
//
//  1. Defense-phase pitches (paying for Defense Reactions)
//  2. Plain blocks
//  3. Defense Reactions
//  4. Attack-phase pitches (paying for this turn's played cards)
//  5. Attack chain — played cards and swung weapons in the order the solver picked
//
// Cards that ended up Held or Arsenal-bound (didn't get played or pitched this turn) are
// summarized on trailing lines so the reader can see what's carrying over.
//
// Pitch-phase assignment is computed here by a simple greedy: smallest pitches first fund the
// defense pool until drCost is covered, the rest fund attack. The solver validated that some
// legal split exists; this chooses one deterministically for display.
func FormatBestTurn(t TurnSummary) string {
	// Partition BestLine entries into role buckets; pitch cards come out as a single pool which
	// then gets split by phase below.
	var pitched, plainBlocks, defenseReactions, held, arsenal []CardAssignment
	var attackCost, drCost int
	for _, a := range t.BestLine {
		switch a.Role {
		case Pitch:
			pitched = append(pitched, a)
		case Attack:
			attackCost += a.Card.Cost()
		case Defend:
			if a.Card.Types().Has(card.TypeDefenseReaction) {
				drCost += a.Card.Cost()
				defenseReactions = append(defenseReactions, a)
			} else {
				plainBlocks = append(plainBlocks, a)
			}
		case Held:
			held = append(held, a)
		case Arsenal:
			arsenal = append(arsenal, a)
		}
	}

	// Greedy split: sort pitches ascending and pour into the defense bucket until drCost is
	// covered; the rest is attack-phase. Stable w.r.t. input order for ties (sort.SliceStable).
	sorted := append([]CardAssignment(nil), pitched...)
	sort.SliceStable(sorted, func(i, j int) bool { return sorted[i].Card.Pitch() < sorted[j].Card.Pitch() })
	var defensePitches, attackPitches []CardAssignment
	covered := 0
	for _, a := range sorted {
		if covered < drCost {
			defensePitches = append(defensePitches, a)
			covered += a.Card.Pitch()
		} else {
			attackPitches = append(attackPitches, a)
		}
	}

	var lines []string
	step := 0
	nextStep := func() int { step++; return step }
	// Pitch lines don't show a damage contribution — pitches generate resource (already factored
	// into the turn's affordability) rather than damage or prevention. Showing "+3" next to a
	// pitch would double-count against the turn's Value.
	appendPitch := func(a CardAssignment, roleLabel string) {
		name := a.Card.Name()
		if a.FromArsenal {
			name += " (from arsenal)"
		}
		lines = append(lines, fmt.Sprintf("  %d. %s: %s", nextStep(), name, roleLabel))
	}
	// Defense lines DO contribute damage-prevention (block share + DR Play return) so their
	// tag is shown.
	appendDefense := func(a CardAssignment, roleLabel string) {
		name := a.Card.Name()
		if a.FromArsenal {
			name += " (from arsenal)"
		}
		lines = append(lines, fmt.Sprintf("  %d. %s: %s (+%s prevented)", nextStep(), name, roleLabel, formatContribution(a.Contribution)))
	}

	for _, a := range defensePitches {
		appendPitch(a, "PITCH (opponent's turn)")
	}
	for _, a := range plainBlocks {
		appendDefense(a, "BLOCK")
	}
	for _, a := range defenseReactions {
		appendDefense(a, "DEFENSE REACTION")
	}
	for _, a := range attackPitches {
		appendPitch(a, "PITCH (my turn)")
	}
	// Attack chain: iterate AttackChain for real play order, cross-referencing BestLine by ID to
	// mark arsenal-played cards. Weapons have no BestLine entry, so they render as plain names.
	// Each entry prints its Play damage; if the hero's OnCardPlayed fired on this entry (e.g.
	// Viserai creating a Runechant on a 2nd+ chain link), append " +M hero trigger" so the
	// attribution is visible rather than silently folded into the card's number.
	used := make([]bool, len(t.BestLine))
	appendAttack := func(label, cardName string, e AttackChainEntry) {
		line := fmt.Sprintf("  %d. %s: %s (+%s)", nextStep(), cardName, label, formatContribution(e.Damage))
		if e.TriggerDamage > 0 {
			line += fmt.Sprintf(" (+%s hero trigger)", formatContribution(e.TriggerDamage))
		}
		lines = append(lines, line)
	}
	for _, e := range t.AttackChain {
		if _, isWeapon := e.Card.(weapon.Weapon); isWeapon {
			appendAttack("WEAPON ATTACK", e.Card.Name(), e)
			continue
		}
		// Find first unused BestLine entry with matching ID so we can detect FromArsenal.
		tag := ""
		for i := range t.BestLine {
			if used[i] || t.BestLine[i].Role != Attack || t.BestLine[i].Card.ID() != e.Card.ID() {
				continue
			}
			if t.BestLine[i].FromArsenal {
				tag = " (from arsenal)"
			}
			used[i] = true
			break
		}
		appendAttack("ATTACK", e.Card.Name()+tag, e)
	}

	// Held / Arsenal footer. These didn't get "played", so they're outside the numbered
	// sequence. Keep them visible so the reader knows the whole turn disposition.
	var footers []string
	for _, a := range held {
		footers = append(footers, fmt.Sprintf("  (held: %s)", a.Card.Name()))
	}
	for _, a := range arsenal {
		label := a.Card.Name()
		if a.FromArsenal {
			label += " (stayed)"
		} else {
			label += " (new)"
		}
		footers = append(footers, fmt.Sprintf("  (arsenal: %s)", label))
	}
	lines = append(lines, footers...)
	return strings.Join(lines, "\n")
}

// Best returns the optimal TurnSummary for the given hand against an opponent that will attack
// for incomingDamage on their next turn. Any equipped weapons may also be swung for their Cost
// if resources allow.
//
// Cards are partitioned into five roles:
//   - Pitch: contributes its Pitch value as resources paying for a played card on this turn or a
//     Defense Reaction on the opponent's turn.
//   - Attack: consumes Cost resources on our turn; the attack is resolved by calling Card.Play in
//     some order the optimizer chooses. Effects on TurnState carry forward to later attacks in
//     the same sequence.
//   - Defend: contributes Defense to damage prevented (capped at incomingDamage; excess block is
//     wasted). Plain blocking is free; Defense-Reaction-typed cards must have their Cost paid to
//     resolve and also contribute any Play() damage.
//   - Held: the card stays in hand to next turn. Contributes nothing this turn.
//   - Arsenal: the card moves into the arsenal slot at end of turn (or, for an arsenal-in card,
//     stays there). Contributes nothing this turn.
//
// Pitch resources are split across two phases because resources don't carry between turns:
// attack-phase pitches pay for our played attack actions and non-attack actions, defence-phase
// pitches pay for our Defense Reactions on the opponent's turn. A card cannot be pitched unless
// at that moment there's an unpaid card on the stack in the matching phase — so a hand with no
// plays in a phase must have no Pitch-role cards tagged to that phase. Any remaining cards that
// can't be legally pitched become Held.
//
// Results are memoized on (hero name, sorted weapon names, sorted card IDs, incomingDamage,
// runechantCarryover, arsenal-in ID) so repeated evaluations of the same hand across shuffles
// short-circuit. The hand is sorted in place into canonical order first — BestLine's hand
// entries align with that post-sort order. Every card in the hand must be registered in package
// cards; Best panics otherwise.
//
// runechantCarryover is the number of Runechant tokens carrying in from the previous turn. The
// returned TurnSummary.LeftoverRunechants is the count remaining at end of the chosen chain,
// which the caller should feed back as the next turn's carryover.
//
// arsenalCardIn is the card sitting in the arsenal slot at the start of this turn (nil if the
// slot is empty). The partition enumerator pulls it into the turn as an extra CardAssignment
// with restricted role options — Arsenal (stay in the slot), Attack (any non-DR card), or
// Defend (only Defense Reactions). Never Pitch or Held. A hand card may also take the Arsenal
// role so long as at most one card in BestLine ends up there.
func Best(hero hero.Hero, weapons []weapon.Weapon, hand []card.Card, incomingDamage int, deck []card.Card, runechantCarryover int, arsenalCardIn card.Card) TurnSummary {
	return sharedEvaluator.Best(hero, weapons, hand, incomingDamage, deck, runechantCarryover, arsenalCardIn)
}

// Best is the method form of the package-level Best: same semantics, but uses this Evaluator's
// memo and bufs so concurrent goroutines can each hold their own state.
func (e *Evaluator) Best(hero hero.Hero, weapons []weapon.Weapon, hand []card.Card, incomingDamage int, deck []card.Card, runechantCarryover int, arsenalCardIn card.Card) TurnSummary {
	// Fetch IDs into a fixed-size stack array to avoid a per-call slice allocation. Hand size is
	// capped at 8 (matches memoKey.cardIDs); larger hands panic out of the inner loops elsewhere.
	n := len(hand)
	var ids [8]card.ID
	memoable := true
	for i, c := range hand {
		ids[i] = c.ID()
		if _, ok := c.(card.NoMemo); ok {
			memoable = false
		}
	}
	if arsenalCardIn != nil {
		if _, ok := arsenalCardIn.(card.NoMemo); ok {
			memoable = false
		}
	}

	sortHandByID(hand, ids[:], n)

	key := makeMemoKey(hero, weapons, &ids, n, incomingDamage, runechantCarryover, arsenalCardIn)
	if memoable {
		memoMu.RLock()
		cached, hit := memo[key]
		memoMu.RUnlock()
		if hit {
			return cached
		}
	}
	result := e.bestUncached(hero, weapons, hand, incomingDamage, deck, runechantCarryover, arsenalCardIn)
	if memoable {
		memoMu.Lock()
		memo[key] = result
		memoMu.Unlock()
	}
	return result
}

// sortHandByID sorts the first n entries of `hand` and `ids` in parallel by ascending id, in
// place. Insertion sort — for n ≤ 8 this is faster than sort.Sort and avoids boxing the slices
// through sort.Interface on every call. Canonicalizing the hand order is what lets the memo key
// collapse permutations of the same cards onto a single entry.
func sortHandByID(hand []card.Card, ids []card.ID, n int) {
	for i := 1; i < n; i++ {
		for j := i; j > 0 && ids[j-1] > ids[j]; j-- {
			ids[j-1], ids[j] = ids[j], ids[j-1]
			hand[j-1], hand[j] = hand[j], hand[j-1]
		}
	}
}

// memoKey is a comparable struct used as the map key for the shared memo. Hand size is capped at
// 8 cards. Hero ID + weapon IDs are in the key so concurrent Evaluators with different (hero,
// weapons) tuples share the memo safely without a scope-wipe step: distinct scopes just produce
// distinct keys. Using hero.ID (uint16) instead of the display name keeps the entire key a fixed-
// size integer struct — no string hashing per lookup.
type memoKey struct {
	heroID             hero.ID
	weaponIDs          [2]card.ID
	cardIDs            [8]card.ID
	cardCount          uint8
	incoming           int
	runechantCarryover int
	// arsenalInID is card.Invalid when the slot is empty, otherwise the ID of the card in the
	// arsenal at the start of the turn — different arsenal-ins give distinct cache entries.
	arsenalInID card.ID
}

// Evaluator owns the per-goroutine mutable state hand.Best threads through a single evaluation:
// a scratch-buffer cache keyed by (handSize, weapons). The memo cache is shared across all
// Evaluators (see `memo` below) so every worker benefits from previously-cached hands — only the
// scratch buffers (which are mutated during a call) must be per-goroutine.
//
// Long-lived Evaluators avoid reallocating ~20 scratch slices on every call. Iterate mode keeps
// one Evaluator per worker goroutine and reuses it across every mutation that worker screens.
type Evaluator struct {
	// bufs holds the pre-allocated scratch slices used by bestUncached / bestAttackWithWeapons /
	// bestSequence. Keyed by (handSize, weaponCount, weaponIDs); recreated when any differ so the
	// scratch sizing stays correct.
	bufs            *attackBufs
	bufsHandSize    int
	bufsWeaponIDs   [2]card.ID
	bufsWeaponCount int
	bufsValid       bool
}

// NewEvaluator returns a fresh Evaluator with empty scratch buffers. Reusing one Evaluator across
// many Best calls amortises the scratch allocation; holding separate Evaluators per goroutine
// keeps them safe for concurrent use while still sharing the global memo.
func NewEvaluator() *Evaluator {
	return &Evaluator{}
}

// sharedEvaluator backs the package-level Best function — single-threaded callers don't need to
// construct their own Evaluator. Parallel callers (iterate's worker pool) create one per
// goroutine for bufs but all of them share the global memo below.
var sharedEvaluator = NewEvaluator()

// memo caches canonical-order TurnSummary results keyed by memoKey. Shared across all Evaluators
// and goroutines; protected by memoMu. Hero name + weapon IDs live in memoKey so distinct (hero,
// weapons) scopes coexist in the same map without a wipe step.
var (
	memo   = map[memoKey]TurnSummary{}
	memoMu sync.RWMutex
)

// makeMemoKey builds a comparable memo key. The hand must already be sorted by card ID.
// sortedIDs is passed as a pointer to the caller's [8]card.ID stack array to avoid a slice-header
// escape. weapon IDs are sorted numerically into the two fixed slots so loadouts passed in
// different orders still hash to the same key.
func makeMemoKey(hero hero.Hero, weapons []weapon.Weapon, sortedIDs *[8]card.ID, n int, incoming int, runechantCarryover int, arsenalCardIn card.Card) memoKey {
	k := memoKey{
		heroID:             hero.ID(),
		incoming:           incoming,
		runechantCarryover: runechantCarryover,
		cardCount:          uint8(n),
		cardIDs:            *sortedIDs,
	}
	if arsenalCardIn != nil {
		k.arsenalInID = arsenalCardIn.ID()
	}
	switch len(weapons) {
	case 1:
		k.weaponIDs[0] = weapons[0].ID()
	case 2:
		a, b := weapons[0].ID(), weapons[1].ID()
		if a > b {
			a, b = b, a
		}
		k.weaponIDs[0], k.weaponIDs[1] = a, b
	}
	return k
}

// attackerMeta caches the scalar card attributes playSequence reads on every permutation. With
// this hoisted to a per-attacker lookup, the hot inner loop doesn't dispatch Types / Cost /
// GoAgain through the Card interface on every iteration — the permutation enumerator evaluates
// up to N! orderings of the same N cards, so one meta build amortises across all of them.
type attackerMeta struct {
	types            card.TypeSet
	cost             int
	printedCost      int // PrintedCost for DiscountPerRunechant cards; 0 otherwise
	isDiscount       bool
	baseGoAgain      bool // Card.GoAgain() — combined with PlayedCard.GrantedGoAgain at read time
	isAttackOrWeapon bool
}

// cardMetaCache and cardMetaReady are the shared, read-only-after-init card metadata tables.
// Populated lazily via cardMetaSlowPath on first encounter and then read freely from multiple
// goroutines without synchronisation. Sized for the entire uint16 ID space so lookups are plain
// bounds-checked reads — ~2 MB of immutable memory, trivially cheap given the target machine.
const cardMetaCacheSize = 1 << 16

var (
	cardMetaCache [cardMetaCacheSize]attackerMeta
	cardMetaReady [cardMetaCacheSize]uint32 // written once (atomically) per ID; 0 = unready, 1 = ready
	cardMetaMu    sync.Mutex
)

// attackerMetaFor returns cached metadata for c, populating the cache on first encounter. Safe to
// call from multiple goroutines: the first writer per ID takes the mutex, subsequent readers see
// the ready flag set with a release barrier and read the (now immutable) meta entry directly.
func attackerMetaFor(c card.Card) attackerMeta {
	id := c.ID()
	if atomic.LoadUint32(&cardMetaReady[id]) == 1 {
		return cardMetaCache[id]
	}
	return cardMetaSlowPath(c, id)
}

// cardMetaSlowPath populates the cache entry for the given card under cardMetaMu. Returns the
// computed meta — no need for a second load of the cache since this goroutine just wrote it.
func cardMetaSlowPath(c card.Card, id card.ID) attackerMeta {
	cardMetaMu.Lock()
	defer cardMetaMu.Unlock()
	// Re-check under lock; another goroutine may have populated between the atomic load above and
	// our Lock below.
	if atomic.LoadUint32(&cardMetaReady[id]) == 1 {
		return cardMetaCache[id]
	}
	t := c.Types()
	m := attackerMeta{
		types:            t,
		cost:             c.Cost(),
		baseGoAgain:      c.GoAgain(),
		isAttackOrWeapon: t.Has(card.TypeAttack) || t.Has(card.TypeWeapon),
	}
	if d, ok := c.(card.DiscountPerRunechant); ok {
		m.isDiscount = true
		m.printedCost = d.PrintedCost()
	}
	cardMetaCache[id] = m
	atomic.StoreUint32(&cardMetaReady[id], 1)
	return m
}

// attackBufs holds pre-allocated buffers for the attack-evaluation pipeline (bestSequence →
// playSequence) and the partition loop in bestUncached. Allocated once and cached globally so a
// whole deck evaluation reuses the same buffers across every partition, mask, and permutation.
type attackBufs struct {
	perm           []card.Card
	pcBuf          []card.PlayedCard
	ptrBuf         []*card.PlayedCard
	cardsPlayedBuf []card.Card
	state          *card.TurnState
	attackerBuf    []card.Card // for bestAttackWithWeapons mask iteration
	// Pre-computed weapon data: weaponCosts[mask] is the total Cost of weapons in that mask;
	// weaponNames[mask] is the pre-built []string of weapon names for the mask. Indexed by bitmask
	// (0 to 2^len(weapons)-1).
	weaponCosts []int
	weaponNames [][]string
	// permMeta is parallel to perm: it carries precomputed scalar metadata for each attacker so
	// playSequence's inner loop can skip interface dispatch on Types / Cost / GoAgain / the
	// DiscountPerRunechant type-assertion. Populated and permuted in lockstep with perm by
	// bestSequence; read in order by playSequence.
	permMeta []attackerMeta
	// Partition-loop buffers, consumed by bestUncached. Sized handSize+1 to cover the optional
	// arsenal-in slot the enumerator treats as index n. defendPrintedVals holds the PrintedCost of
	// DiscountPerRunechant defense reactions for the leaf's post-attack discount re-pricing;
	// non-discount cards get 0. isDRBuf caches TypeDefenseReaction membership so the leaf skips
	// Types().Has calls.
	rolesBuf          []Role
	pitchVals         []int
	costVals          []int
	defendCostVals    []int
	defendPrintedVals []int
	defenseVals       []int
	isDRBuf           []bool
	// pitchedValsScratch backs the per-leaf "pitched values" slice the feasibility check
	// inspects. Re-sliced to [:0] at the start of every leaf and appended to inside it, so
	// lifting the alloc to bufs eliminates one make([]int, …, totalN) per bestUncached call.
	pitchedValsScratch []int
	pitchedBuf         []card.Card
	attackersBuf       []card.Card
	defendersBuf       []card.Card
	// perCardScratch is sized maxAttackers (handSize + weaponCount). Only written by playSequence
	// when the caller passes a non-nil perCardOut, and read by bestSequence to snapshot the
	// winning permutation's per-card damage into the caller's output buffer. Untracked callers
	// (the partition-loop hot path) pass nil and never touch this slice.
	perCardScratch []float64
	// perCardTriggerScratch parallels perCardScratch for hero-trigger damage (OnCardPlayed
	// return). Same sizing / same "only written when tracking" discipline.
	perCardTriggerScratch []float64
	// fillContribWinnerOrder / fillContribPerCard are output buffers for bestSequence when
	// fillContributions runs the tracked replay after the partition loop. Kept on attackBufs so
	// each Best call reuses the same underlying slab instead of allocating a fresh pair.
	fillContribWinnerOrder []card.Card
	fillContribPerCard     []float64
	fillContribTriggerDmg  []float64
	// fillContribUsed marks hand indices already assigned during chain→hand mapping. Sized
	// handSize; the caller resets it with clear before each fillContributions pass.
	fillContribUsed []bool
}

func newAttackBufs(handSize, weaponCount int, weapons []weapon.Weapon) *attackBufs {
	// +1 reserves a slot for the arsenal-in card, which joins `attackers` or `defenders` when
	// the partition enumerator plays it out of arsenal. Without it, all-attack hands + arsenal
	// overflow attackerBuf (slice bounds panic in bestAttackWithWeapons).
	maxAttackers := handSize + weaponCount + 1
	numMasks := 1 << weaponCount
	weaponCosts := make([]int, numMasks)
	weaponNames := make([][]string, numMasks)
	for mask := 0; mask < numMasks; mask++ {
		cost := 0
		var names []string
		for i, w := range weapons {
			if mask&(1<<i) != 0 {
				cost += w.Cost()
				names = append(names, w.Name())
			}
		}
		weaponCosts[mask] = cost
		weaponNames[mask] = names
	}
	return &attackBufs{
		perm:           make([]card.Card, maxAttackers),
		permMeta:       make([]attackerMeta, maxAttackers),
		pcBuf:          make([]card.PlayedCard, maxAttackers),
		ptrBuf:         make([]*card.PlayedCard, maxAttackers),
		cardsPlayedBuf: make([]card.Card, 0, maxAttackers),
		state:          &card.TurnState{},
		attackerBuf:    make([]card.Card, maxAttackers),
		weaponCosts:    weaponCosts,
		weaponNames:    weaponNames,
		rolesBuf:          make([]Role, handSize+1),
		pitchVals:         make([]int, handSize+1),
		costVals:          make([]int, handSize+1),
		defendCostVals:    make([]int, handSize+1),
		defendPrintedVals: make([]int, handSize+1),
		defenseVals:       make([]int, handSize+1),
		isDRBuf:            make([]bool, handSize+1),
		pitchedValsScratch: make([]int, 0, handSize+1),
		pitchedBuf:             make([]card.Card, 0, handSize+1),
		attackersBuf:           make([]card.Card, 0, handSize+1),
		defendersBuf:           make([]card.Card, 0, handSize+1),
		perCardScratch:         make([]float64, maxAttackers),
		perCardTriggerScratch:  make([]float64, maxAttackers),
		fillContribWinnerOrder: make([]card.Card, maxAttackers),
		fillContribPerCard:     make([]float64, maxAttackers),
		fillContribTriggerDmg:  make([]float64, maxAttackers),
		fillContribUsed:        make([]bool, handSize),
	}
}

// getAttackBufs returns this Evaluator's scratch-buffer set, rebuilding it when (handSize,
// weapons) changes. The cache is single-slot per Evaluator — iterate runs tens of thousands of
// hands with the same shape, so a slot outperforms a keyed pool for our workload.
func (e *Evaluator) getAttackBufs(handSize int, weapons []weapon.Weapon) *attackBufs {
	var wids [2]card.ID
	for i, w := range weapons {
		if i >= len(wids) {
			break
		}
		wids[i] = w.ID()
	}
	if e.bufsValid &&
		e.bufsHandSize == handSize &&
		e.bufsWeaponCount == len(weapons) &&
		e.bufsWeaponIDs == wids {
		return e.bufs
	}
	e.bufs = newAttackBufs(handSize, len(weapons), weapons)
	e.bufsHandSize = handSize
	e.bufsWeaponCount = len(weapons)
	e.bufsWeaponIDs = wids
	e.bufsValid = true
	return e.bufs
}

func (e *Evaluator) bestUncached(hero hero.Hero, weapons []weapon.Weapon, hand []card.Card, incomingDamage int, deck []card.Card, runechantCarryover int, arsenalCardIn card.Card) TurnSummary {
	n := len(hand)
	// The partition recurse treats the arsenal-in card as an extra entry at index n with a
	// restricted role menu (Arsenal / Attack / Defend), so everything about it — its cost,
	// its damage, whether it stays or plays — is decided inside the enumeration rather than
	// by a wrapping enumerator. totalN is the effective size of BestLine.
	totalN := n
	if arsenalCardIn != nil {
		totalN = n + 1
	}

	// Seed best.LeftoverRunechants with the carryover — partitions with no attacks don't reduce
	// the count, so carryover is the natural baseline to beat. BestLine is initialised with
	// every hand card Held and the arsenal-in card (if any) staying in the slot, so a hand with
	// no Value-adding partition (everything pruned by cost/pitch-feasibility) still reports
	// sensible assignments — nothing got played or pitched.
	best := TurnSummary{BestLine: make([]CardAssignment, totalN), LeftoverRunechants: runechantCarryover}
	// bestSwung holds the winning partition's swung weapon names so fillContributions can
	// rebuild the chain it runs bestSequence over. It lives outside TurnSummary because
	// weapons are recoverable from AttackChain once fillContributions finishes.
	var bestSwung []string
	for i := 0; i < n; i++ {
		best.BestLine[i] = CardAssignment{Card: hand[i], Role: Held}
	}
	if arsenalCardIn != nil {
		best.BestLine[n] = CardAssignment{Card: arsenalCardIn, Role: Arsenal, FromArsenal: true}
		best.ArsenalCard = arsenalCardIn
	}

	// bufs is the pooled scratch space for this deck evaluation (see getAttackBufs). Partition
	// scratch is sized handSize+1, big enough for totalN when an arsenal-in card inflates the
	// effective hand. Each field is reset via the totalN-slice init loop below, so carry-over
	// bytes from prior calls never leak into this partition's values.
	bufs := e.getAttackBufs(n, weapons)
	rolesBuf := bufs.rolesBuf[:totalN]
	pvals := bufs.pitchVals[:totalN]
	cvals := bufs.costVals[:totalN]
	dvals := bufs.defenseVals[:totalN]
	dCostVals := bufs.defendCostVals[:totalN]
	dPrintedVals := bufs.defendPrintedVals[:totalN]
	isDR := bufs.isDRBuf[:totalN]

	// Pre-compute per-card pitch / cost / defense values so the recurse doesn't re-invoke
	// card-method interface calls on each partition leaf. defendCostVals holds Cost only for
	// Defense Reactions; non-reactions block for free. defendPrintedVals holds PrintedCost for
	// DiscountPerRunechant defenders and zero otherwise — used by the post-attack discount
	// re-pricing check at the leaf. Both are populated conditionally below, so clear prior-call
	// residue from the pooled scratch before the loop.
	clear(dCostVals)
	clear(dPrintedVals)
	hasReactions := false
	hasDiscountReactions := false
	for i := 0; i < totalN; i++ {
		var c card.Card
		if i < n {
			c = hand[i]
		} else {
			c = arsenalCardIn
		}
		pvals[i] = c.Pitch()
		cvals[i] = c.Cost()
		dvals[i] = c.Defense()
		isDR[i] = c.Types().Has(card.TypeDefenseReaction)
		if isDR[i] {
			dCostVals[i] = cvals[i]
			hasReactions = true
			if dp, ok := c.(card.DiscountPerRunechant); ok {
				dPrintedVals[i] = dp.PrintedCost()
				hasDiscountReactions = true
			}
		}
	}
	pitched := bufs.pitchedBuf
	attackers := bufs.attackersBuf
	defenders := bufs.defendersBuf

	// recurse tracks defenderCostSum separately so the leaf can hand the attack pipeline a
	// resource budget of (pitchSum - defenderCostSum) to deduct chain-card effective costs from.
	// arsenalCount is at most 1 across the whole partition; any branch that would push it past
	// one is pruned at the role-selection step. pitchedVals is scratch reused across leaves;
	// lifted to bufs.pitchedValsScratch so it doesn't allocate on every bestUncached call.
	pitchedValsScratch := bufs.pitchedValsScratch[:0]
	var recurse func(i, pitchSum, costSum, defenseSum, defenderCostSum, arsenalCount int)
	recurse = func(i, pitchSum, costSum, defenseSum, defenderCostSum, arsenalCount int) {
		if i == totalN {
			attackCardCost := costSum - defenderCostSum
			drCost := defenderCostSum
			pitchedVals := pitchedValsScratch[:0]
			for j := 0; j < totalN; j++ {
				if rolesBuf[j] == Pitch {
					pitchedVals = append(pitchedVals, pvals[j])
				}
			}
			if !anyMaskFeasible(pitchedVals, attackCardCost, drCost, bufs.weaponCosts, len(weapons)) {
				return
			}
			prevented := defenseSum
			if prevented > incomingDamage {
				prevented = incomingDamage
			}
			// Group roles into played / pitched / defending buckets. The grouping iterates the
			// hand (size n) for the usual buckets, then layers in the arsenal slot (index n)
			// based on its assigned role. Arsenal-role cards — whether from hand or from the
			// slot — contribute nothing this turn.
			var p, a []card.Card
			var defenseDealt int
			hasAnyDefender := hasReactions
			if !hasAnyDefender && arsenalCardIn != nil && rolesBuf[n] == Defend {
				hasAnyDefender = true
			}
			if hasAnyDefender {
				var d []card.Card
				p, a, d = groupByRoleInto(hand, rolesBuf[:n], pitched[:0], attackers[:0], defenders[:0])
				if arsenalCardIn != nil {
					switch rolesBuf[n] {
					case Attack:
						a = append(a, arsenalCardIn)
					case Defend:
						d = append(d, arsenalCardIn)
					}
				}
				defenseDealt = defenseReactionDamage(d, p, deck, bufs.state)
			} else {
				p, a = groupPitchAttack(hand, rolesBuf[:n], pitched[:0], attackers[:0])
				if arsenalCardIn != nil && rolesBuf[n] == Attack {
					a = append(a, arsenalCardIn)
				}
			}
			attackDealt, leftoverRunechants, residualBudget, swung := bestAttackWithWeapons(hero, weapons, a, p, deck, bufs, pitchSum, costSum, defenderCostSum, runechantCarryover)

			// DiscountPerRunechant defense reactions reserved 0 in defenderCostSum (their Cost()
			// is the fully-discounted minimum). Re-price them now that the attack chain has
			// resolved and left `leftoverRunechants` runechants available. Arsenal-in defenders
			// aren't checked here; the roster has no DiscountPerRunechant cards that would also
			// want to live in the arsenal slot, so the cheap hand-only scan suffices.
			if hasDiscountReactions {
				extraCost := 0
				for j := 0; j < n; j++ {
					if rolesBuf[j] != Defend || dPrintedVals[j] == 0 {
						continue
					}
					actualCost := dPrintedVals[j] - leftoverRunechants
					if actualCost < 0 {
						actualCost = 0
					}
					extraCost += actualCost
				}
				if extraCost > residualBudget {
					return
				}
			}

			v := attackDealt + defenseDealt + prevented
			arsenalCard := findArsenalCard(rolesBuf, hand, arsenalCardIn, n, totalN)
			if !beatsBest(v, leftoverRunechants, arsenalCard != nil, best) {
				return
			}
			best.Value = v
			bestSwung = swung
			best.LeftoverRunechants = leftoverRunechants
			best.ArsenalCard = arsenalCard
			// Write the winning roles into BestLine. Cards and FromArsenal flags were populated
			// at construction; only Role varies per partition. Contribution is cleared here and
			// written by fillContributions below for the winning line.
			for j := 0; j < totalN; j++ {
				best.BestLine[j].Role = rolesBuf[j]
				best.BestLine[j].Contribution = 0
			}
			return
		}
		isArsenalSlot := i == n && arsenalCardIn != nil
		for r := Role(0); r <= Arsenal; r++ {
			if !roleAllowed(r, isArsenalSlot, isDR[i], arsenalCount) {
				continue
			}
			rolesBuf[i] = r
			switch r {
			case Pitch:
				recurse(i+1, pitchSum+pvals[i], costSum, defenseSum, defenderCostSum, arsenalCount)
			case Attack:
				recurse(i+1, pitchSum, costSum+cvals[i], defenseSum, defenderCostSum, arsenalCount)
			case Defend:
				recurse(i+1, pitchSum, costSum+dCostVals[i], defenseSum+dvals[i], defenderCostSum+dCostVals[i], arsenalCount)
			case Held:
				recurse(i+1, pitchSum, costSum, defenseSum, defenderCostSum, arsenalCount)
			case Arsenal:
				recurse(i+1, pitchSum, costSum, defenseSum, defenderCostSum, arsenalCount+1)
			}
		}
	}
	recurse(0, 0, 0, 0, 0, 0)
	// Once per Best call, on the winning line only, attribute per-card contribution.
	if len(best.BestLine) > 0 {
		fillContributions(&best, hero, weapons, bestSwung, deck, bufs, incomingDamage, runechantCarryover)
	}
	return best
}

// canCoverPhasesAllUsed decides whether every pitched value can be split between the attack
// phase (covering attackCost) and the defence phase (covering drCost) while respecting FaB's
// pitch-timing rule: a card can only be pitched while some played card on the stack is still
// unpaid. That means the Pitch-role cards the partition has committed to must all pay for
// something — any "extra" card would have to be Held instead.
//
// Per-phase legality uses a sufficient condition: if sum(pool) == phaseCost the pool exactly
// covers its costs (assumed legal); if sum(pool) > phaseCost the excess has to be absorbable in
// a single over-paying pitch, so max(pool) must strictly exceed the excess. Partitions that
// would require multiple pitches to each push a cost above full are rejected.
//
// With both phases having positive cost we enumerate every non-empty, non-full attack-pool mask
// (2^k - 2 for k pitched cards) and accept the first that satisfies both phase checks. k is
// bounded by the hand size so this stays cheap relative to the outer 4^n partition search.
func canCoverPhasesAllUsed(pitchedVals []int, attackCost, drCost int) bool {
	k := len(pitchedVals)
	if k == 0 {
		// No pitches. Legal only if both phases cost nothing.
		return attackCost == 0 && drCost == 0
	}
	if attackCost == 0 && drCost == 0 {
		// Nothing to pitch for — any Pitch-role card would have been Held instead.
		return false
	}
	total := 0
	for _, v := range pitchedVals {
		total += v
	}
	if total < attackCost+drCost {
		return false
	}
	full := uint32(1<<uint(k)) - 1
	if attackCost == 0 {
		return phaseLegal(pitchedVals, full, drCost)
	}
	if drCost == 0 {
		return phaseLegal(pitchedVals, full, attackCost)
	}
	// Both phases have cost, so each pool must be non-empty.
	for aMask := uint32(1); aMask < full; aMask++ {
		if !phaseLegal(pitchedVals, aMask, attackCost) {
			continue
		}
		if phaseLegal(pitchedVals, full^aMask, drCost) {
			return true
		}
	}
	return false
}

// phaseLegal returns true iff the pitch values selected by subsetMask can legally cover phaseCost
// for a single phase. sum < phaseCost means the pool can't pay; sum == phaseCost is trivially
// legal (exact coverage); sum > phaseCost needs one pitch large enough to absorb the whole
// over-pay in a single final pitch for some card in the phase.
func phaseLegal(pitchedVals []int, subsetMask uint32, phaseCost int) bool {
	sum, maxP := 0, 0
	for i, v := range pitchedVals {
		if subsetMask&(1<<uint(i)) != 0 {
			sum += v
			if v > maxP {
				maxP = v
			}
		}
	}
	if sum < phaseCost {
		return false
	}
	if sum == phaseCost {
		return true
	}
	return maxP > sum-phaseCost
}

// groupByRoleInto is like groupByRole but appends into caller-provided slices (which should be
// passed pre-reset to length 0) to avoid per-partition heap allocation.
func groupByRoleInto(hand []card.Card, roles []Role, pitched, attackers, defenders []card.Card) ([]card.Card, []card.Card, []card.Card) {
	for i, c := range hand {
		switch roles[i] {
		case Pitch:
			pitched = append(pitched, c)
		case Attack:
			attackers = append(attackers, c)
		case Defend:
			defenders = append(defenders, c)
		}
	}
	return pitched, attackers, defenders
}

// groupPitchAttack is the reaction-free leaf's grouping step: skips the defenders bucket since
// we only need it for Defense-Reaction-Play dispatch, which this path doesn't run.
func groupPitchAttack(hand []card.Card, roles []Role, pitched, attackers []card.Card) ([]card.Card, []card.Card) {
	for i, c := range hand {
		switch roles[i] {
		case Pitch:
			pitched = append(pitched, c)
		case Attack:
			attackers = append(attackers, c)
		}
	}
	return pitched, attackers
}

// anyMaskFeasible returns true if at least one weapon-swing subset can be paid for by some split
// of the pitched values between the attack phase and the defense phase. Called at every
// partition leaf — phase-legality is the cheap-to-check screen that lets the expensive
// bestAttackWithWeapons work run only on partitions that could legally pay for something.
func anyMaskFeasible(pitchedVals []int, attackCardCost, drCost int, weaponCosts []int, weaponCount int) bool {
	masks := 1 << weaponCount
	for mask := 0; mask < masks; mask++ {
		if canCoverPhasesAllUsed(pitchedVals, attackCardCost+weaponCosts[mask], drCost) {
			return true
		}
	}
	return false
}

// findArsenalCard picks out the card assigned the Arsenal role in the current partition, if any.
// Index j < n points into the hand; j == n points at the arsenal-in card. Returns nil when no
// slot in the partition is Arsenal.
func findArsenalCard(rolesBuf []Role, hand []card.Card, arsenalCardIn card.Card, n, totalN int) card.Card {
	for j := 0; j < totalN; j++ {
		if rolesBuf[j] != Arsenal {
			continue
		}
		if j < n {
			return hand[j]
		}
		return arsenalCardIn
	}
	return nil
}

// beatsBest decides whether a candidate partition displaces the current best under the solver's
// tiebreak rules: higher Value first, then more leftover runechants (they convert to future
// arcane damage), then preferring a partition that fills the arsenal slot (arsenal saves a hand
// slot in the next refill — smaller but real upside).
func beatsBest(v, leftoverRunechants int, hasArsenal bool, best TurnSummary) bool {
	if v > best.Value {
		return true
	}
	if v < best.Value {
		return false
	}
	if leftoverRunechants > best.LeftoverRunechants {
		return true
	}
	if leftoverRunechants < best.LeftoverRunechants {
		return false
	}
	return hasArsenal && best.ArsenalCard == nil
}

// roleAllowed decides whether the partition enumerator may assign role r to the current card. The
// arsenal slot (i == n, arsenal-in card present) may only take Arsenal (stay), Attack (any
// non-DR card — auras and non-attack actions play fine from arsenal on your turn), or Defend
// (Defense Reactions only — plain-blocking from arsenal isn't allowed). Hand cards take any role
// except Attack for Defense Reactions (strict FaB timing — DRs only fire on the opponent's turn).
// arsenalCount tracks how many Arsenal assignments the partition already has; the cap is 1.
func roleAllowed(r Role, isArsenalSlot, isDefenseReaction bool, arsenalCount int) bool {
	if r == Arsenal && arsenalCount >= 1 {
		return false
	}
	if isArsenalSlot {
		switch r {
		case Pitch, Held:
			return false
		case Attack:
			return !isDefenseReaction
		case Defend:
			return isDefenseReaction
		}
		return true
	}
	return !(r == Attack && isDefenseReaction)
}

// defenseReactionDamage runs the Play() hook of every Defense Reaction in `defenders` and sums
// the damage they deal back to the attacker (e.g. Weeping Battleground's 1 arcane on banish).
// Played in isolation — no attack ordering; TurnState only carries Pitched/Deck so effects that
// check "what was pitched" work. Uncapped: this damage is dealt, not prevented.
//
// state is caller-provided (from attackBufs) and reset per call. Passing a reused pointer lets
// the state stay on the heap-allocated buffer rather than escaping a fresh stack value through
// the interface method on every partition leaf.
func defenseReactionDamage(defenders, pitched, deck []card.Card, state *card.TurnState) int {
	total := 0
	for _, d := range defenders {
		if !d.Types().Has(card.TypeDefenseReaction) {
			continue
		}
		*state = card.TurnState{Pitched: pitched, Deck: deck}
		total += d.Play(state)
	}
	return total
}

// bestAttackWithWeapons enumerates every subset of `weapons` to swing alongside `attackers` and
// returns the max damage over all (affordable) weapon masks, plus the runechant leftover, the
// residual chain budget (pitch not consumed by the winning line — used by bestUncached to re-check
// DiscountPerRunechant defense affordability), and the swung weapons (in input order).
//
// resourceBudget passed down to the chain pipeline is pitchSum - defenderCostSum - weaponCosts
// for this mask; the chain then deducts each attacker's effective cost (which for
// DiscountPerRunechant cards is max(0, PrintedCost() - runechantCount at play-time)) and rejects
// orderings that run negative.
func bestAttackWithWeapons(hero hero.Hero, weapons []weapon.Weapon, attackers, pitched, deck []card.Card, bufs *attackBufs, pitchSum, costSum, defenderCostSum, runechantCarryover int) (int, int, int, []string) {
	// Build the sequence-eval context once per partition leaf. All three values (resourceBudget,
	// runechantCarryover, and the stable pitched/deck refs) are constant across the weapon-mask
	// loop, so sharing one context keeps the inner calls narrow.
	ctx := &sequenceContext{
		hero:               hero,
		pitched:            pitched,
		deck:               deck,
		bufs:               bufs,
		resourceBudget:     pitchSum - defenderCostSum,
		runechantCarryover: runechantCarryover,
	}
	best := 0
	bestLeftoverRunechants := runechantCarryover
	bestResidualBudget := ctx.resourceBudget
	var bestSwung []string
	// Reuse the shared attacker buffer across mask iterations.
	copy(bufs.attackerBuf, attackers)
	// pitchedVals is derived from the pitched cards' Pitch() so the per-mask phase-feasibility
	// check can include each mask's weapon cost in the attack-phase total. A locally-allocated
	// slice keeps this call re-entrant (Best can be invoked concurrently through memoCache once
	// that's exposed).
	attackCardCost := costSum - defenderCostSum
	drCost := defenderCostSum
	pitchedVals := make([]int, len(pitched))
	for i, c := range pitched {
		pitchedVals[i] = c.Pitch()
	}
	for mask := 0; mask < 1<<len(weapons); mask++ {
		// bufs.weaponCosts[mask] is the pre-summed Cost of the selected weapons — avoids an
		// interface dispatch per weapon on every mask. The phase-feasibility check ensures every
		// Pitch-role card can legally pay for something (attack cost + weapon cost + DR cost).
		if !canCoverPhasesAllUsed(pitchedVals, attackCardCost+bufs.weaponCosts[mask], drCost) {
			continue
		}
		allAttackers := bufs.attackerBuf[:len(attackers)]
		for i, w := range weapons {
			if mask&(1<<i) != 0 {
				allAttackers = append(allAttackers, w)
			}
		}
		// Every sequence card (attackers AND the swung weapons) deducts its own effective cost
		// from resourceBudget inside playSequence, so we don't pre-deduct weapon cost here.
		dealt, leftoverRunechants, residualBudget := ctx.bestSequence(allAttackers, nil, nil, nil)
		// Prefer higher damage; on ties prefer more leftover runechants; then more residual
		// budget — both are extra slack that can enable discount defense reactions.
		if dealt > best ||
			(dealt == best && leftoverRunechants > bestLeftoverRunechants) ||
			(dealt == best && leftoverRunechants == bestLeftoverRunechants && residualBudget > bestResidualBudget) {
			best = dealt
			bestLeftoverRunechants = leftoverRunechants
			bestResidualBudget = residualBudget
			bestSwung = bufs.weaponNames[mask]
		}
	}
	return best, bestLeftoverRunechants, bestResidualBudget, bestSwung
}

// sequenceContext carries the stable per-partition-leaf environment for evaluating played-card
// sequences: the hero (for OnCardPlayed triggers), the pitched / deck reference slices read by
// Card.Play, the shared scratch buffers, and the numeric budgets that persist across permutation
// and mask iterations. One context is built once per leaf and used for every permutation the
// solver evaluates below it, so the hot inner calls (playSequence, bestSequence) shrink to just
// their varying inputs and tracking outputs.
type sequenceContext struct {
	hero               hero.Hero
	pitched, deck      []card.Card
	bufs               *attackBufs
	resourceBudget     int
	runechantCarryover int
}

// bestSequence tries every ordering of attackers and returns the max total damage after Play
// is called on each in turn plus the runechant count at the end of that winning permutation.
// Between each card's Play() and its append to CardsPlayed, the hero's OnCardPlayed hook fires
// so triggered abilities (e.g. Viserai's Runechants) contribute.
//
// Uses Heap's algorithm (iterative) for permutation generation. That saves a closure/callback
// allocation and a recursive call per permutation vs. a callback-style permuter.
//
// When winnerOrderOut is non-nil (len >= len(attackers)) the winning permutation is copied into
// it. perCardOut / perCardTriggerOut (same size rule) receive the winning line's per-card Play
// damage and per-card hero-trigger damage respectively. Used once per Best call by
// fillContributions; the hot partition-loop caller (bestAttackWithWeapons) passes nil for all
// three so the permutation search stays allocation-free.
func (ctx *sequenceContext) bestSequence(attackers, winnerOrderOut []card.Card, perCardOut, perCardTriggerOut []float64) (int, int, int) {
	n := len(attackers)
	if n == 0 {
		return 0, ctx.runechantCarryover, ctx.resourceBudget
	}
	perm := ctx.bufs.perm[:n]
	permMeta := ctx.bufs.permMeta[:n]
	copy(perm, attackers)
	for idx, c := range attackers {
		permMeta[idx] = attackerMetaFor(c)
	}

	// Scratch buffers are playSequence's per-card outputs, overwritten on every permutation; on
	// a new winner we copy them into the caller's perCardOut / perCardTriggerOut. Only used
	// when the caller asked to track.
	var scratch, triggerScratch []float64
	if perCardOut != nil {
		scratch = ctx.bufs.perCardScratch[:n]
	}
	if perCardTriggerOut != nil {
		if cap(ctx.bufs.perCardTriggerScratch) < n {
			ctx.bufs.perCardTriggerScratch = make([]float64, n)
		}
		triggerScratch = ctx.bufs.perCardTriggerScratch[:n]
	}

	best := 0
	bestLeftoverRunechants := ctx.runechantCarryover
	bestResidualBudget := ctx.resourceBudget
	eval := func() {
		dmg, leftoverRunechants, residualBudget, legal := ctx.playSequenceWithMeta(perm, scratch, triggerScratch)
		if !legal {
			return
		}
		// Prefer higher damage; on ties prefer more leftover runechants; then more residual
		// budget — both are extra slack that can enable discount defense reactions.
		if dmg > best ||
			(dmg == best && leftoverRunechants > bestLeftoverRunechants) ||
			(dmg == best && leftoverRunechants == bestLeftoverRunechants && residualBudget > bestResidualBudget) {
			best = dmg
			bestLeftoverRunechants = leftoverRunechants
			bestResidualBudget = residualBudget
			if winnerOrderOut != nil {
				copy(winnerOrderOut[:n], perm)
			}
			if perCardOut != nil {
				copy(perCardOut[:n], scratch)
			}
			if perCardTriggerOut != nil {
				copy(perCardTriggerOut[:n], triggerScratch)
			}
		}
	}
	eval()
	// Heap's algorithm, non-recursive: c[] counts how many times each stack frame has iterated.
	// perm and permMeta are swapped together so playSequence always sees meta aligned with the
	// current permutation.
	var c [8]int
	i := 0
	for i < n {
		if c[i] < i {
			if i&1 == 0 {
				perm[0], perm[i] = perm[i], perm[0]
				permMeta[0], permMeta[i] = permMeta[i], permMeta[0]
			} else {
				perm[c[i]], perm[i] = perm[i], perm[c[i]]
				permMeta[c[i]], permMeta[i] = permMeta[i], permMeta[c[i]]
			}
			eval()
			c[i]++
			i = 0
		} else {
			c[i] = 0
			i++
		}
	}
	return best, bestLeftoverRunechants, bestResidualBudget
}

// playSequence plays `order` as a sequence of cards, reusing ctx.bufs' pooled buffers to avoid
// per-permutation heap allocation. The buffers are mutated in place; the caller must not read
// them concurrently.
//
// When perCardOut is non-nil (len >= n) each entry is set to the card's Play return for that
// position; perCardTriggerOut (same size rule) receives the hero OnCardPlayed return for that
// same position. The hot partition-loop callers pass nil for both to skip these writes; the
// winning-line replay from fillContributions passes real slices.
//
// Runechant flow:
//   - state.Runechants starts at ctx.runechantCarryover (tokens from the previous turn).
//   - Each card's Play / hero OnCardPlayed may call CreateRunechants, incrementing the count AND
//     returning n damage — tokens are credited exactly once, at creation.
//   - After each Attack- or Weapon-typed card's Play+OnCardPlayed resolve, all current tokens
//     fire and are destroyed: state.Runechants is zeroed but damage is NOT re-added (that would
//     double-count tokens whose value was already credited on creation).
//   - At end of the sequence, state.Runechants is the leftover count that carries into the next
//     turn.
//
// Resource flow:
//   - ctx.resourceBudget is the starting pool; each card deducts its effective cost. For cards
//     implementing DiscountPerRunechant, effective cost is max(0, PrintedCost() - state.
//     Runechants) at the moment it's played; for everyone else it's Cost(). A negative
//     remaining budget returns legal=false (the caller treats this ordering as zero damage).
// playSequence populates permMeta from order and then calls playSequenceWithMeta. Test-only
// callers and fillContributions use this when they don't go through bestSequence's permutation
// loop. The hot path (bestSequence) builds meta once and calls playSequenceWithMeta directly so
// the interface dispatch for scalar attributes amortises across the N! permutations it evaluates.
func (ctx *sequenceContext) playSequence(order []card.Card, perCardOut, perCardTriggerOut []float64) (damage int, leftoverRunechants int, residualBudget int, legal bool) {
	meta := ctx.bufs.permMeta[:len(order)]
	for i, c := range order {
		meta[i] = attackerMetaFor(c)
	}
	return ctx.playSequenceWithMeta(order, perCardOut, perCardTriggerOut)
}

// playSequenceWithMeta runs a specific attacker ordering and returns damage + leftover runechants.
// Assumes ctx.bufs.permMeta[:len(order)] holds metadata aligned with order; the caller is
// responsible for keeping them in lockstep (bestSequence swaps meta whenever it swaps perm).
func (ctx *sequenceContext) playSequenceWithMeta(order []card.Card, perCardOut, perCardTriggerOut []float64) (damage int, leftoverRunechants int, residualBudget int, legal bool) {
	n := len(order)
	pcBuf := ctx.bufs.pcBuf
	ptrBuf := ctx.bufs.ptrBuf
	meta := ctx.bufs.permMeta[:n]
	for i, c := range order {
		pcBuf[i] = card.PlayedCard{Card: c}
		ptrBuf[i] = &pcBuf[i]
		if perCardOut != nil {
			perCardOut[i] = 0
		}
		if perCardTriggerOut != nil {
			perCardTriggerOut[i] = 0
		}
	}
	played := ptrBuf[:n]
	state := ctx.bufs.state
	// Reset only the fields playSequence is allowed to mutate. Pitched and Deck are stable across
	// permutations within a partition leaf; CardsRemaining and Self are rewritten per-card in the
	// loop below. A full-struct replace ran in the profile at ~850ms / call because every field
	// (including the big slice headers) got memcpy'd from a zero template; skipping the fields
	// the caller doesn't read afterward shaves most of that.
	state.Pitched = ctx.pitched
	state.Deck = ctx.deck
	state.CardsPlayed = ctx.bufs.cardsPlayedBuf[:0]
	state.Runechants = ctx.runechantCarryover
	state.DelayedRunechants = 0
	state.ArcaneDamageDealt = false
	state.AuraCreated = false
	state.Overpower = false
	resources := ctx.resourceBudget
	for i, pc := range played {
		m := meta[i]
		// Effective cost: discount cards drop by runechant count (floored at 0); others pay printed.
		var effCost int
		if m.isDiscount {
			effCost = m.printedCost - state.Runechants
			if effCost < 0 {
				effCost = 0
			}
		} else {
			effCost = m.cost
		}
		resources -= effCost
		if resources < 0 {
			return 0, 0, 0, false
		}

		state.CardsRemaining = played[i+1:]
		state.Self = pc

		// If this card is an attack or weapon and any Runechant is currently live, those tokens
		// will fire on its damage step. Set ArcaneDamageDealt now — *before* Play and the hero's
		// OnCardPlayed trigger — so Play effects that read "if you've dealt arcane damage this
		// turn" see the flag for same-hand triggers. Cards that deal arcane damage directly via
		// their Play text flip the flag themselves inside Play.
		isAttackOrWeapon := m.isAttackOrWeapon
		if isAttackOrWeapon && state.Runechants > 0 {
			state.ArcaneDamageDealt = true
		}

		playDmg := pc.Card.Play(state)
		triggerDmg := ctx.hero.OnCardPlayed(pc.Card, state)
		damage += playDmg + triggerDmg
		if perCardOut != nil {
			perCardOut[i] = float64(playDmg)
		}
		if perCardTriggerOut != nil {
			perCardTriggerOut[i] = float64(triggerDmg)
		}
		state.CardsPlayed = append(state.CardsPlayed, pc.Card)

		// Attacks and weapon swings consume all runechants in play. Damage isn't re-added here:
		// each token was already credited as +1 at creation time (see CreateRunechants), so
		// consuming them is purely state cleanup.
		if isAttackOrWeapon {
			state.Runechants = 0
		}

		if i < n-1 && !(m.baseGoAgain || pc.GrantedGoAgain) {
			return 0, 0, 0, false
		}
	}
	// Delayed tokens (e.g. from Blessing of Occult) skip this turn and go straight to next
	// turn's carryover.
	return damage, state.Runechants + state.DelayedRunechants, resources, true
}

// fillContributions populates each BestLine entry's Contribution from the winning line. Pitch
// role cards credit Card.Pitch() (resource value); Defend role cards credit their proportional
// share of Prevented blocking plus their own Play return if they're a defense reaction; Attack
// role cards credit the per-card damage tracked during the winning attack-chain replay.
// Contributions on Held / Arsenal entries stay at zero — those cards contributed nothing this
// turn.
//
// Called once per Best call from bestUncached after the partition loop picks the winner. All
// transient slices (pitched/attackers/chain/winnerOrder/perCard/used) borrow attackBufs slots
// so nothing allocates here.
func fillContributions(summary *TurnSummary, hero hero.Hero, weapons []weapon.Weapon, swungNames []string, deck []card.Card, bufs *attackBufs, incomingDamage, runechantCarryover int) {
	line := summary.BestLine
	total := len(line)

	// Reconstruct pitched, attackers, swung weapons, and the sequence's resourceBudget from the
	// winning line. The arsenal-in entry (FromArsenal=true, last slot) participates in
	// attackers / defenders when its role is Attack / Defend, identically to hand entries.
	pitched := bufs.pitchedBuf[:0]
	attackers := bufs.attackersBuf[:0]
	var sumDef int
	pitchSum := 0
	defenderCostSum := 0
	for _, a := range line {
		switch a.Role {
		case Pitch:
			pitched = append(pitched, a.Card)
			pitchSum += a.Card.Pitch()
		case Attack:
			attackers = append(attackers, a.Card)
		case Defend:
			sumDef += a.Card.Defense()
			if a.Card.Types().Has(card.TypeDefenseReaction) {
				if _, disc := a.Card.(card.DiscountPerRunechant); !disc {
					defenderCostSum += a.Card.Cost()
				}
			}
		}
	}
	resourceBudget := pitchSum - defenderCostSum

	// Pitch contributions.
	for i := range line {
		if line[i].Role == Pitch {
			line[i].Contribution = float64(line[i].Card.Pitch())
		}
	}

	// Defense: prevented share (proportional to each defender's Defense()) plus defense-reaction
	// Play return when applicable.
	prevented := sumDef
	if prevented > incomingDamage {
		prevented = incomingDamage
	}
	for i := range line {
		if line[i].Role != Defend {
			continue
		}
		c := line[i].Card
		if sumDef > 0 {
			line[i].Contribution = float64(c.Defense()) * float64(prevented) / float64(sumDef)
		}
		if c.Types().Has(card.TypeDefenseReaction) {
			*bufs.state = card.TurnState{Pitched: pitched, Deck: deck}
			line[i].Contribution += float64(c.Play(bufs.state))
		}
	}

	// Attack chain: re-run the sequence search with tracking turned on so we recover the
	// winning permutation and each chain position's contribution. Weapons in the chain don't
	// map back to BestLine (they aren't cards in hand or arsenal) — their damage is counted in
	// Value but not credited per-card here.
	chain := bufs.attackerBuf[:0]
	chain = append(chain, attackers...)
	for _, name := range swungNames {
		for _, w := range weapons {
			if w.Name() == name {
				chain = append(chain, w)
				break
			}
		}
	}
	if len(chain) > 0 {
		winnerOrder := bufs.fillContribWinnerOrder[:len(chain)]
		perCardDmg := bufs.fillContribPerCard[:len(chain)]
		if cap(bufs.fillContribTriggerDmg) < len(chain) {
			bufs.fillContribTriggerDmg = make([]float64, len(chain))
		}
		perCardTrigger := bufs.fillContribTriggerDmg[:len(chain)]
		ctx := &sequenceContext{
			hero:               hero,
			pitched:            pitched,
			deck:               deck,
			bufs:               bufs,
			resourceBudget:     resourceBudget,
			runechantCarryover: runechantCarryover,
		}
		ctx.bestSequence(chain, winnerOrder, perCardDmg, perCardTrigger)
		// Snapshot the winning chain order into TurnSummary so display callers can show cards
		// and weapons in the actual play order (matters for Go-again / trigger chains). Fresh
		// slice — winnerOrder aliases attackBuf storage that the next partition would clobber.
		summary.AttackChain = make([]AttackChainEntry, len(winnerOrder))
		for i := range winnerOrder {
			summary.AttackChain[i] = AttackChainEntry{
				Card:          winnerOrder[i],
				Damage:        perCardDmg[i],
				TriggerDamage: perCardTrigger[i],
			}
		}
		// Map chain-position damage back to BestLine indices by scanning for the first unused
		// Attack-role entry with a matching ID. Duplicate printings played as twin attacks
		// are disambiguated by scan order. Hand-card Contribution bundles Play return plus hero
		// trigger so the per-card stats aggregate the card's total this-turn impact.
		if cap(bufs.fillContribUsed) < total {
			bufs.fillContribUsed = make([]bool, total)
		}
		used := bufs.fillContribUsed[:total]
		for i := range used {
			used[i] = false
		}
		for k, c := range winnerOrder {
			if _, isWeapon := c.(weapon.Weapon); isWeapon {
				continue
			}
			for i := range line {
				if used[i] || line[i].Role != Attack || line[i].Card.ID() != c.ID() {
					continue
				}
				line[i].Contribution = perCardDmg[k] + perCardTrigger[k]
				used[i] = true
				break
			}
		}
	}
}


