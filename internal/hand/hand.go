// Package hand evaluates the value of a hand of Flesh and Blood cards played in isolation.
package hand

import (
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
)

// Play is the chosen partition for a hand: one role per card plus a total value score. Best sorts
// the caller's hand into canonical order in place, and Roles are aligned to that post-sort order.
type Play struct {
	Roles []Role
	// Weapons holds the names of equipped weapons that were swung in the optimal attack sequence,
	// in the order they appear in the input weapons slice. Empty if no weapon was swung.
	Weapons []string
	// Value is the play's total score (damage dealt + damage prevented). Equals the sum of
	// Contributions plus any weapon-swing damage (weapons aren't in the hand so they aren't in
	// Contributions).
	Value int
	// Contributions is per-hand-card credit toward Value, aligned with Roles. Attack role cards
	// carry their Play() return (including any hero OnCardPlayed triggers chained off them) at
	// the moment they resolved in the winning chain; Defend role cards carry their own Play
	// return for defense reactions plus their share of the Prevented block; Pitch role cards
	// carry their Pitch() resource value. Populated once per Best call on the winning line.
	Contributions []float64
	// LeftoverRunechants is the number of Runechant tokens in play at the end of the chosen
	// chain, which carry into the next turn's Best call. For partitions with no attacks, this
	// equals the carryover the caller passed in.
	LeftoverRunechants int
}

// String returns a human-readable role name ("PITCH", "ATTACK", "DEFEND").
func (r Role) String() string {
	switch r {
	case Pitch:
		return "PITCH"
	case Attack:
		return "ATTACK"
	case Defend:
		return "DEFEND"
	}
	return "UNKNOWN"
}

// FormatRoles pairs each card in hand with its assigned role for debug output, e.g.
// "Hocus Pocus (Blue): PITCH, Runic Reaping (Red): ATTACK".
func FormatRoles(hand []card.Card, roles []Role) string {
	parts := make([]string, len(hand))
	for i, c := range hand {
		parts[i] = c.Name() + ": " + roles[i].String()
	}
	return strings.Join(parts, ", ")
}

// Best returns the optimal Play for the given hand against an opponent that will attack for
// incomingDamage on their next turn. Any equipped weapons may also be swung for their Cost if
// resources allow.
//
// Cards are partitioned into three roles:
//   - Pitch: contributes its Pitch value as resources.
//   - Attack: consumes Cost resources; the attack is resolved by calling Card.Play in some order
//     the optimizer chooses. Effects on TurnState carry forward to later attacks in the same
//     sequence.
//   - Defend: contributes Defense to damage prevented (capped at incomingDamage; excess block is
//     wasted). Plain blocking is free; Defense-Reaction-typed cards must have their Cost paid to
//     resolve and also contribute any Play() damage.
//
// Pitch resources are a shared pool — attackers and Defense Reactions draw from the same
// pitchSum, so partitions whose combined cost exceeds pitchSum are illegal and pruned at the
// leaf.
//
// The optimizer brute-forces all 3^N partitions, then for each legal partition enumerates every
// subset of weapons to swing and every ordering of the combined attacker list. For N=4 with 0–2
// weapons that remains well under 10k evaluations.
//
// Results are memoized on (hero name, sorted weapon names, sorted card IDs, incomingDamage,
// runechantCarryover) so that repeated evaluations of the same hand across shuffles
// short-circuit. The hand is sorted in place into canonical order first — Roles in the returned
// Play align with that post-sort order. Every card in the hand must be registered in package
// cards; Best panics otherwise.
//
// runechantCarryover is the number of Runechant tokens carrying in from the previous turn. The
// returned Play.Leftover is the count remaining at end of the chosen chain, which the caller
// should feed back as the next turn's carryover.
func Best(hero hero.Hero, weapons []weapon.Weapon, hand []card.Card, incomingDamage int, deck []card.Card, runechantCarryover int) Play {
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

	sortHandByID(hand, ids[:], n)

	// Unmemoable hands skip the cache read but still write — the stale entry is harmless since
	// future unmemoable lookups will skip the read too.
	key := makeMemoKey(hero, weapons, &ids, n, incomingDamage, runechantCarryover)
	if memoable {
		if cached, hit := memo[key]; hit {
			// Returned Play aliases the cached slices — callers must not mutate Roles or Weapons.
			return cached
		}
	}
	result := bestUncached(hero, weapons, hand, incomingDamage, deck, runechantCarryover)
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
}

