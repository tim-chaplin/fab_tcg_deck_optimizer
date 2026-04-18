// Package hand evaluates the value of a hand of Flesh and Blood cards played in isolation.
package hand

import (
	"fmt"
	"sort"
	"strings"

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
	// with any swung weapons at the positions the solver chose to swing them. Swung weapons
	// can be recovered by filtering AttackChain for weapon.Weapon values, so no separate
	// Weapons field is needed. Empty when no attacks were played.
	AttackChain []card.Card
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
	appendCard := func(a CardAssignment, roleLabel string) {
		name := a.Card.Name()
		if a.FromArsenal {
			name += " (from arsenal)"
		}
		lines = append(lines, fmt.Sprintf("  %d. %s: %s", nextStep(), name, roleLabel))
	}

	for _, a := range defensePitches {
		appendCard(a, "PITCH (opponent's turn)")
	}
	for _, a := range plainBlocks {
		appendCard(a, "BLOCK")
	}
	for _, a := range defenseReactions {
		appendCard(a, "DEFENSE REACTION")
	}
	for _, a := range attackPitches {
		appendCard(a, "PITCH (my turn)")
	}
	// Attack chain: iterate AttackChain for real play order, cross-referencing BestLine by ID to
	// mark arsenal-played cards. Weapons have no BestLine entry, so they render as plain names.
	used := make([]bool, len(t.BestLine))
	for _, c := range t.AttackChain {
		if _, isWeapon := c.(weapon.Weapon); isWeapon {
			lines = append(lines, fmt.Sprintf("  %d. %s: WEAPON ATTACK", nextStep(), c.Name()))
			continue
		}
		// Find first unused BestLine entry with matching ID so we can detect FromArsenal.
		tag := ""
		for i := range t.BestLine {
			if used[i] || t.BestLine[i].Role != Attack || t.BestLine[i].Card.ID() != c.ID() {
				continue
			}
			if t.BestLine[i].FromArsenal {
				tag = " (from arsenal)"
			}
			used[i] = true
			break
		}
		lines = append(lines, fmt.Sprintf("  %d. %s%s: ATTACK", nextStep(), c.Name(), tag))
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
		if cached, hit := memo[key]; hit {
			return cached
		}
	}
	result := bestUncached(hero, weapons, hand, incomingDamage, deck, runechantCarryover, arsenalCardIn)
	memo[key] = result
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

// memoKey is a comparable struct used as the map key for memo. Hand size is capped at 8 cards;
// weapon count at 2. Weapons reference by card.ID (not name) so map hashing and equality only
// touch fixed-size integer fields plus one hero name — no string hashing per card.
type memoKey struct {
	hero               string
	weaponIDs          [2]card.ID
	cardIDs            [8]card.ID
	cardCount          uint8
	incoming           int
	runechantCarryover int
	// arsenalInID is card.Invalid when the slot is empty, otherwise the ID of the card in the
	// arsenal at the start of the turn — different arsenal-ins give distinct cache entries.
	arsenalInID card.ID
}

// memo caches canonical-order TurnSummary results keyed by memoKey. Not goroutine-safe — the
// simulator is single-threaded so no lock is needed.
var memo = map[memoKey]TurnSummary{}

// makeMemoKey builds a comparable memo key. The hand must already be sorted by card ID; weapon
// IDs are sorted numerically into the two fixed slots. sortedIDs is passed as a pointer to the
// caller's [8]card.ID stack array to avoid a slice-header escape.
func makeMemoKey(hero hero.Hero, weapons []weapon.Weapon, sortedIDs *[8]card.ID, n int, incoming int, runechantCarryover int, arsenalCardIn card.Card) memoKey {
	k := memoKey{hero: hero.Name(), incoming: incoming, runechantCarryover: runechantCarryover, cardCount: uint8(n), cardIDs: *sortedIDs}
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

// attackBufs holds pre-allocated buffers for the attack-evaluation pipeline (bestAttackDamage →
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
	// Partition-loop buffers. rolesBuf/pitchVals/costVals/defendCostVals/defendPrintedVals are
	// sized exactly handSize; pitchedBuf/attackersBuf/defendersBuf are sized 0 with cap handSize
	// and re-sliced empty at the start of each partition leaf. defendPrintedVals holds the
	// PrintedCost of DiscountPerRunechant defense reactions for the leaf's post-attack discount
	// check; slot is 0 for cards that aren't both a defense reaction AND DiscountPerRunechant.
	rolesBuf          []Role
	pitchVals         []int
	costVals          []int
	defendCostVals    []int
	defendPrintedVals []int
	defenseVals       []int
	pitchedBuf        []card.Card
	attackersBuf      []card.Card
	defendersBuf      []card.Card
	// perCardScratch is sized maxAttackers (handSize + weaponCount). Only written by playSequence
	// when the caller passes a non-nil perCardOut, and read by bestAttackDamage to snapshot the
	// winning permutation's per-card damage into the caller's output buffer. Untracked callers
	// (the partition-loop hot path) pass nil and never touch this slice.
	perCardScratch []float64
	// fillContribWinnerOrder / fillContribPerCard are output buffers for bestAttackDamage when
	// fillContributions runs the tracked replay after the partition loop. Kept on attackBufs so
	// each Best call reuses the same underlying slab instead of allocating a fresh pair.
	fillContribWinnerOrder []card.Card
	fillContribPerCard     []float64
	// fillContribUsed marks hand indices already assigned during chain→hand mapping. Sized
	// handSize; the caller resets it with clear before each fillContributions pass.
	fillContribUsed []bool
}

func newAttackBufs(handSize, weaponCount int, weapons []weapon.Weapon) *attackBufs {
	maxAttackers := handSize + weaponCount
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
		pcBuf:          make([]card.PlayedCard, maxAttackers),
		ptrBuf:         make([]*card.PlayedCard, maxAttackers),
		cardsPlayedBuf: make([]card.Card, 0, maxAttackers),
		state:          &card.TurnState{},
		attackerBuf:    make([]card.Card, maxAttackers),
		weaponCosts:    weaponCosts,
		weaponNames:    weaponNames,
		rolesBuf:          make([]Role, handSize),
		pitchVals:         make([]int, handSize),
		costVals:          make([]int, handSize),
		defendCostVals:    make([]int, handSize),
		defendPrintedVals: make([]int, handSize),
		defenseVals:       make([]int, handSize),
		pitchedBuf:             make([]card.Card, 0, handSize),
		attackersBuf:           make([]card.Card, 0, handSize),
		defendersBuf:           make([]card.Card, 0, handSize),
		perCardScratch:         make([]float64, maxAttackers),
		fillContribWinnerOrder: make([]card.Card, maxAttackers),
		fillContribPerCard:     make([]float64, maxAttackers),
		fillContribUsed:        make([]bool, handSize),
	}
}

