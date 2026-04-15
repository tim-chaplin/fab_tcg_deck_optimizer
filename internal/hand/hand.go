// Package hand evaluates the value of a hand of Flesh and Blood cards played in isolation.
package hand

import (
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

// Play is the chosen partition for a hand: one role per card, plus the resulting damage dealt and
// damage prevented. Roles are aligned to the caller's hand order. (Weapon swing decisions are not
// reported in Roles — they're consumed only for their damage contribution.)
type Play struct {
	Roles     []Role
	Dealt     int
	Prevented int
}

// Value returns the total value of the play (damage dealt + damage prevented).
func (p Play) Value() int { return p.Dealt + p.Prevented }

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
func Best(h hero.Hero, weapons []weapon.Weapon, hand []card.Card, incomingDamage int) Play {
	n := len(hand)
	best := Play{Roles: make([]Role, n)}
	roles := make([]Role, n)

	var recurse func(i int)
	recurse = func(i int) {
		if i == n {
			evalPartition(h, weapons, hand, roles, incomingDamage, &best)
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

func evalPartition(h hero.Hero, weapons []weapon.Weapon, hand []card.Card, roles []Role, incoming int, best *Play) {
	var resources, cardCosts, defense int
	var cardAttackers []card.Card
	for i, c := range hand {
		switch roles[i] {
		case Pitch:
			resources += c.Pitch()
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
		if dealt := bestAttackDamage(h, attackers); dealt > bestDealt {
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

// bestAttackDamage tries every ordering of attackers and returns the max total damage after Play is
// called on each in sequence. Between each attacker's Play() and its append to CardsPlayed, the
// hero's OnCardPlayed hook fires so triggered abilities (e.g. Viserai's Runechants) contribute.
func bestAttackDamage(h hero.Hero, attackers []card.Card) int {
	if len(attackers) == 0 {
		return 0
	}
	perm := make([]card.Card, len(attackers))
	copy(perm, attackers)
	best := 0
	permute(perm, 0, func(order []card.Card) {
		// Illegal: any card without Go again ends the chain, so all but the last must have GoAgain.
		for i := 0; i < len(order)-1; i++ {
			if !order[i].GoAgain() {
				return
			}
		}
		var state card.TurnState
		total := 0
		for i, c := range order {
			state.CardsRemaining = order[i+1:]
			total += c.Play(&state)
			total += h.OnCardPlayed(c, &state)
			state.CardsPlayed = append(state.CardsPlayed, c)
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
