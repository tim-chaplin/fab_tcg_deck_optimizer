// Package hand evaluates the value of a hand of Flesh and Blood cards played in isolation.
package hand

import (
	"sort"
	"strconv"
	"strings"
	"sync"

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

// Play is the chosen partition for a hand: one role per card, plus the resulting damage dealt and
// damage prevented. Best sorts the caller's hand into canonical order in place, and Roles are
// aligned to that post-sort order. (Weapon swing decisions are not reported in Roles — they're
// consumed only for their damage contribution.)
type Play struct {
	Roles     []Role
	Dealt     int
	Prevented int
}

// Value returns the total value of the play (damage dealt + damage prevented).
func (p Play) Value() int { return p.Dealt + p.Prevented }

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
func Best(hero hero.Hero, weapons []weapon.Weapon, hand []card.Card, incomingDamage int) Play {
	ids := handIDs(hand)
	sort.Sort(&handByID{hand: hand, ids: ids})
	key := formatMemoKey(hero, weapons, ids, incomingDamage)

	memoMu.RLock()
	cached, hit := memo[key]
	memoMu.RUnlock()
	if hit {
		// Clone Roles so callers can't mutate the cached slice.
		roles := make([]Role, len(cached.Roles))
		copy(roles, cached.Roles)
		return Play{Roles: roles, Dealt: cached.Dealt, Prevented: cached.Prevented}
	}

	result := bestUncached(hero, weapons, hand, incomingDamage)

	memoMu.Lock()
	memo[key] = result
	memoMu.Unlock()

	return result
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

// memo caches canonical-order results keyed by formatMemoKey.
var (
	memoMu sync.RWMutex
	memo   = map[string]Play{}
)

// formatMemoKey serializes the canonical key fields into a string. The hand must already be
// sorted by card ID; weapon names are sorted here.
func formatMemoKey(hero hero.Hero, weapons []weapon.Weapon, sortedIDs []cards.ID, incoming int) string {
	wnames := make([]string, len(weapons))
	for i, w := range weapons {
		wnames[i] = w.Name()
	}
	sort.Strings(wnames)

	var b strings.Builder
	b.WriteString(hero.Name())
	b.WriteByte('|')
	b.WriteString(strings.Join(wnames, ","))
	b.WriteByte('|')
	for i, id := range sortedIDs {
		if i > 0 {
			b.WriteByte(',')
		}
		b.WriteString(strconv.FormatUint(uint64(id), 10))
	}
	b.WriteByte('|')
	b.WriteString(strconv.Itoa(incoming))
	return b.String()
}

func bestUncached(hero hero.Hero, weapons []weapon.Weapon, hand []card.Card, incomingDamage int) Play {
	n := len(hand)
	best := Play{Roles: make([]Role, n)}
	roles := make([]Role, n)

	var recurse func(i int)
	recurse = func(i int) {
		if i == n {
			evalPartition(hero, weapons, hand, roles, incomingDamage, &best)
			return
		}
		for r := Role(0); r <= Defend; r++ {
			roles[i] = r
			recurse(i + 1)
		}
	}
	recurse(0)
	return best
}

func evalPartition(hero hero.Hero, weapons []weapon.Weapon, hand []card.Card, roles []Role, incoming int, best *Play) {
	var resources, cardCosts, defense int
	var cardAttackers, pitched []card.Card
	for i, c := range hand {
		switch roles[i] {
		case Pitch:
			resources += c.Pitch()
			pitched = append(pitched, c)
		case Attack:
			cardCosts += c.Cost()
			cardAttackers = append(cardAttackers, c)
		case Defend:
			defense += c.Defense()
		}
	}
	if resources < cardCosts {
		return
	}
	prevented := defense
	if prevented > incoming {
		prevented = incoming
	}

	// Enumerate every subset of weapons to swing. Each selected weapon adds its Cost and joins the
	// attacker permutation.
	bestDealt := 0
	for mask := 0; mask < 1<<len(weapons); mask++ {
		totalCost := cardCosts
		attackers := append([]card.Card(nil), cardAttackers...)
		for i, w := range weapons {
			if mask&(1<<i) != 0 {
				totalCost += w.Cost()
				attackers = append(attackers, w)
			}
		}
		if resources < totalCost {
			continue
		}
		if dealt := bestAttackDamage(hero, attackers, pitched); dealt > bestDealt {
			bestDealt = dealt
		}
	}

	v := bestDealt + prevented
	if v > best.Dealt+best.Prevented {
		best.Dealt = bestDealt
		best.Prevented = prevented
		copy(best.Roles, roles)
	}
}

// bestAttackDamage tries every ordering of attackers and returns the max total damage after Play
// is called on each in sequence. Between each attacker's Play() and its append to CardsPlayed,
// the hero's OnCardPlayed hook fires so triggered abilities (e.g. Viserai's Runechants)
// contribute.
//
// Each permutation gets its own freshly-allocated []*PlayedCard so that any grants applied via
// CardsRemaining mutation don't leak across permutations.
func bestAttackDamage(hero hero.Hero, attackers, pitched []card.Card) int {
	if len(attackers) == 0 {
		return 0
	}
	perm := make([]card.Card, len(attackers))
	copy(perm, attackers)
	best := 0
	permute(perm, 0, func(order []card.Card) {
		// Fresh PlayedCard wrappers per permutation so prior permutations' grants don't carry over.
		played := make([]*card.PlayedCard, len(order))
		for i, c := range order {
			played[i] = &card.PlayedCard{Card: c}
		}

		state := card.TurnState{Pitched: pitched}
		total := 0
		for i, pc := range played {
			// Chain legality: non-final cards need an action point, which comes from printed Go
			// again OR from a grant applied to this PlayedCard by a prior card's Play (e.g.
			// Mauvrion Skies flipping pc.GrantedGoAgain on a matching later entry).
			state.CardsRemaining = played[i+1:]
			total += pc.Card.Play(&state)
			total += hero.OnCardPlayed(pc.Card, &state)
			state.CardsPlayed = append(state.CardsPlayed, pc.Card)

			if i < len(played)-1 && !pc.EffectiveGoAgain() {
				return
			}
		}
		if total > best {
			best = total
		}
	})
	return best
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
