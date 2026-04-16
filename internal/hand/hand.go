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
	// Value is the play's total score (damage dealt + damage prevented). The breakdown is not
	// tracked on Play directly — a future stats object may reintroduce it.
	Value int
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
	// Partition-loop buffers. rolesBuf/pitchVals/costVals/defendCostVals are sized exactly
	// handSize; pitchedBuf/attackersBuf/defendersBuf are sized 0 with cap handSize and re-sliced
	// empty at the start of each partition leaf.
	rolesBuf       []Role
	pitchVals      []int
	costVals       []int
	defendCostVals []int
	defenseVals    []int
	pitchedBuf     []card.Card
	attackersBuf   []card.Card
	defendersBuf   []card.Card
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
		rolesBuf:       make([]Role, handSize),
		pitchVals:      make([]int, handSize),
		costVals:       make([]int, handSize),
		defendCostVals: make([]int, handSize),
		defenseVals:    make([]int, handSize),
		pitchedBuf:     make([]card.Card, 0, handSize),
		attackersBuf:   make([]card.Card, 0, handSize),
		defendersBuf:   make([]card.Card, 0, handSize),
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
	defenseVals := bufs.defenseVals[:n]

	// Pre-compute per-card pitch, cost, and defense values so the recursion can track sums
	// incrementally without interface dispatch at each partition leaf. defendCostVals carries
	// Cost only for Defense-Reaction-typed cards; non-reactions defend for free (plain blocking).
	// hasReactions lets the leaf skip groupByRoleInto's defenders bucket and the reaction-Play
	// dispatch when no card in the hand can fire a reaction — the common case.
	hasReactions := false
	for i, c := range hand {
		pitchVals[i] = c.Pitch()
		costVals[i] = c.Cost()
		defenseVals[i] = c.Defense()
		if c.Types().Has(card.TypeDefenseReaction) {
			defendCostVals[i] = costVals[i]
			hasReactions = true
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
			attackDealt, leftoverRunechants, swung := bestAttackWithWeapons(hero, weapons, a, p, deck, bufs, pitchSum, costSum, defenderCostSum, runechantCarryover)

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
		for r := Role(0); r <= Defend; r++ {
			roles[i] = r
			switch r {
			case Pitch:
				recurse(i+1, pitchSum+pitchVals[i], costSum, defenseSum, defenderCostSum)
			case Attack:
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
// returns the max damage over all (affordable) weapon masks, plus the runechant leftover from
// that winning mask and the swung weapons (in input order). Each selected weapon adds its Cost
// and joins the attacker permutation inside bestAttackDamage.
//
// resourceBudget passed down to the chain pipeline is pitchSum - defenderCostSum - weaponCosts
// for this mask; the chain then deducts each attacker's effective cost (which for
// DiscountPerRunechant cards is max(0, PrintedCost() - runechantCount at play-time)) and rejects
// orderings that run negative.
func bestAttackWithWeapons(hero hero.Hero, weapons []weapon.Weapon, attackers, pitched, deck []card.Card, bufs *attackBufs, pitchSum, costSum, defenderCostSum, runechantCarryover int) (int, int, []string) {
	// Baseline: no attacks played. Leftover runechants stay at the caller's carryover.
	best := 0
	bestLeftoverRunechants := runechantCarryover
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
		// Chain budget = pitchSum minus defender costs. Every chain card (attackers AND the
		// swung weapons) deducts its own effective cost from this budget inside playSequence,
		// so we don't pre-deduct weapon cost here.
		chainBudget := pitchSum - defenderCostSum
		dealt, leftoverRunechants := bestAttackDamage(hero, allAttackers, pitched, deck, bufs, chainBudget, runechantCarryover)
		// Prefer higher damage; on ties prefer more leftover runechants for future turns.
		if dealt > best || (dealt == best && leftoverRunechants > bestLeftoverRunechants) {
			best = dealt
			bestLeftoverRunechants = leftoverRunechants
			bestSwung = bufs.weaponNames[mask]
		}
	}
	return best, bestLeftoverRunechants, bestSwung
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
func bestAttackDamage(hero hero.Hero, attackers, pitched, deck []card.Card, bufs *attackBufs, chainBudget, runechantCarryover int) (int, int) {
	n := len(attackers)
	if n == 0 {
		return 0, runechantCarryover
	}
	perm := bufs.perm[:n]
	copy(perm, attackers)

	best := 0
	bestLeftoverRunechants := runechantCarryover
	eval := func() {
		dmg, leftoverRunechants, legal := playSequence(hero, pitched, deck, perm, bufs.pcBuf, bufs.ptrBuf, bufs.cardsPlayedBuf, bufs.state, chainBudget, runechantCarryover)
		if !legal {
			return
		}
		// Prefer higher damage; on ties prefer more leftover runechants for future turns.
		if dmg > best || (dmg == best && leftoverRunechants > bestLeftoverRunechants) {
			best = dmg
			bestLeftoverRunechants = leftoverRunechants
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
	return best, bestLeftoverRunechants
}

// playSequence plays `order` as an attack chain, reusing caller-provided buffers to avoid
// per-permutation heap allocation. pcBuf and ptrBuf must be at least len(order); cardsPlayedBuf
// is reset to length 0 each call. state is reset and reused each call. The buffers are mutated
// in place; the caller must not read them concurrently.
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
func playSequence(hero hero.Hero, pitched, deck, order []card.Card, pcBuf []card.PlayedCard, ptrBuf []*card.PlayedCard, cardsPlayedBuf []card.Card, state *card.TurnState, chainBudget, runechantCarryover int) (damage int, leftoverRunechants int, legal bool) {
	n := len(order)
	for i, c := range order {
		pcBuf[i] = card.PlayedCard{Card: c}
		ptrBuf[i] = &pcBuf[i]
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
			return 0, 0, false
		}

		state.CardsRemaining = played[i+1:]
		state.Self = pc
		damage += pc.Card.Play(state)
		damage += hero.OnCardPlayed(pc.Card, state)
		state.CardsPlayed = append(state.CardsPlayed, pc.Card)

		// Attacks and weapon swings consume all runechants in play. Damage isn't re-added here:
		// each token was already credited as +1 at creation time (see CreateRunechants), so
		// consuming them is purely state cleanup.
		t := pc.Card.Types()
		if t.Has(card.TypeAttack) || t.Has(card.TypeWeapon) {
			state.Runechants = 0
		}

		if i < n-1 && !pc.EffectiveGoAgain() {
			return 0, 0, false
		}
	}
	// Delayed tokens (e.g. from Blessing of Occult) skip this turn and go straight to next
	// turn's carryover.
	return damage, state.Runechants + state.DelayedRunechants, true
}