// attackBufsCache is a single-slot cache: the simulator is single-threaded and calls bestUncached
// many times per deck with the same handSize / weapon set, so a global slot avoids rebuilding the
// ~7-slice attackBufs on every unique hand. Keyed by (handSize, weaponCount, weapon IDs).
var (
	cachedBufs       *attackBufs
	cachedHandSize   int
	cachedWeaponIDs  [2]card.ID
	cachedWeaponCt   int
	cachedBufsValid  bool
)

func getAttackBufs(handSize int, weapons []weapon.Weapon) *attackBufs {
	var wids [2]card.ID
	for i, w := range weapons {
		if i >= len(wids) {
			break
		}
		wids[i] = w.ID()
	}
	if cachedBufsValid &&
		cachedHandSize == handSize &&
		cachedWeaponCt == len(weapons) &&
		cachedWeaponIDs == wids {
		return cachedBufs
	}
	cachedBufs = newAttackBufs(handSize, len(weapons), weapons)
	cachedHandSize = handSize
	cachedWeaponCt = len(weapons)
	cachedWeaponIDs = wids
	cachedBufsValid = true
	return cachedBufs
}

func bestUncached(hero hero.Hero, weapons []weapon.Weapon, hand []card.Card, incomingDamage int, deck []card.Card, runechantCarryover int, arsenalCardIn card.Card) TurnSummary {
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
	// rebuild the chain it runs bestAttackDamage over. It lives outside TurnSummary because
	// weapons are recoverable from AttackChain once fillContributions finishes.
	var bestSwung []string
	for i := 0; i < n; i++ {
		best.BestLine[i] = CardAssignment{Card: hand[i], Role: Held}
	}
	if arsenalCardIn != nil {
		best.BestLine[n] = CardAssignment{Card: arsenalCardIn, Role: Arsenal, FromArsenal: true}
		best.ArsenalCard = arsenalCardIn
	}

	// bufs is the pooled scratch space for this deck evaluation (see getAttackBufs). Arrays are
	// sized handSize (n); the optional arsenal slot at index n uses local scratch below — small
	// enough that per-call allocation is cheap and keeps attackBufs untouched.
	bufs := getAttackBufs(n, weapons)
	// Use local scratch with capacity totalN for the per-card computed values. Small (≤ 9).
	rolesBuf := make([]Role, totalN)
	pvals := make([]int, totalN)
	cvals := make([]int, totalN)
	dvals := make([]int, totalN)
	dCostVals := make([]int, totalN)
	dPrintedVals := make([]int, totalN)
	isDR := make([]bool, totalN)

	// Pre-compute per-card pitch / cost / defense values so the recurse doesn't re-invoke
	// card-method interface calls on each partition leaf. defendCostVals holds Cost only for
	// Defense Reactions; non-reactions block for free. defendPrintedVals holds PrintedCost for
	// DiscountPerRunechant defenders and zero otherwise — used by the post-attack discount
	// re-pricing check at the leaf.
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
	// one is pruned at the role-selection step. pitchedVals is scratch reused across leaves.
	pitchedValsScratch := make([]int, 0, totalN)
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
			feasible := false
			for mask := 0; mask < 1<<len(weapons); mask++ {
				if canCoverPhasesAllUsed(pitchedVals, attackCardCost+bufs.weaponCosts[mask], drCost) {
					feasible = true
					break
				}
			}
			if !feasible {
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
			// Identify the arsenal-role card (if any) up front — we use its presence as a
			// tiebreaker: on equal Value and equal leftover runechants, a partition that puts
			// something in the arsenal beats one that doesn't, since the arsenal card carries
			// to next turn without eating a hand slot in the refill.
			var arsenalCard card.Card
			for j := 0; j < totalN; j++ {
				if rolesBuf[j] == Arsenal {
					if j < n {
						arsenalCard = hand[j]
					} else {
						arsenalCard = arsenalCardIn
					}
					break
				}
			}
			bestHasArsenal := best.ArsenalCard != nil
			hasArsenal := arsenalCard != nil
			// Tiebreak order: higher Value first, then more leftover runechants (they convert
			// to future arcane damage), then preferring a partition that uses the arsenal slot
			// (arsenal saves a hand slot in the next refill, smaller but real upside).
			better := v > best.Value
			if !better && v == best.Value {
				if leftoverRunechants > best.LeftoverRunechants {
					better = true
				} else if leftoverRunechants == best.LeftoverRunechants && hasArsenal && !bestHasArsenal {
					better = true
				}
			}
			if better {
				best.Value = v
				bestSwung = swung
				best.LeftoverRunechants = leftoverRunechants
				best.ArsenalCard = arsenalCard
				// Write the winning roles into BestLine. Cards and FromArsenal flags were
				// populated at construction; only Role varies per partition. Contribution is
				// cleared here and written by fillContributions below for the winning line.
				for j := 0; j < totalN; j++ {
					best.BestLine[j].Role = rolesBuf[j]
					best.BestLine[j].Contribution = 0
				}
			}
			return
		}
		isArsenalSlot := i == n && arsenalCardIn != nil
		for r := Role(0); r <= Arsenal; r++ {
			// Role restrictions: the arsenal slot may only take Arsenal (stay), Attack (any
			// non-DR card — auras and non-attack actions can also be played from arsenal on
			// your turn), or Defend (only if DR — plain-blocking from arsenal isn't allowed).
			// Hand cards can take any role, with Attack forbidden for DRs (strict FaB timing —
			// DRs only fire on the opponent's turn).
			if isArsenalSlot {
				switch r {
				case Pitch, Held:
					continue
				case Attack:
					if isDR[i] {
						continue
					}
				case Defend:
					if !isDR[i] {
						continue
					}
				}
			} else {
				if r == Attack && isDR[i] {
					continue
				}
			}
			if r == Arsenal && arsenalCount >= 1 {
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

// preventedDamage is the damage a wall of `defenders` blocks against `incoming`: the sum of
// printed Defense, capped at incoming (excess block is wasted).
func preventedDamage(defenders []card.Card, incoming int) int {
	total := 0
	for _, d := range defenders {
		total += d.Defense()
	}
	if total > incoming {
		return incoming
	}
	return total
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
	// Baseline: no attacks played — carryover runechants stay, and the whole chainBudget is
	// available for defense affordability checks.
	chainBudget := pitchSum - defenderCostSum
	best := 0
	bestLeftoverRunechants := runechantCarryover
	bestResidualBudget := chainBudget
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
		// Every chain card (attackers AND the swung weapons) deducts its own effective cost
		// from chainBudget inside playSequence, so we don't pre-deduct weapon cost here.
		dealt, leftoverRunechants, residualBudget := bestAttackDamage(hero, allAttackers, pitched, deck, bufs, chainBudget, runechantCarryover, nil, nil)
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

// bestAttackDamage tries every ordering of attackers and returns the max total damage after Play
// is called on each in sequence plus the runechant count at end of that winning permutation.
// Between each attacker's Play() and its append to CardsPlayed, the hero's OnCardPlayed hook
// fires so triggered abilities (e.g. Viserai's Runechants) contribute.
//
// Uses Heap's algorithm (iterative) for permutation generation. That saves a closure/callback
// allocation and a recursive call per permutation vs. a callback-style permuter.
//
// chainBudget is the resource pool available to cover chain-card effective costs. For orderings
// that run out of resources partway through, playSequence returns legal=false and contributes 0.
//
// When winnerOrderOut and perCardOut are non-nil (both must have len >= len(attackers)), they're
// filled with the winning permutation and its per-card damage respectively. Used once per Best
// call by fillContributions to attribute per-card damage; the hot partition-loop caller
// (bestAttackWithWeapons) passes nil for both so the permutation search stays allocation-free.
func bestAttackDamage(hero hero.Hero, attackers, pitched, deck []card.Card, bufs *attackBufs, chainBudget, runechantCarryover int, winnerOrderOut []card.Card, perCardOut []float64) (int, int, int) {
	n := len(attackers)
	if n == 0 {
		return 0, runechantCarryover, chainBudget
	}
	perm := bufs.perm[:n]
	copy(perm, attackers)

	// Scratch is the playSequence per-card output buffer, overwritten on every permutation; on a
	// new winner we copy it into the caller's perCardOut. Only used when the caller asked to track.
	var scratch []float64
	if perCardOut != nil {
		scratch = bufs.perCardScratch[:n]
	}

	best := 0
	bestLeftoverRunechants := runechantCarryover
	bestResidualBudget := chainBudget
	eval := func() {
		dmg, leftoverRunechants, residualBudget, legal := playSequence(hero, pitched, deck, perm, bufs.pcBuf, bufs.ptrBuf, bufs.cardsPlayedBuf, bufs.state, chainBudget, runechantCarryover, scratch)
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
		}
	}
	eval()
	// Heap's algorithm, non-recursive: c[] counts how many times each stack frame has iterated.
	var c [8]int
	i := 0
	for i < n {
		if c[i] < i {
			if i&1 == 0 {
				perm[0], perm[i] = perm[i], perm[0]
			} else {
				perm[c[i]], perm[i] = perm[i], perm[c[i]]
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

// playSequence plays `order` as an attack chain, reusing caller-provided buffers to avoid
// per-permutation heap allocation. pcBuf and ptrBuf must be at least len(order); cardsPlayedBuf
// is reset to length 0 each call. state is reset and reused each call. The buffers are mutated
// in place; the caller must not read them concurrently.
//
// When perCardOut is non-nil (len >= n) each entry is set to the damage attributed to the
// corresponding position in `order` — the card's Play return plus the hero OnCardPlayed trigger
// it chained. The hot partition-loop callers pass nil to skip this write; the winning-line
// replay from fillContributions passes a real slice.
//
// Runechant flow:
//   - state.Runechants starts at runechantCarryover (tokens from the previous turn).
//   - Each card's Play / hero OnCardPlayed may call CreateRunechants, incrementing the count AND
//     returning n damage — tokens are credited exactly once, at creation.
//   - After each Attack- or Weapon-typed card's Play+OnCardPlayed resolve, all current tokens
//     fire and are destroyed: state.Runechants is zeroed but damage is NOT re-added (that would
//     double-count tokens whose value was already credited on creation).
//   - At end of chain, state.Runechants is the leftover count that carries into the next turn.
//
// Resource flow:
//   - chainBudget starts the chain; each card deducts its effective cost. For cards implementing
//     DiscountPerRunechant, effective cost is max(0, PrintedCost() - state.Runechants) at the
//     moment it's played; for everyone else it's Cost(). A negative remaining budget returns
//     legal=false (the caller treats this ordering as zero damage).
func playSequence(hero hero.Hero, pitched, deck, order []card.Card, pcBuf []card.PlayedCard, ptrBuf []*card.PlayedCard, cardsPlayedBuf []card.Card, state *card.TurnState, chainBudget, runechantCarryover int, perCardOut []float64) (damage int, leftoverRunechants int, residualBudget int, legal bool) {
	n := len(order)
	for i, c := range order {
		pcBuf[i] = card.PlayedCard{Card: c}
		ptrBuf[i] = &pcBuf[i]
		if perCardOut != nil {
			perCardOut[i] = 0
		}
	}
	played := ptrBuf[:n]
	*state = card.TurnState{Pitched: pitched, Deck: deck, CardsPlayed: cardsPlayedBuf[:0], Runechants: runechantCarryover}
	resources := chainBudget
	for i, pc := range played {
		// Effective cost: discount cards drop by runechant count (floored at 0); others pay printed.
		var effCost int
		if d, ok := pc.Card.(card.DiscountPerRunechant); ok {
			effCost = d.PrintedCost() - state.Runechants
			if effCost < 0 {
				effCost = 0
			}
		} else {
			effCost = pc.Card.Cost()
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
		t := pc.Card.Types()
		isAttackOrWeapon := t.Has(card.TypeAttack) || t.Has(card.TypeWeapon)
		if isAttackOrWeapon && state.Runechants > 0 {
			state.ArcaneDamageDealt = true
		}

		playDmg := pc.Card.Play(state)
		triggerDmg := hero.OnCardPlayed(pc.Card, state)
		damage += playDmg + triggerDmg
		if perCardOut != nil {
			perCardOut[i] = float64(playDmg + triggerDmg)
		}
		state.CardsPlayed = append(state.CardsPlayed, pc.Card)

		// Attacks and weapon swings consume all runechants in play. Damage isn't re-added here:
		// each token was already credited as +1 at creation time (see CreateRunechants), so
		// consuming them is purely state cleanup.
		if isAttackOrWeapon {
			state.Runechants = 0
		}

		if i < n-1 && !pc.EffectiveGoAgain() {
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

	// Reconstruct pitched, attackers, swung weapons, and chainBudget from the winning line. The
	// arsenal-in entry (FromArsenal=true, last slot) participates in attackers / defenders when
	// its role is Attack / Defend, identically to hand entries.
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
	chainBudget := pitchSum - defenderCostSum

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

	// Attack chain: re-run bestAttackDamage with tracking turned on so we recover the winning
	// permutation and each chain position's contribution. Weapons in the chain don't map back
	// to BestLine (they aren't cards in hand or arsenal) — their damage is counted in Value but
	// not credited per-card here.
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
		bestAttackDamage(hero, chain, pitched, deck, bufs, chainBudget, runechantCarryover, winnerOrder, perCardDmg)
		// Snapshot the winning chain order into TurnSummary so display callers can show cards
		// and weapons in the actual play order (matters for Go-again / trigger chains). Fresh
		// slice — winnerOrder aliases attackBuf storage that the next partition would clobber.
		summary.AttackChain = append([]card.Card(nil), winnerOrder...)
		// Map chain-position damage back to BestLine indices by scanning for the first unused
		// Attack-role entry with a matching ID. Duplicate printings played as twin attacks
		// are disambiguated by scan order.
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
				line[i].Contribution = perCardDmg[k]
				used[i] = true
				break
			}
		}
	}
}


