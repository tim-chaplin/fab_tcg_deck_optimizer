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
	// Arsenal marks the card placed into the arsenal at end of turn. Contributes nothing to this
	// turn's Value (it scores on the future turn it's played) and carries across the turn boundary
	// as Play.ArsenalCard. At most one card in BestLine takes this role; when the partition
	// enumerator leaves arsenal empty, Best post-hoc promotes one Held card.
	Arsenal
)

// CardAssignment is a single card + the role it took this turn. Hand cards produce one per card;
// an arsenal-in card contributes one with FromArsenal set so a turn fits in one slice.
// Contribution is the per-card credit toward TurnSummary.Value (damage dealt, block share, or
// pitch resource depending on Role), filled by fillContributions once the winner is picked.
type CardAssignment struct {
	Card         card.Card
	Role         Role
	FromArsenal  bool
	Contribution float64
}

// TurnSummary is the result of running Best on a hand: the winning card-role assignments plus
// aggregate metadata about the turn.
type TurnSummary struct {
	// BestLine is the winning partition. Hand cards come first in canonical (post-sort) order;
	// the previous-turn arsenal card, if any, is the last entry with FromArsenal=true. Never
	// mutate — memoized results alias this slice.
	BestLine []CardAssignment
	// AttackChain is the winning chain in play order: attack-role cards from BestLine interleaved
	// with any swung weapons at the positions the solver picked. Each entry carries its Play-time
	// damage (plus hero-trigger damage for cards) so callers can attribute contribution for
	// weapons, which have no BestLine entry. Swung weapons are recoverable by type-asserting
	// AttackChainEntry.Card to weapon.Weapon. Empty when no attacks were played.
	AttackChain []AttackChainEntry
	// Value is the turn's total score (damage dealt + damage prevented). Equals the sum of
	// BestLine[].Contribution plus any weapon-swing damage (weapons aren't in BestLine).
	Value int
	// LeftoverRunechants is the Runechant token count at end of the chosen chain; the caller
	// feeds it back as the next turn's carryover. For partitions with no attacks, equals the
	// carryover the caller passed in.
	LeftoverRunechants int
	// ArsenalCard is the card occupying the arsenal slot at end of turn — either a hand card
	// just arsenaled (role=Arsenal) or a previous-turn arsenal card that stayed. Nil when empty.
	// The caller feeds this back as the next turn's arsenalCardIn.
	ArsenalCard card.Card
}

// AttackChainEntry is a single played attack — a card with role=Attack or a swung weapon —
// carrying the damage it contributed when it resolved in the winning chain. Damage is the Play()
// return; TriggerDamage is the hero's OnCardPlayed contribution (e.g. Viserai creating a
// Runechant) so callers can surface hero attribution on its own line. For BestLine Attack entries
// Damage + TriggerDamage equals CardAssignment.Contribution; weapons live only here.
type AttackChainEntry struct {
	Card          card.Card
	Damage        float64
	TriggerDamage float64
}

// ArsenalIn returns the assignment for the card that started the turn in the arsenal, if any.
// Lets callers treat the arsenal-in card differently from hand cards without scanning BestLine.
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

// formatContribution renders a contribution/damage value for the best-turn printout. Integers
// render bare; fractional values (proportional defense share when multiple blockers split an
// incoming attack) show one decimal place.
func formatContribution(v float64) string {
	if v == float64(int(v)) {
		return fmt.Sprintf("%d", int(v))
	}
	return fmt.Sprintf("%.1f", v)
}

// assignmentName returns the card name, suffixed with " (from arsenal)" when the assignment came
// from the arsenal slot — that tag tells readers why the card isn't in the dealt-hand list the
// optimiser reports alongside.
func assignmentName(a CardAssignment) string {
	if a.FromArsenal {
		return a.Card.Name() + " (from arsenal)"
	}
	return a.Card.Name()
}

// FormatBestLine pairs each card in BestLine with its assigned role for debug output, e.g.
// "Hocus Pocus (Blue): PITCH, Runic Reaping (Red): ATTACK". Compact one-line form; use
// FormatBestTurn for chronological play order.
func FormatBestLine(line []CardAssignment) string {
	parts := make([]string, len(line))
	for i, a := range line {
		parts[i] = assignmentName(a) + ": " + a.Role.String()
	}
	return strings.Join(parts, ", ")
}

// splitPitchesByPhase assigns each pitch card to the defense or attack phase, simulating the
// order FaB prompts them in. Smallest pitches fund the defense bucket until drCost is covered;
// the rest pay for this turn's attacks. Stable on ties so display order is deterministic.
func splitPitchesByPhase(pitched []CardAssignment, drCost int) (defensePitches, attackPitches []CardAssignment) {
	sorted := append([]CardAssignment(nil), pitched...)
	sort.SliceStable(sorted, func(i, j int) bool { return sorted[i].Card.Pitch() < sorted[j].Card.Pitch() })
	covered := 0
	for _, a := range sorted {
		if covered < drCost {
			defensePitches = append(defensePitches, a)
			covered += a.Card.Pitch()
		} else {
			attackPitches = append(attackPitches, a)
		}
	}
	return defensePitches, attackPitches
}

