// Package hand evaluates the value of a hand of Flesh and Blood cards
// played in isolation.
package hand

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

// Role is what a card does on a given turn cycle.
type Role uint8

const (
	Pitch Role = iota
	Attack
	Defend
)

// HandSize is the fixed number of cards in a hand.
const HandSize = 4

// Play is the chosen partition for a hand: one role per card, plus the
// resulting damage dealt and damage prevented. Roles are aligned to the
// caller's hand order.
type Play struct {
	Roles     []Role
	Dealt     int
	Prevented int
}

// Value returns the total value of the play (damage dealt + damage prevented).
func (p Play) Value() int { return p.Dealt + p.Prevented }

// Best returns the optimal Play for the given hand against an opponent
// that will attack for incomingDamage on their next turn.
//
// Cards are partitioned into three roles:
//   - Pitch: contributes its Pitch value as resources.
//   - Attack: consumes Cost resources; the attack is resolved by calling
//     Card.Play in some order the optimizer chooses. Effects on TurnState
//     carry forward to later attacks in the same sequence.
//   - Defend: contributes Defense to damage prevented (capped at
//     incomingDamage; excess block is wasted).
//
// The optimizer brute-forces all 3^N partitions, and for each legal
// partition tries every ordering of the attackers to let order-sensitive
// Play effects compound. For N=4 this is ~81 × 24 = ~2000 evaluations.
func Best(hand []card.Card, incomingDamage int) Play {
	n := len(hand)
	best := Play{Roles: make([]Role, n)}
	roles := make([]Role, n)

	var recurse func(i int)
	recurse = func(i int) {
		if i == n {
			evalPartition(hand, roles, incomingDamage, &best)
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

func evalPartition(hand []card.Card, roles []Role, incoming int, best *Play) {
	var resources, costs, defense int
	var attackers []int
	for i, c := range hand {
		switch roles[i] {
		case Pitch:
			resources += c.Pitch()
		case Attack:
			costs += c.Cost()
			attackers = append(attackers, i)
		case Defend:
			defense += c.Defense()
		}
	}
	if resources < costs {
		return
	}
	prevented := defense
	if prevented > incoming {
		prevented = incoming
	}

	dealt := bestAttackDamage(hand, attackers)

	v := dealt + prevented
	if v > best.Dealt+best.Prevented {
		best.Dealt = dealt
		best.Prevented = prevented
		copy(best.Roles, roles)
	}
}

// bestAttackDamage tries every ordering of `attackers` (indices into
// hand) and returns the max total damage after Play is called on each
// in sequence.
func bestAttackDamage(hand []card.Card, attackers []int) int {
	if len(attackers) == 0 {
		return 0
	}
	perm := append([]int(nil), attackers...)
	best := 0
	permute(perm, 0, func(order []int) {
		var state card.TurnState
		total := 0
		for _, idx := range order {
			c := hand[idx]
			total += c.Play(&state)
			state.CardsPlayed = append(state.CardsPlayed, c)
		}
		if total > best {
			best = total
		}
	})
	return best
}

func permute(a []int, k int, emit func([]int)) {
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
