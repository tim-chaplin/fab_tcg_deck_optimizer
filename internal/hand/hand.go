// Package hand evaluates the value of a hand of Flesh and Blood cards played in isolation.
package hand

import (
	"sort"
	"strings"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/cards"
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
//     wasted).
//
// The optimizer brute-forces all 3^N partitions, then for each legal partition enumerates every
// subset of weapons to swing and every ordering of the combined attacker list. For N=4 with 0–2
// weapons that remains well under 10k evaluations.
//
// Results are memoized on (hero name, sorted weapon names, sorted card IDs, incomingDamage) so
// that repeated evaluations of the same hand across shuffles short-circuit. The hand is sorted
// in place into canonical order first — Roles in the returned Play align with that post-sort
// order. Every card in the hand must be registered in package cards; Best panics otherwise.
func Best(hero hero.Hero, weapons []weapon.Weapon, hand []card.Card, incomingDamage int, deck []card.Card) Play {
	ids := handIDs(hand)
	sort.Sort(&handByID{hand: hand, ids: ids})

	// Unmemoable hands skip the cache read but still write — the stale entry is harmless since
	// future unmemoable lookups will skip the read too.
	key := makeMemoKey(hero, weapons, ids, incomingDamage)
	if isMemoable(hand) {
		if cached, hit := memo[key]; hit {
			// Returned Play aliases the cached slices — callers must not mutate Roles or Weapons.
			return cached
		}
	}
	result := bestUncached(hero, weapons, hand, incomingDamage, deck)
	memo[key] = result
	return result
}

// isMemoable reports whether a hand's Best result can be cached. Any card implementing
// card.NoMemo (e.g. one whose Play depends on deck composition not captured by the memo key)
// disqualifies the whole hand.
func isMemoable(hand []card.Card) bool {
	for _, c := range hand {
		if _, ok := c.(card.NoMemo); ok {
			return false
		}
	}
	return true
}

// handIDs returns the registry IDs for each card in hand, preserving order. Panics if any card
// isn't registered in package cards.
func handIDs(hand []card.Card) []cards.ID {
	ids := make([]cards.ID, len(hand))
	for i, c := range hand {
		id, found := cards.ByName(c.Name())
		if !found {
			panic("hand: card not in index: " + c.Name())
		}
		ids[i] = id
	}
	return ids
}

// handByID sorts a parallel (hand, ids) pair by ascending ID.
type handByID struct {
	hand []card.Card
	ids  []cards.ID
}

func (h *handByID) Len() int           { return len(h.ids) }
func (h *handByID) Less(i, j int) bool { return h.ids[i] < h.ids[j] }
func (h *handByID) Swap(i, j int) {
	h.ids[i], h.ids[j] = h.ids[j], h.ids[i]
	h.hand[i], h.hand[j] = h.hand[j], h.hand[i]
}

// memoKey is a comparable struct used as the map key for memo — avoids the string allocations
// that a formatted key would require on every hand evaluation. Hand size is capped at 8 cards;
// weapon count at 2.
type memoKey struct {
	hero      string
	weapon0   string
	weapon1   string
	cardIDs   [8]cards.ID
	cardCount uint8
	incoming  int
}

// memo caches canonical-order results keyed by memoKey. Not goroutine-safe — the simulator is
// single-threaded so no lock is needed.
var memo = map[memoKey]Play{}

// makeMemoKey builds a comparable memo key. The hand must already be sorted by card ID; weapon
// names are sorted lexicographically into the two fixed slots.
func makeMemoKey(hero hero.Hero, weapons []weapon.Weapon, sortedIDs []cards.ID, incoming int) memoKey {
	k := memoKey{hero: hero.Name(), incoming: incoming, cardCount: uint8(len(sortedIDs))}
	for i, id := range sortedIDs {
		k.cardIDs[i] = id
	}
	switch len(weapons) {
	case 1:
		k.weapon0 = weapons[0].Name()
	case 2:
		a, b := weapons[0].Name(), weapons[1].Name()
		if a > b {
			a, b = b, a
		}
		k.weapon0, k.weapon1 = a, b
	}
	return k
}

// attackBufs holds pre-allocated buffers for the attack-evaluation pipeline (bestAttackDamage →
// playSequence). Allocated once per bestUncached call and reused across all partitions and weapon
// masks to avoid repeated heap allocation.
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
	}
}