// memo caches canonical-order results keyed by memoKey. Not goroutine-safe — the simulator is
// single-threaded so no lock is needed.
var memo = map[memoKey]Play{}

// makeMemoKey builds a comparable memo key. The hand must already be sorted by card ID; weapon
// IDs are sorted numerically into the two fixed slots. sortedIDs is passed as a pointer to the
// caller's [8]card.ID stack array to avoid a slice-header escape.
func makeMemoKey(hero hero.Hero, weapons []weapon.Weapon, sortedIDs *[8]card.ID, n int, incoming int, runechantCarryover int) memoKey {
	k := memoKey{hero: hero.Name(), incoming: incoming, runechantCarryover: runechantCarryover, cardCount: uint8(n), cardIDs: *sortedIDs}
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

func bestUncached(hero hero.Hero, weapons []weapon.Weapon, hand []card.Card, incomingDamage int, deck []card.Card, runechantCarryover int) Play {
	n := len(hand)
	// Seed best.LeftoverRunechants with the carryover — partitions with no attacks don't reduce
	// the count, so carryover is the natural baseline to beat.
	best := Play{Roles: make([]Role, n), LeftoverRunechants: runechantCarryover}

	// bufs is the pooled scratch space for this deck evaluation (see getAttackBufs).
	bufs := getAttackBufs(n, weapons)
	roles := bufs.rolesBuf[:n]
	pitchVals := bufs.pitchVals[:n]
	costVals := bufs.costVals[:n]
	defendCostVals := bufs.defendCostVals[:n]
	defendPrintedVals := bufs.defendPrintedVals[:n]
	defenseVals := bufs.defenseVals[:n]

	// Pre-compute per-card pitch, cost, and defense values so the recursion can track sums
	// incrementally without interface dispatch at each partition leaf. defendCostVals carries
	// Cost only for Defense-Reaction-typed cards; non-reactions defend for free (plain blocking).
	// hasReactions lets the leaf skip groupByRoleInto's defenders bucket and the reaction-Play
	// dispatch when no card in the hand can fire a reaction — the common case.
	// defendPrintedVals[i] carries PrintedCost for defense reactions that implement
	// DiscountPerRunechant (e.g. Reduce to Runechant); hasDiscountReactions flags that any such
	// card is present, so the leaf re-checks affordability using leftoverRunechants. Slot is 0
	// for every other card, including non-discount reactions (their cost is already in
	// defendCostVals and doesn't change post-attack).
	hasReactions := false
	hasDiscountReactions := false
	for i, c := range hand {
		pitchVals[i] = c.Pitch()
		costVals[i] = c.Cost()
		defenseVals[i] = c.Defense()
		defendPrintedVals[i] = 0
		if c.Types().Has(card.TypeDefenseReaction) {
			defendCostVals[i] = costVals[i]
			hasReactions = true
			if d, ok := c.(card.DiscountPerRunechant); ok {
				defendPrintedVals[i] = d.PrintedCost()
				hasDiscountReactions = true
			}
		} else {
			defendCostVals[i] = 0
		}
	}
	pitched := bufs.pitchedBuf
	attackers := bufs.attackersBuf
	defenders := bufs.defendersBuf

	// recurse tracks defenderCostSum separately so the leaf can hand the attack pipeline a
	// resource budget of (pitchSum - defenderCostSum) to deduct chain-card effective costs from.
	var recurse func(i, pitchSum, costSum, defenseSum, defenderCostSum int)
	recurse = func(i, pitchSum, costSum, defenseSum, defenderCostSum int) {
		if i == n {
			if pitchSum < costSum {
				return
			}
			prevented := defenseSum
			if prevented > incomingDamage {
				prevented = incomingDamage
			}
			var p, a []card.Card
			var defenseDealt int
			if hasReactions {
				var d []card.Card
				p, a, d = groupByRoleInto(hand, roles, pitched[:0], attackers[:0], defenders[:0])
				defenseDealt = defenseReactionDamage(d, p, deck, bufs.state)
			} else {
				p, a = groupPitchAttack(hand, roles, pitched[:0], attackers[:0])
			}
			attackDealt, leftoverRunechants, residualBudget, swung := bestAttackWithWeapons(hero, weapons, a, p, deck, bufs, pitchSum, costSum, defenderCostSum, runechantCarryover)

			// DiscountPerRunechant defense reactions reserved 0 in defenderCostSum (their Cost() is
			// the fully-discounted minimum). Now that the attack chain has resolved and left
			// `leftoverRunechants` runechants available for defense, re-price each discounted
			// defender at its actual effective cost and reject this partition if the chain didn't
			// leave enough residual budget to cover the delta. Runechants aren't consumed by the
			// discount check — multiple discount defenders can all read the same leftover pool.
			if hasDiscountReactions {
				extraCost := 0
				for j := 0; j < n; j++ {
					if roles[j] != Defend || defendPrintedVals[j] == 0 {
						continue
					}
					actualCost := defendPrintedVals[j] - leftoverRunechants
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
			// Prefer higher Value; on ties prefer more leftover runechants — they're future
			// damage on the next turn, so they're strictly additional value over the carryover
			// baseline even when this turn's Value doesn't differentiate.
			if v > best.Value || (v == best.Value && leftoverRunechants > best.LeftoverRunechants) {
				best.Value = v
				copy(best.Roles, roles)
				best.Weapons = swung
				best.LeftoverRunechants = leftoverRunechants
			}
			return
		}
		// Defense Reactions have strict timing in FaB — they can only be played in response to an
		// opponent's attack, not chained during our own attack phase. Pruning the Attack role here
		// keeps the solver from picking partitions that use a reaction's Play effect for its
		// damage/runechant credit in offense.
		isDefenseReaction := hand[i].Types().Has(card.TypeDefenseReaction)
		for r := Role(0); r <= Defend; r++ {
			roles[i] = r
			switch r {
			case Pitch:
				recurse(i+1, pitchSum+pitchVals[i], costSum, defenseSum, defenderCostSum)
			case Attack:
				if isDefenseReaction {
					continue
				}
				recurse(i+1, pitchSum, costSum+costVals[i], defenseSum, defenderCostSum)
			case Defend:
				// Plain blocking (any card's Defense used to absorb damage) costs nothing.
				// Defense Reactions must have their Cost paid to resolve — defendCostVals carries
				// that cost for reaction cards and is zero otherwise.
				recurse(i+1, pitchSum, costSum+defendCostVals[i], defenseSum+defenseVals[i], defenderCostSum+defendCostVals[i])
			}
		}
	}
	recurse(0, 0, 0, 0, 0)
	// Once per Best call, on the winning line only, attribute per-card contribution. We skip
	// when nothing has been played (best.Roles empty — a pathological degenerate path Best never
	// actually returns from, but guard anyway).
	if len(best.Roles) > 0 {
		fillContributions(&best, hero, hand, weapons, deck, bufs, incomingDamage, runechantCarryover)
	}
	return best
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
	for mask := 0; mask < 1<<len(weapons); mask++ {
		// bufs.weaponCosts[mask] is the pre-summed Cost of the selected weapons — avoids an
		// interface dispatch per weapon on every mask.
		if pitchSum < costSum+bufs.weaponCosts[mask] {
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
		t := pc.Card.Types()
		if t.Has(card.TypeAttack) || t.Has(card.TypeWeapon) {
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

// fillContributions populates play.Contributions from the winning line. Pitch role cards credit
// Card.Pitch() (resource value); Defend role cards credit their proportional share of Prevented
// blocking plus their own Play return if they're a defense reaction; Attack role cards credit
// the per-card damage tracked during the winning attack-chain replay.
//
// Called once per Best call from bestUncached after the partition loop picks the winner. All
// transient slices (pitched/attackers/chain/winnerOrder/perCard/used) borrow attackBufs slots
// so only the returned Contributions slice allocates.
func fillContributions(play *Play, hero hero.Hero, hand []card.Card, weapons []weapon.Weapon, deck []card.Card, bufs *attackBufs, incomingDamage, runechantCarryover int) {
	n := len(hand)
	contribs := make([]float64, n)

	// Reconstruct pitched, attackers, swung weapons, and chainBudget from best.Roles plus the
	// original hand/weapons arguments. Reuses pitchedBuf/attackersBuf/attackerBuf from bufs —
	// the partition-loop writes to them are no longer needed.
	pitched := bufs.pitchedBuf[:0]
	attackers := bufs.attackersBuf[:0]
	var sumDef int
	pitchSum := 0
	defenderCostSum := 0
	for i, c := range hand {
		switch play.Roles[i] {
		case Pitch:
			pitched = append(pitched, c)
			pitchSum += c.Pitch()
		case Attack:
			attackers = append(attackers, c)
		case Defend:
			sumDef += c.Defense()
			if c.Types().Has(card.TypeDefenseReaction) {
				if _, disc := c.(card.DiscountPerRunechant); !disc {
					defenderCostSum += c.Cost()
				}
			}
		}
	}
	chainBudget := pitchSum - defenderCostSum

	// Pitch cards: Card.Pitch() as resource-value contribution.
	for i, c := range hand {
		if play.Roles[i] == Pitch {
			contribs[i] = float64(c.Pitch())
		}
	}

	// Defense: prevented share (proportional to each defender's Defense()) plus defense-reaction
	// Play return when applicable. defense-reaction Play is called with a freshly-minted state —
	// the same pattern defenseReactionDamage uses during normal scoring.
	prevented := sumDef
	if prevented > incomingDamage {
		prevented = incomingDamage
	}
	for i, c := range hand {
		if play.Roles[i] != Defend {
			continue
		}
		if sumDef > 0 {
			contribs[i] = float64(c.Defense()) * float64(prevented) / float64(sumDef)
		}
		if c.Types().Has(card.TypeDefenseReaction) {
			*bufs.state = card.TurnState{Pitched: pitched, Deck: deck}
			contribs[i] += float64(c.Play(bufs.state))
		}
	}

	// Attack chain: re-run bestAttackDamage with tracking turned on — it rediscovers the winning
	// permutation (same scoring as the partition loop's untracked call) and fills winnerOrder
	// and perCardDmg for the line that wins. Chain is the hand's attack-role cards followed by
	// the swung weapons, assembled into attackerBuf (a bufs slot, no allocation).
	chain := bufs.attackerBuf[:0]
	chain = append(chain, attackers...)
	for _, name := range play.Weapons {
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
		// Map chain-position damage back to hand indices. Weapons aren't in the hand so their
		// damage is dropped here (it's already in play.Value). For hand cards, find the first
		// unassigned Attack-role index with matching card.ID — duplicate printings played as
		// twin attacks are disambiguated by scan order.
		used := bufs.fillContribUsed[:n]
		for i := range used {
			used[i] = false
		}
		for k, c := range winnerOrder {
			if _, isWeapon := c.(weapon.Weapon); isWeapon {
				continue
			}
			for i, h := range hand {
				if used[i] || play.Roles[i] != Attack || h.ID() != c.ID() {
					continue
				}
				contribs[i] = perCardDmg[k]
				used[i] = true
				break
			}
		}
	}

	play.Contributions = contribs
}