// appendAttackChainLines renders the Attack phase: one numbered line per AttackChain entry in
// solver-chosen play order, with the shared step counter advanced via stepPtr so later sections
// keep numbering contiguous. Non-weapon entries cross-reference BestLine by ID so arsenal-played
// cards get a "(from arsenal)" tag; weapons skip the match since they have no BestLine entry. A
// non-zero TriggerDamage adds a trailing " (+M hero trigger)" so the attribution is visible
// instead of silently folded into the card's own damage number.
func appendAttackChainLines(lines []string, t TurnSummary, stepPtr *int) []string {
	used := make([]bool, len(t.BestLine))
	appendAttack := func(label, cardName string, e AttackChainEntry) {
		*stepPtr++
		line := fmt.Sprintf("  %d. %s: %s (+%s)", *stepPtr, cardName, label, formatContribution(e.Damage))
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
		// Match the first unused Attack-role BestLine entry by ID so we can detect FromArsenal.
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
	return lines
}

// FormatBestTurn renders a TurnSummary as a numbered play-order list, one card per line,
// matching the actual FaB turn sequence:
//
//  1. Defense-phase pitches (paying for Defense Reactions)
//  2. Plain blocks
//  3. Defense Reactions
//  4. Attack-phase pitches (paying for this turn's played cards)
//  5. Attack chain — played cards and swung weapons in the order the solver picked
//
// Held / Arsenal cards are summarized on trailing lines so the reader sees what's carrying over.
//
// Pitch-phase assignment uses a greedy split for display: smallest pitches first fund the defense
// pool until drCost is covered, the rest fund attack. The solver already validated some legal
// split exists; this picks one deterministically.
func FormatBestTurn(t TurnSummary) string {
	// Partition BestLine into role buckets; pitch cards pool and are split by phase below.
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

	defensePitches, attackPitches := splitPitchesByPhase(pitched, drCost)

	var lines []string
	step := 0
	nextStep := func() int { step++; return step }
	// Pitch lines don't show a damage contribution: pitches generate resource, not damage —
	// showing "+3" next to a pitch would double-count against the turn's Value.
	appendPitch := func(a CardAssignment, roleLabel string) {
		lines = append(lines, fmt.Sprintf("  %d. %s: %s", nextStep(), assignmentName(a), roleLabel))
	}
	// Defense lines contribute damage-prevention (block share + DR Play return), so the tag is
	// shown.
	appendDefense := func(a CardAssignment, roleLabel string) {
		lines = append(lines, fmt.Sprintf("  %d. %s: %s (+%s prevented)", nextStep(), assignmentName(a), roleLabel, formatContribution(a.Contribution)))
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
	lines = appendAttackChainLines(lines, t, &step)

	// Held / Arsenal footer: unplayed cards, outside the numbered sequence, shown so the reader
	// sees the whole turn disposition.
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
// for incomingDamage on their next turn. Equipped weapons may be swung for their Cost if
// resources allow.
//
// Cards partition into five roles:
//   - Pitch: contributes its Pitch value as resource paying for a played card this turn or a
//     Defense Reaction on the opponent's turn.
//   - Attack: consumes Cost on our turn; resolved by calling Card.Play in an order the optimizer
//     chooses. TurnState effects carry forward to later attacks in the same sequence.
//   - Defend: contributes Defense to damage prevented (capped at incomingDamage; excess wasted).
//     Plain blocking is free; Defense Reactions must pay Cost and contribute any Play() damage.
//   - Held: stays in hand for next turn. Contributes nothing this turn.
//   - Arsenal: moves into the arsenal slot at end of turn, or for an arsenal-in card, stays.
//     Contributes nothing this turn.
//
// Pitch resources split across two phases because resources don't carry between turns: attack
// pitches pay for this turn's played cards, defense pitches pay for Defense Reactions on the
// opponent's turn. A card can only be pitched while some unpaid card is on the stack in the
// matching phase, so a hand with no plays in a phase has no Pitch-role cards in that phase; any
// card that can't be legally pitched becomes Held.
//
// Results are memoized on (heroID, sorted weapon IDs, sorted card IDs, incomingDamage,
// runechantCarryover, arsenal-in ID) so repeat evaluations short-circuit. The hand is sorted in
// place into canonical order first; BestLine's hand entries align with that post-sort order.
// Every card in the hand must be registered in package cards or Best panics.
//
// runechantCarryover is the Runechant token count carrying in from the previous turn.
// TurnSummary.LeftoverRunechants is the count at end of the chosen chain; feed it back as the
// next turn's carryover.
//
// arsenalCardIn is the card sitting in the arsenal slot at start of turn (nil if empty). The
// enumerator pulls it in as an extra CardAssignment with restricted role options — Arsenal
// (stay), Attack (any non-DR card), or Defend (Defense Reactions only). Never Pitch or Held. A
// hand card may also take the Arsenal role so long as at most one BestLine entry ends up there.
func Best(hero hero.Hero, weapons []weapon.Weapon, hand []card.Card, incomingDamage int, deck []card.Card, runechantCarryover int, arsenalCardIn card.Card) TurnSummary {
	return sharedEvaluator.Best(hero, weapons, hand, incomingDamage, deck, runechantCarryover, arsenalCardIn)
}

// Best is the method form of the package-level Best: same semantics, uses this Evaluator's
// scratch buffers so concurrent goroutines can each hold their own.
func (e *Evaluator) Best(hero hero.Hero, weapons []weapon.Weapon, hand []card.Card, incomingDamage int, deck []card.Card, runechantCarryover int, arsenalCardIn card.Card) TurnSummary {
	// IDs go into a fixed-size stack array to avoid a per-call slice alloc. Hand size is capped
	// at 8 (matches memoKey.cardIDs); larger hands panic out of the inner loops.
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

// sortHandByID sorts the first n entries of hand and ids in parallel by ascending id, in place.
// Insertion sort — for n ≤ 8 this beats sort.Sort and avoids boxing slices through sort.Interface.
// Canonicalizing the hand order is what lets the memo key collapse permutations onto one entry.
func sortHandByID(hand []card.Card, ids []card.ID, n int) {
	for i := 1; i < n; i++ {
		for j := i; j > 0 && ids[j-1] > ids[j]; j-- {
			ids[j-1], ids[j] = ids[j], ids[j-1]
			hand[j-1], hand[j] = hand[j], hand[j-1]
		}
	}
}

// memoKey is the comparable map key for the shared memo. Hand size is capped at 8 cards. Hero
// ID + weapon IDs live in the key so Evaluators with different (hero, weapons) tuples coexist in
// the memo without a scope wipe — distinct scopes just produce distinct keys. The uint16 hero.ID
// keeps the whole key a fixed-size integer struct; no string hashing per lookup.
type memoKey struct {
	heroID             hero.ID
	weaponIDs          [2]card.ID
	cardIDs            [8]card.ID
	cardCount          uint8
	incoming           int
	runechantCarryover int
	// arsenalInID is card.Invalid when the slot is empty, otherwise the ID of the starting
	// arsenal card — different arsenal-ins give distinct cache entries.
	arsenalInID card.ID
}

// Evaluator owns the per-goroutine mutable state hand.Best threads through an evaluation: a
// scratch-buffer cache keyed by (handSize, weapons). The memo cache is shared across all
// Evaluators so every worker benefits from cached hands — only the mutated scratch buffers
// must be per-goroutine. A long-lived Evaluator avoids reallocating ~20 scratch slices per call.
type Evaluator struct {
	// bufs holds the pre-allocated scratch slices for bestUncached / bestAttackWithWeapons /
	// bestSequence. Keyed by (handSize, weaponCount, weaponIDs); recreated when any differ so the
	// scratch sizing stays correct.
	bufs            *attackBufs
	bufsHandSize    int
	bufsWeaponIDs   [2]card.ID
	bufsWeaponCount int
	bufsValid       bool
}

// NewEvaluator returns a fresh Evaluator with empty scratch buffers. One Evaluator per goroutine
// is safe for concurrent use and still shares the global memo.
func NewEvaluator() *Evaluator {
	return &Evaluator{}
}

// sharedEvaluator backs the package-level Best — single-threaded callers don't need to construct
// their own. Parallel callers create one per goroutine for bufs but all share the global memo.
var sharedEvaluator = NewEvaluator()

// memo caches canonical-order TurnSummary results keyed by memoKey. Shared across all Evaluators
// and goroutines; protected by memoMu. Hero + weapon IDs live in memoKey so distinct (hero,
// weapons) scopes coexist without a wipe step.
var (
	memo   = map[memoKey]TurnSummary{}
	memoMu sync.RWMutex
)

// ClearMemo drops every cached TurnSummary. Callers (iterate-mode, benchmarks) use this to cap
// memo growth across unrelated runs; cross-run hit rate is near zero so nothing of value is lost.
func ClearMemo() {
	memoMu.Lock()
	clear(memo)
	memoMu.Unlock()
}

// MemoLen returns the current number of cached entries, for diagnostic logging.
func MemoLen() int {
	memoMu.RLock()
	n := len(memo)
	memoMu.RUnlock()
	return n
}

// makeMemoKey builds a comparable memo key. The hand must already be sorted by card ID.
// sortedIDs is a pointer to the caller's stack [8]card.ID to avoid a slice-header escape. Weapon
// IDs are sorted into the two fixed slots so loadouts in any order hash to the same key.
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
// this hoisted to a per-attacker lookup, the hot inner loop skips Types / Cost / GoAgain
// interface dispatch; the one meta build amortises across all N! permutations.
type attackerMeta struct {
	types            card.TypeSet
	cost             int
	printedCost      int // PrintedCost for DiscountPerRunechant cards; 0 otherwise
	isDiscount       bool
	baseGoAgain      bool // Card.GoAgain() — combined with PlayedCard.GrantedGoAgain at read time
	isAttackOrWeapon bool
}

// cardMetaCache / cardMetaReady are shared, read-only-after-init card metadata tables. Populated
// lazily via cardMetaSlowPath on first encounter, then read from all goroutines without sync.
// Sized for the full uint16 ID space so lookups are plain bounds-checked reads (~2 MB total).
const cardMetaCacheSize = 1 << 16

var (
	cardMetaCache [cardMetaCacheSize]attackerMeta
	cardMetaReady [cardMetaCacheSize]uint32 // written once (atomically) per ID; 0 = unready, 1 = ready
	cardMetaMu    sync.Mutex
)

// attackerMetaFor returns cached metadata for c, populating on first encounter. Safe from
// multiple goroutines: the first writer per ID holds the mutex, later readers see the ready flag
// set with a release barrier and read the immutable meta entry directly.
func attackerMetaFor(c card.Card) attackerMeta {
	id := c.ID()
	if atomic.LoadUint32(&cardMetaReady[id]) == 1 {
		return cardMetaCache[id]
	}
	return cardMetaSlowPath(c, id)
}

// attackerMetaPtrFor is the pointer-returning counterpart of attackerMetaFor: it hands back a
// direct pointer into the global cache so permutation swaps move 8 bytes instead of a full
// attackerMeta struct. The target is read-only after initialisation.
func attackerMetaPtrFor(c card.Card) *attackerMeta {
	id := c.ID()
	if atomic.LoadUint32(&cardMetaReady[id]) == 1 {
		return &cardMetaCache[id]
	}
	cardMetaSlowPath(c, id)
	return &cardMetaCache[id]
}

// cardMetaSlowPath populates the cache entry under cardMetaMu and returns the computed meta.
func cardMetaSlowPath(c card.Card, id card.ID) attackerMeta {
	cardMetaMu.Lock()
	defer cardMetaMu.Unlock()
	// Re-check under lock: another goroutine may have populated between the atomic load and here.
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
// playSequence) and the partition loop in bestUncached. Allocated once and cached on the
// Evaluator so a deck eval reuses them across every partition, mask, and permutation.
type attackBufs struct {
	perm           []card.Card
	pcBuf          []card.PlayedCard
	ptrBuf         []*card.PlayedCard
	cardsPlayedBuf []card.Card
	state          *card.TurnState
	attackerBuf    []card.Card // for bestAttackWithWeapons mask iteration
	// Pre-computed per-mask weapon data. Indexed by bitmask (0 to 2^len(weapons)-1):
	// weaponCosts[mask] is total Cost; weaponNames[mask] is the pre-built []string of names.
	weaponCosts []int
	weaponNames [][]string
	// permMeta parallels perm: each entry points into the global cardMetaCache so playSequence's
	// inner loop skips interface dispatch on Types / Cost / GoAgain / DiscountPerRunechant.
	// Pointer-valued so bestSequence's permutation swaps move 8 bytes instead of a full struct.
	permMeta []*attackerMeta
	// Partition-loop buffers, consumed by bestUncached. Sized handSize+1 to cover the optional
	// arsenal-in slot the enumerator treats as index n. defendPrintedVals holds PrintedCost for
	// DiscountPerRunechant defense reactions (post-attack discount re-pricing); non-discount
	// cards get 0. isDRBuf caches TypeDefenseReaction membership to skip Types().Has calls.
	rolesBuf          []Role
	pitchVals         []int
	costVals          []int
	defendCostVals    []int
	defendPrintedVals []int
	defenseVals       []int
	isDRBuf           []bool
	// pitchedValsScratch backs the per-leaf "pitched values" slice the feasibility check reads.
	// Re-sliced to [:0] at the start of every leaf to eliminate a make([]int, …) per leaf.
	pitchedValsScratch []int
	pitchedBuf         []card.Card
	attackersBuf       []card.Card
	defendersBuf       []card.Card
	// perCardScratch is sized maxAttackers (handSize + weaponCount). Written by playSequence only
	// when the caller passes a non-nil perCardOut; bestSequence snapshots the winning
	// permutation's per-card damage from here into the caller's output buffer. The partition-loop
	// hot path passes nil and never touches this slice.
	perCardScratch []float64
	// perCardTriggerScratch parallels perCardScratch for hero-trigger damage (OnCardPlayed
	// return). Only written when the caller tracks.
	perCardTriggerScratch []float64
	// fillContribWinnerOrder / fillContribPerCard are output buffers for bestSequence during
	// fillContributions's tracked replay. Kept on attackBufs so each Best call reuses the slab.
	fillContribWinnerOrder []card.Card
	fillContribPerCard     []float64
	fillContribTriggerDmg  []float64
	// fillContribUsed marks hand indices already assigned during chain→hand mapping. Sized
	// handSize; reset with clear before each fillContributions pass.
	fillContribUsed []bool
}

func newAttackBufs(handSize, weaponCount int, weapons []weapon.Weapon) *attackBufs {
	// +1 reserves a slot for the arsenal-in card, which joins attackers or defenders when the
	// enumerator plays it from arsenal. Without it, all-attack hands + arsenal would overflow
	// attackerBuf in bestAttackWithWeapons.
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
	pcBuf := make([]card.PlayedCard, maxAttackers)
	ptrBuf := make([]*card.PlayedCard, maxAttackers)
	// Wire the ptrBuf entries to their pcBuf slots once — the mapping is stable across every
	// permutation so playSequenceWithMeta doesn't need to rewrite it per call.
	for i := range pcBuf {
		ptrBuf[i] = &pcBuf[i]
	}
	return &attackBufs{
		perm:                   make([]card.Card, maxAttackers),
		permMeta:               make([]*attackerMeta, maxAttackers),
		pcBuf:                  pcBuf,
		ptrBuf:                 ptrBuf,
		cardsPlayedBuf:         make([]card.Card, 0, maxAttackers),
		state:                  &card.TurnState{},
		attackerBuf:            make([]card.Card, maxAttackers),
		weaponCosts:            weaponCosts,
		weaponNames:            weaponNames,
		rolesBuf:               make([]Role, handSize+1),
		pitchVals:              make([]int, handSize+1),
		costVals:               make([]int, handSize+1),
		defendCostVals:         make([]int, handSize+1),
		defendPrintedVals:      make([]int, handSize+1),
		defenseVals:            make([]int, handSize+1),
		isDRBuf:                make([]bool, handSize+1),
		pitchedValsScratch:     make([]int, 0, handSize+1),
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

// getAttackBufs returns this Evaluator's scratch-buffer set, rebuilding when (handSize, weapons)
// changes. Single-slot per Evaluator: iterate runs tens of thousands of same-shape hands, so a
// slot outperforms a keyed pool for this workload.
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
	// restricted role menu (Arsenal / Attack / Defend), so everything about it is decided inside
	// the enumeration. totalN is the effective size of BestLine.
	totalN := n
	if arsenalCardIn != nil {
		totalN = n + 1
	}

	// Seed best.LeftoverRunechants with the carryover: partitions with no attacks don't reduce
	// it, so carryover is the baseline to beat. BestLine starts with every hand card Held and
	// the arsenal-in card (if any) staying in the slot, so a hand with no Value-adding partition
	// still reports sensible "nothing played, nothing pitched" assignments.
	best := TurnSummary{BestLine: make([]CardAssignment, totalN), LeftoverRunechants: runechantCarryover}
	// bestSwung holds the winning partition's swung weapon names so fillContributions can rebuild
	// the chain it runs bestSequence over. Lives outside TurnSummary since weapons are
	// recoverable from AttackChain once fillContributions finishes.
	var bestSwung []string
	// bestHasHeld tracks whether the current best has at least one Held hand card — lets
	// beatsBest distinguish "arsenal will be occupied post-hoc" from "arsenal will be empty."
	// Seeded true when the hand is non-empty: the initial best puts every hand card into Held,
	// so a post-hoc promotion would fill arsenal. Candidates need both a Value/leftover tie and
	// some way to end with arsenal occupied to displace it.
	bestHasHeld := n > 0
	for i := 0; i < n; i++ {
		best.BestLine[i] = CardAssignment{Card: hand[i], Role: Held}
	}
	if arsenalCardIn != nil {
		best.BestLine[n] = CardAssignment{Card: arsenalCardIn, Role: Arsenal, FromArsenal: true}
		best.ArsenalCard = arsenalCardIn
	}

	// bufs is the pooled scratch space for this deck evaluation. Partition scratch is sized
	// handSize+1, big enough for totalN when an arsenal-in card inflates the effective hand.
	// Each field is re-sliced and rewritten below, so carry-over from prior calls can't leak.
	bufs := e.getAttackBufs(n, weapons)
	rolesBuf := bufs.rolesBuf[:totalN]
	pvals := bufs.pitchVals[:totalN]
	cvals := bufs.costVals[:totalN]
	dvals := bufs.defenseVals[:totalN]
	dCostVals := bufs.defendCostVals[:totalN]
	dPrintedVals := bufs.defendPrintedVals[:totalN]
	isDR := bufs.isDRBuf[:totalN]

	// Pre-compute per-card pitch / cost / defense values so the recurse doesn't re-invoke the
	// card-method interface calls per leaf. defendCostVals holds Cost only for Defense Reactions
	// (non-reactions block free). defendPrintedVals holds PrintedCost for DiscountPerRunechant
	// defenders (used by the post-attack discount re-pricing check) and zero otherwise. Both are
	// populated conditionally, so clear prior-call residue from the pooled scratch first.
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
	// pitchedValsScratch is reused across leaves to avoid per-leaf allocation.
	pitchedValsScratch := bufs.pitchedValsScratch[:0]
	var recurse func(i, pitchSum, costSum, defenseSum, defenderCostSum int)
	recurse = func(i, pitchSum, costSum, defenseSum, defenderCostSum int) {
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
			// Group roles into played / pitched / defending buckets. Iterates the hand (size n),
			// then layers in the arsenal slot (index n) based on its assigned role. Arsenal-role
			// cards contribute nothing this turn whether they came from hand or the slot.
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
			attackDealt, leftoverRunechants, residualBudget, swung := bestAttackWithWeapons(hero, weapons, a, p, deck, bufs, pitchSum, costSum, defenderCostSum, runechantCarryover, incomingDamage, defenseSum)

			// DiscountPerRunechant defense reactions reserved 0 in defenderCostSum (their Cost()
			// is the fully-discounted minimum). Re-price them now that the attack chain has
			// resolved with `leftoverRunechants` tokens available. Arsenal-in defenders aren't
			// checked: no DiscountPerRunechant cards currently want to live in arsenal.
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
			arsenalCard := findArsenalCard(rolesBuf, arsenalCardIn, n)
			// Hand cards never take Arsenal role during enumeration, so arsenalCard is only set
			// when arsenal-in stayed; post-hoc promotion potential is tracked via hasHeld.
			hasHeld := false
			for j := 0; j < n; j++ {
				if rolesBuf[j] == Held {
					hasHeld = true
					break
				}
			}
			willOccupy := arsenalCard != nil || hasHeld
			bestWillOccupy := best.ArsenalCard != nil || bestHasHeld
			if !beatsBest(v, leftoverRunechants, willOccupy, best, bestWillOccupy) {
				return
			}
			best.Value = v
			bestSwung = swung
			best.LeftoverRunechants = leftoverRunechants
			best.ArsenalCard = arsenalCard
			bestHasHeld = hasHeld
			// Write the winning roles into BestLine. Cards and FromArsenal flags were populated
			// at construction; only Role varies. Contribution is cleared here and filled by
			// fillContributions below for the winning line.
			for j := 0; j < totalN; j++ {
				best.BestLine[j].Role = rolesBuf[j]
				best.BestLine[j].Contribution = 0
			}
			return
		}
		isArsenalSlot := i == n && arsenalCardIn != nil
		// Hand cards can't take Arsenal role (post-hoc promotion handles that). Cap the range at
		// Held for hand slots to skip the roleAllowed-rejection work for Arsenal.
		maxRole := Held
		if isArsenalSlot {
			maxRole = Arsenal
		}
		for r := Role(0); r <= maxRole; r++ {
			if !roleAllowed(r, isArsenalSlot, isDR[i]) {
				continue
			}
			rolesBuf[i] = r
			switch r {
			case Pitch:
				recurse(i+1, pitchSum+pvals[i], costSum, defenseSum, defenderCostSum)
			case Attack:
				recurse(i+1, pitchSum, costSum+cvals[i], defenseSum, defenderCostSum)
			case Defend:
				recurse(i+1, pitchSum, costSum+dCostVals[i], defenseSum+dvals[i], defenderCostSum+dCostVals[i])
			case Held, Arsenal:
				recurse(i+1, pitchSum, costSum, defenseSum, defenderCostSum)
			}
		}
	}
	recurse(0, 0, 0, 0, 0)
	// Once per Best call, on the winning line only, attribute per-card contribution.
	if len(best.BestLine) > 0 {
		fillContributions(&best, hero, weapons, bestSwung, deck, bufs, incomingDamage, runechantCarryover)
	}
	// If the arsenal slot is empty after enumeration, promote a Held card. The pick is
	// deterministic per-hand (hashed from sorted card IDs) so the memo stays consistent, but
	// spreads across Held positions across different hands — avoiding the lowest-ID bias that
	// picking BestLine[0] would introduce under sort-by-ID canonicalisation.
	if best.ArsenalCard == nil {
		promoteRandomHeldToArsenal(&best, hand, n, arsenalCardIn)
	}
	return best
}

// promoteRandomHeldToArsenal picks one Held hand card in best.BestLine and flips its role to
// Arsenal. Selection hashes the sorted hand IDs (plus arsenal-in ID) modulo the Held count:
//   - Same hand → same promotion (memo-safe).
//   - Different hands spread across Held positions (no systematic lowest-ID preference).
func promoteRandomHeldToArsenal(best *TurnSummary, hand []card.Card, n int, arsenalCardIn card.Card) {
	// Collect Held indices (hand slots only; arsenal-in can't be Held by construction).
	heldIndices := make([]int, 0, n)
	for i := 0; i < n; i++ {
		if best.BestLine[i].Role == Held {
			heldIndices = append(heldIndices, i)
		}
	}
	if len(heldIndices) == 0 {
		return
	}
	// FNV-1a-flavoured hash over the sorted hand IDs + arsenal-in ID. Just needs to spread
	// across bucket counts 1..n.
	var h uint64 = 1469598103934665603 // FNV offset basis
	for _, c := range hand {
		h ^= uint64(c.ID())
		h *= 1099511628211 // FNV prime
	}
	if arsenalCardIn != nil {
		h ^= uint64(arsenalCardIn.ID())
		h *= 1099511628211
	}
	pick := heldIndices[h%uint64(len(heldIndices))]
	best.BestLine[pick].Role = Arsenal
	best.ArsenalCard = best.BestLine[pick].Card
}

// canCoverPhasesAllUsed decides whether every pitched value can be split between the attack
// phase (covering attackCost) and the defense phase (covering drCost) while respecting FaB's
// pitch-timing rule: a card can only be pitched while some unpaid card is on the stack. So every
// Pitch-role card must pay for something — any "extra" card would have to be Held instead.
//
// Per-phase legality uses a sufficient condition: sum == phaseCost is trivially legal; sum >
// phaseCost needs the excess absorbable by one over-paying pitch, i.e. max(pool) > excess.
// Partitions needing multiple pitches each pushing past a cost's ceiling are rejected.
//
// With both phases positive we enumerate every non-empty, non-full attack-pool mask (2^k - 2 for
// k pitched cards) and take the first that satisfies both phase checks. k is bounded by hand
// size so this stays cheap relative to the outer 4^n partition search.
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

// phaseLegal returns true iff the pitch values selected by subsetMask can legally cover
// phaseCost. sum < phaseCost can't pay; sum == phaseCost is exact (legal); sum > phaseCost needs
// one pitch large enough to absorb the excess as a single over-paying final pitch.
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

// groupByRoleInto appends hand cards into caller-provided pitched/attackers/defenders slices
// (passed pre-reset to length 0) to avoid per-partition heap allocation.
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

// groupPitchAttack is the reaction-free leaf's grouping step: skips the defenders bucket (only
// needed for Defense-Reaction-Play dispatch, which this path doesn't run).
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
// of the pitched values between attack and defense phases. Called at every partition leaf —
// phase-legality is the cheap screen that keeps bestAttackWithWeapons off infeasible partitions.
func anyMaskFeasible(pitchedVals []int, attackCardCost, drCost int, weaponCosts []int, weaponCount int) bool {
	masks := 1 << weaponCount
	for mask := 0; mask < masks; mask++ {
		if canCoverPhasesAllUsed(pitchedVals, attackCardCost+weaponCosts[mask], drCost) {
			return true
		}
	}
	return false
}

// findArsenalCard returns the arsenal-in card when it stays in the arsenal slot, nil otherwise.
// Hand cards never take Arsenal role during enumeration (post-hoc promotion handles that), so
// the only slot that can be Arsenal is the arsenal-in slot at index n.
func findArsenalCard(rolesBuf []Role, arsenalCardIn card.Card, n int) card.Card {
	if arsenalCardIn != nil && rolesBuf[n] == Arsenal {
		return arsenalCardIn
	}
	return nil
}

// beatsBest decides whether a candidate partition displaces the current best. Tiebreak order:
// higher Value, then more leftover runechants (they convert to future arcane damage), then
// preferring a partition that ends with the arsenal slot occupied (saves a hand slot next
// refill). "Occupied" covers both an arsenal-in card that stayed and a Held hand card slated
// for post-hoc promotion.
func beatsBest(v, leftoverRunechants int, willOccupyArsenal bool, best TurnSummary, bestWillOccupyArsenal bool) bool {
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
	return willOccupyArsenal && !bestWillOccupyArsenal
}

// roleAllowed decides whether the partition enumerator may assign role r to the current card.
// The arsenal-in slot may only take Arsenal (stay), Attack (any non-DR card — non-attack actions
// play fine from arsenal on your turn), or Defend (Defense Reactions only — plain-blocking from
// arsenal isn't legal). Hand cards take any role except Attack for Defense Reactions (DRs only
// fire on the opponent's turn); their role loop caps at Held, so the "which Held card gets
// arsenaled" choice happens post-hoc and doesn't bias toward low-ID slots.
func roleAllowed(r Role, isArsenalSlot, isDefenseReaction bool) bool {
	if isArsenalSlot {
		switch r {
		case Pitch, Held:
			return false
		case Attack:
			return !isDefenseReaction
		case Defend:
			return isDefenseReaction
		}
		return true // Arsenal is always allowed on the arsenal-in slot.
	}
	return !(r == Attack && isDefenseReaction)
}

// defenseReactionDamage runs Play() for every Defense Reaction in defenders and sums the damage
// they deal back to the attacker (e.g. Weeping Battleground's 1 arcane on banish). Played in
// isolation — no attack ordering; TurnState carries only Pitched/Deck so effects that read
// "what was pitched" work. Uncapped: this damage is dealt, not prevented.
//
// state is caller-provided (from attackBufs) and reset per call. Reusing the pointer keeps the
// state on the heap buffer rather than escaping a fresh stack value per partition leaf.
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

// bestAttackWithWeapons enumerates every subset of weapons to swing alongside attackers and
// returns the max damage over all affordable masks, the runechant leftover, the residual chain
// budget (pitch not consumed by the winning line — the caller re-checks DiscountPerRunechant
// defense affordability against it), and the swung weapons in input order.
//
// resourceBudget is pitchSum - defenderCostSum; the chain further deducts each attacker's
// effective cost (DiscountPerRunechant: max(0, PrintedCost - runechants at play-time); others:
// Cost()) and rejects orderings that run negative.
func bestAttackWithWeapons(hero hero.Hero, weapons []weapon.Weapon, attackers, pitched, deck []card.Card, bufs *attackBufs, pitchSum, costSum, defenderCostSum, runechantCarryover, incomingDamage, blockTotal int) (int, int, int, []string) {
	// Build the sequence-eval context once per partition leaf: resourceBudget, runechantCarryover,
	// and pitched/deck refs are constant across the weapon-mask loop.
	ctx := &sequenceContext{
		hero:               hero,
		pitched:            pitched,
		deck:               deck,
		bufs:               bufs,
		resourceBudget:     pitchSum - defenderCostSum,
		runechantCarryover: runechantCarryover,
		incomingDamage:     incomingDamage,
		blockTotal:         blockTotal,
	}
	// Hoist the leaf-constant TurnState fields out of the per-permutation reset in
	// playSequenceWithMeta. Pitched, Deck, IncomingDamage, BlockTotal don't change across a
	// partition's permutations; setting them once per ctx saves four stores per playSequence call.
	ctx.seedState()
	best := 0
	bestLeftoverRunechants := runechantCarryover
	bestResidualBudget := ctx.resourceBudget
	var bestSwung []string
	// Reuse the shared attacker buffer across mask iterations.
	copy(bufs.attackerBuf, attackers)
	// pitchedVals feeds the per-mask phase-feasibility check so each mask's weapon cost folds
	// into the attack-phase total. Locally allocated to keep this call re-entrant.
	attackCardCost := costSum - defenderCostSum
	drCost := defenderCostSum
	pitchedVals := make([]int, len(pitched))
	for i, c := range pitched {
		pitchedVals[i] = c.Pitch()
	}
	for mask := 0; mask < 1<<len(weapons); mask++ {
		// bufs.weaponCosts[mask] is the pre-summed Cost, avoiding per-weapon interface dispatch.
		// The feasibility check ensures every Pitch-role card can legally pay for something.
		if !canCoverPhasesAllUsed(pitchedVals, attackCardCost+bufs.weaponCosts[mask], drCost) {
			continue
		}
		allAttackers := bufs.attackerBuf[:len(attackers)]
		for i, w := range weapons {
			if mask&(1<<i) != 0 {
				allAttackers = append(allAttackers, w)
			}
		}
		// Every sequence card (attackers AND swung weapons) deducts its own effective cost inside
		// playSequence, so don't pre-deduct weapon cost here.
		dealt, leftoverRunechants, residualBudget := ctx.bestSequence(allAttackers, nil, nil, nil)
		// Tiebreaks: prefer more leftover runechants, then more residual budget — both are extra
		// slack that can enable discount defense reactions.
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

// sequenceContext carries the stable per-partition-leaf environment: hero (for OnCardPlayed
// triggers), pitched / deck refs for Card.Play, shared scratch buffers, and the numeric budgets
// that persist across permutation and mask iterations. Built once per leaf so the hot inner
// calls (playSequence, bestSequence) shrink to their varying inputs and tracking outputs.
type sequenceContext struct {
	hero               hero.Hero
	pitched, deck      []card.Card
	bufs               *attackBufs
	resourceBudget     int
	runechantCarryover int
	incomingDamage     int
	blockTotal         int
}

// seedState writes the TurnState fields that are constant across a partition's permutations
// (pitched / deck references, incoming damage, block total). Called once per ctx so
// playSequenceWithMeta's hot per-permutation reset can skip them.
func (ctx *sequenceContext) seedState() {
	s := ctx.bufs.state
	s.Pitched = ctx.pitched
	s.Deck = ctx.deck
	s.IncomingDamage = ctx.incomingDamage
	s.BlockTotal = ctx.blockTotal
}

// bestSequence tries every ordering of attackers and returns the max total damage plus the
// runechant count at the end of the winning permutation. Between each card's Play() and its
// append to CardsPlayed, the hero's OnCardPlayed hook fires so triggered abilities contribute.
//
// Uses Heap's algorithm (iterative) — no closure/callback alloc, no recursive call per perm.
//
// When winnerOrderOut is non-nil (len >= len(attackers)) the winning permutation is copied into
// it. perCardOut / perCardTriggerOut (same size rule) receive the winning line's per-card Play
// damage and hero-trigger damage. fillContributions uses these; the partition-loop caller
// passes nil for all three so the permutation search stays allocation-free.
func (ctx *sequenceContext) bestSequence(attackers, winnerOrderOut []card.Card, perCardOut, perCardTriggerOut []float64) (int, int, int) {
	n := len(attackers)
	if n == 0 {
		return 0, ctx.runechantCarryover, ctx.resourceBudget
	}
	perm := ctx.bufs.perm[:n]
	permMeta := ctx.bufs.permMeta[:n]
	copy(perm, attackers)
	for idx, c := range attackers {
		permMeta[idx] = attackerMetaPtrFor(c)
	}

	// Scratch buffers are playSequence's per-card outputs, overwritten every permutation. On a
	// new winner we copy them into the caller's perCardOut / perCardTriggerOut. Only populated
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
		// Tiebreaks: more leftover runechants, then more residual budget — both are slack that
		// can enable discount defense reactions.
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
	// Heap's algorithm, iterative: c[] counts how many times each stack frame has iterated.
	// perm and permMeta swap together so playSequence sees meta aligned with the permutation.
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

// playSequence plays `order` as a sequence of cards, reusing ctx.bufs' pooled buffers. Buffers
// are mutated in place; the caller must not read them concurrently.
//
// When perCardOut is non-nil (len >= n) each entry is the card's Play return for that position;
// perCardTriggerOut (same size rule) receives the hero's OnCardPlayed return for that position.
// The hot partition-loop callers pass nil for both; fillContributions's replay passes real slices.
//
// Runechant flow:
//   - state.Runechants starts at ctx.runechantCarryover (tokens from the previous turn).
//   - Each card's Play / hero OnCardPlayed may call CreateRunechants, incrementing the count AND
//     returning n damage — tokens are credited exactly once, at creation.
//   - After each Attack- or Weapon-typed card's Play+OnCardPlayed resolve, all current tokens
//     fire and are destroyed: state.Runechants is zeroed but damage is NOT re-added (that would
//     double-count tokens credited at creation).
//   - At end of the sequence, state.Runechants is the leftover count carrying into next turn.
//
// Resource flow:
//   - ctx.resourceBudget is the starting pool; each card deducts its effective cost. For cards
//     implementing DiscountPerRunechant, effective cost is max(0, PrintedCost -
//     state.Runechants) at play time; otherwise Cost(). Negative remaining budget returns
//     legal=false (the caller treats that ordering as zero damage).
//
// Populates permMeta from order and then calls playSequenceWithMeta. The hot path (bestSequence)
// builds meta once and calls playSequenceWithMeta directly so interface dispatch for scalar
// attributes amortises across the N! permutations it evaluates.
func (ctx *sequenceContext) playSequence(order []card.Card, perCardOut, perCardTriggerOut []float64) (damage int, leftoverRunechants int, residualBudget int, legal bool) {
	ctx.seedState()
	meta := ctx.bufs.permMeta[:len(order)]
	for i, c := range order {
		meta[i] = attackerMetaPtrFor(c)
	}
	return ctx.playSequenceWithMeta(order, perCardOut, perCardTriggerOut)
}

// playSequenceWithMeta runs a specific attacker ordering. Assumes ctx.bufs.permMeta[:len(order)]
// holds metadata aligned with order; the caller keeps them in lockstep (bestSequence swaps meta
// whenever it swaps perm).
func (ctx *sequenceContext) playSequenceWithMeta(order []card.Card, perCardOut, perCardTriggerOut []float64) (damage int, leftoverRunechants int, residualBudget int, legal bool) {
	n := len(order)
	pcBuf := ctx.bufs.pcBuf
	ptrBuf := ctx.bufs.ptrBuf
	meta := ctx.bufs.permMeta[:n]
	// ptrBuf entries point at the matching pcBuf slots permanently (wired once in newAttackBufs),
	// so only the per-permutation Card and the zeroed GrantedGoAgain need refreshing here.
	for i, c := range order {
		pcBuf[i] = card.PlayedCard{Card: c}
		if perCardOut != nil {
			perCardOut[i] = 0
		}
		if perCardTriggerOut != nil {
			perCardTriggerOut[i] = 0
		}
	}
	played := ptrBuf[:n]
	state := ctx.bufs.state
	// Reset only the fields playSequence mutates. Pitched and Deck are stable across permutations
	// in a partition leaf; CardsRemaining and Self are rewritten per-card in the loop. A
	// full-struct replace would memcpy every field (including big slice headers) and profiled at
	// ~850ms/call; skipping caller-unread fields cuts most of that.
	// Pitched / Deck / IncomingDamage / BlockTotal are seeded once per ctx (see seedState); cards
	// don't mutate them, so we skip the per-permutation reset.
	state.CardsPlayed = ctx.bufs.cardsPlayedBuf[:0]
	state.Runechants = ctx.runechantCarryover
	state.DelayedRunechants = 0
	state.ArcaneDamageDealt = false
	state.AuraCreated = false
	state.Overpower = false
	resources := ctx.resourceBudget
	for i, pc := range played {
		m := meta[i]
		// Effective cost: discount cards drop by runechant count (floored at 0); others pay Cost.
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

		// If this card is an attack or weapon and any Runechant is live, those tokens fire on
		// its damage step. Set ArcaneDamageDealt now — before Play and OnCardPlayed — so Play
		// effects that read "if you've dealt arcane damage this turn" see the flag for same-hand
		// triggers. Cards that deal arcane damage via their Play text flip the flag themselves.
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

		// Attacks and weapon swings consume all runechants in play. Damage isn't re-added: each
		// token was credited +1 at creation time, so this is pure state cleanup.
		if isAttackOrWeapon {
			state.Runechants = 0
		}

		if i < n-1 && !(m.baseGoAgain || pc.GrantedGoAgain) {
			return 0, 0, 0, false
		}
	}
	// Delayed tokens skip this turn and go straight to next turn's carryover.
	return damage, state.Runechants + state.DelayedRunechants, resources, true
}

// fillDefenseContributions writes Contribution on each Defend-role entry. The block-prevention
// share is proportional to the card's Defense() out of sumDef, capped by incomingDamage so
// over-blocking doesn't inflate attribution past what actually stopped. Defense Reactions add
// their own Play return on top, evaluated against a fresh TurnState seeded with the turn's
// pitched pool and remaining deck so card effects see the same context the solver scored them in.
func fillDefenseContributions(line []CardAssignment, pitched []card.Card, deck []card.Card, bufs *attackBufs, sumDef, incomingDamage int) {
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
}

// fillContributions populates each BestLine entry's Contribution from the winning line:
//   - Pitch:  Card.Pitch() as resource value.
//   - Defend: proportional share of Prevented plus own Play return if a Defense Reaction.
//   - Attack: per-card damage from the winning attack-chain replay.
//   - Held / Arsenal: zero (contributed nothing this turn).
//
// Called once per Best call after the partition loop picks the winner. All transient slices
// (pitched/attackers/chain/winnerOrder/perCard/used) borrow attackBufs slots so nothing
// allocates here.
func fillContributions(summary *TurnSummary, hero hero.Hero, weapons []weapon.Weapon, swungNames []string, deck []card.Card, bufs *attackBufs, incomingDamage, runechantCarryover int) {
	line := summary.BestLine

	// Reconstruct pitched, attackers, swung weapons, and resourceBudget from the winning line.
	// The arsenal-in entry (FromArsenal=true, last slot) participates in attackers / defenders
	// identically to hand entries when its role is Attack / Defend.
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

	fillDefenseContributions(line, pitched, deck, bufs, sumDef, incomingDamage)

	chain := buildAttackChain(bufs.attackerBuf[:0], attackers, weapons, swungNames)
	if len(chain) > 0 {
		ctx := &sequenceContext{
			hero:               hero,
			pitched:            pitched,
			deck:               deck,
			bufs:               bufs,
			resourceBudget:     resourceBudget,
			runechantCarryover: runechantCarryover,
			incomingDamage:     incomingDamage,
			blockTotal:         sumDef,
		}
		ctx.seedState()
		fillAttackChainContributions(summary, chain, ctx)
	}
}

// buildAttackChain appends attackers first, then the weapons named in swungNames in that order,
// so the sequence search sees the same chain composition the partition loop priced. Uses the
// passed-in slice's backing array (typically bufs.attackerBuf) to stay allocation-free.
func buildAttackChain(dst []card.Card, attackers []card.Card, weapons []weapon.Weapon, swungNames []string) []card.Card {
	dst = append(dst, attackers...)
	for _, name := range swungNames {
		for _, w := range weapons {
			if w.Name() == name {
				dst = append(dst, w)
				break
			}
		}
	}
	return dst
}

// fillAttackChainContributions re-runs the sequence search with tracking enabled to recover the
// winning permutation, snapshots it into summary.AttackChain (fresh slice — the buf-backed
// winnerOrder is reused on later calls), and maps each position's damage back to BestLine's
// Attack-role entries. Weapons don't map back since they have no BestLine entry; their damage
// is already in summary.Value. Duplicate printings played as twin attacks disambiguate by
// scan order. Contribution bundles Play return + hero-trigger so per-card stats reflect the
// card's total this-turn impact.
func fillAttackChainContributions(summary *TurnSummary, chain []card.Card, ctx *sequenceContext) {
	line := summary.BestLine
	total := len(line)
	bufs := ctx.bufs
	winnerOrder := bufs.fillContribWinnerOrder[:len(chain)]
	perCardDmg := bufs.fillContribPerCard[:len(chain)]
	if cap(bufs.fillContribTriggerDmg) < len(chain) {
		bufs.fillContribTriggerDmg = make([]float64, len(chain))
	}
	perCardTrigger := bufs.fillContribTriggerDmg[:len(chain)]
	ctx.bestSequence(chain, winnerOrder, perCardDmg, perCardTrigger)
	summary.AttackChain = make([]AttackChainEntry, len(winnerOrder))
	for i := range winnerOrder {
		summary.AttackChain[i] = AttackChainEntry{
			Card:          winnerOrder[i],
			Damage:        perCardDmg[i],
			TriggerDamage: perCardTrigger[i],
		}
	}
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