func bestUncached(hero hero.Hero, weapons []weapon.Weapon, hand []card.Card, incomingDamage int, deck []card.Card) Play {
	n := len(hand)
	best := Play{Roles: make([]Role, n)}
	roles := make([]Role, n)

	// Pre-compute per-card pitch and cost values so the recursion can track sums incrementally
	// without interface dispatch at each partition leaf.
	pitchVals := make([]int, n)
	costVals := make([]int, n)
	for i, c := range hand {
		pitchVals[i] = c.Pitch()
		costVals[i] = c.Cost()
	}

	// Pre-allocate all buffers once for the entire partition search.
	bufs := newAttackBufs(n, len(weapons), weapons)
	pitched := make([]card.Card, 0, n)
	attackers := make([]card.Card, 0, n)
	defenders := make([]card.Card, 0, n)

	var recurse func(i, pitchSum, costSum int)
	recurse = func(i, pitchSum, costSum int) {
		if i == n {
			if pitchSum < costSum {
				return
			}
			p, a, d := groupByRoleInto(hand, roles, pitched[:0], attackers[:0], defenders[:0])
			prevented := preventedDamage(d, incomingDamage)
			defenseDealt := defenseReactionDamage(d, p, deck)
			attackDealt, swung := bestAttackWithWeapons(hero, weapons, a, p, deck, bufs, pitchSum, costSum)

			v := attackDealt + defenseDealt + prevented
			if v > best.Value {
				best.Value = v
				copy(best.Roles, roles)
				best.Weapons = swung
			}
			return
		}
		for r := Role(0); r <= Defend; r++ {
			roles[i] = r
			switch r {
			case Pitch:
				recurse(i+1, pitchSum+pitchVals[i], costSum)
			case Attack:
				recurse(i+1, pitchSum, costSum+costVals[i])
			default:
				recurse(i+1, pitchSum, costSum)
			}
		}
	}
	recurse(0, 0, 0)
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

// canAfford reports whether the pitch resources produced by pitched cover the combined Cost of
// every card in attackers. A partition is legal only if its attackers (hand cards plus any weapon
// swings joined into the same list) can all be paid for.
func canAfford(pitched, attackers []card.Card) bool {
	resources := 0
	for _, c := range pitched {
		resources += c.Pitch()
	}
	cost := 0
	for _, c := range attackers {
		cost += c.Cost()
	}
	return resources >= cost
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
func defenseReactionDamage(defenders, pitched, deck []card.Card) int {
	total := 0
	for _, d := range defenders {
		if !d.Types().Has(card.TypeDefenseReaction) {
			continue
		}
		state := card.TurnState{Pitched: pitched, Deck: deck}
		total += d.Play(&state)
	}
	return total
}

// bestAttackWithWeapons enumerates every subset of `weapons` to swing alongside `attackers` and
// returns the max damage over all (affordable) weapon masks, plus the swung weapons from the
// winning mask (in input order). Each selected weapon adds its Cost and joins the attacker
// permutation inside bestAttackDamage.
func bestAttackWithWeapons(hero hero.Hero, weapons []weapon.Weapon, attackers, pitched, deck []card.Card, bufs *attackBufs, pitchSum, costSum int) (int, []string) {
	best := 0
	var bestSwung []string
	// Reuse the shared attacker buffer across mask iterations.
	copy(bufs.attackerBuf, attackers)
	for mask := 0; mask < 1<<len(weapons); mask++ {
		// Use pre-computed weapon costs instead of iterating through interface calls.
		if pitchSum < costSum+bufs.weaponCosts[mask] {
			continue
		}
		allAttackers := bufs.attackerBuf[:len(attackers)]
		for i, w := range weapons {
			if mask&(1<<i) != 0 {
				allAttackers = append(allAttackers, w)
			}
		}
		if dealt := bestAttackDamage(hero, allAttackers, pitched, deck, bufs); dealt > best {
			best = dealt
			bestSwung = bufs.weaponNames[mask]
		}
	}
	return best, bestSwung
}

// bestAttackDamage tries every ordering of attackers and returns the max total damage after Play
// is called on each in sequence. Between each attacker's Play() and its append to CardsPlayed,
// the hero's OnCardPlayed hook fires so triggered abilities (e.g. Viserai's Runechants)
// contribute.
//
// All buffers come from the shared attackBufs allocated once per bestUncached call; only the
// permutation array is reset per call.
func bestAttackDamage(hero hero.Hero, attackers, pitched, deck []card.Card, bufs *attackBufs) int {
	n := len(attackers)
	if n == 0 {
		return 0
	}
	perm := bufs.perm[:n]
	copy(perm, attackers)

	best := 0
	permute(perm, 0, func(order []card.Card) {
		if dmg, legal := playSequence(hero, pitched, deck, order, bufs.pcBuf, bufs.ptrBuf, bufs.cardsPlayedBuf, bufs.state); legal && dmg > best {
			best = dmg
		}
	})
	return best
}

// playSequence plays `order` as an attack chain, reusing caller-provided buffers to avoid
// per-permutation heap allocation. pcBuf and ptrBuf must be at least len(order); cardsPlayedBuf
// is reset to length 0 each call. state is reset and reused each call. The buffers are mutated
// in place; the caller must not read them concurrently.
func playSequence(hero hero.Hero, pitched, deck, order []card.Card, pcBuf []card.PlayedCard, ptrBuf []*card.PlayedCard, cardsPlayedBuf []card.Card, state *card.TurnState) (damage int, legal bool) {
	n := len(order)
	for i, c := range order {
		pcBuf[i] = card.PlayedCard{Card: c}
		ptrBuf[i] = &pcBuf[i]
	}
	played := ptrBuf[:n]
	*state = card.TurnState{Pitched: pitched, Deck: deck, CardsPlayed: cardsPlayedBuf[:0]}
	for i, pc := range played {
		state.CardsRemaining = played[i+1:]
		state.Self = pc
		damage += pc.Card.Play(state)
		damage += hero.OnCardPlayed(pc.Card, state)
		state.CardsPlayed = append(state.CardsPlayed, pc.Card)
		if i < n-1 && !pc.EffectiveGoAgain() {
			return 0, false
		}
	}
	return damage, true
}


func permute(a []card.Card, k int, emit func([]card.Card)) {
	if k == len(a)-1 {
		emit(a)
		return
	}
	for i := k; i < len(a); i++ {
		a[k], a[i] = a[i], a[k]
		permute(a, k+1, emit)
		a[k], a[i] = a[i], a[k]
	}
}
