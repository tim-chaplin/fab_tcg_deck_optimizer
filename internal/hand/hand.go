// Package hand evaluates the value of a hand of Flesh and Blood cards
// played in isolation.
package hand

import "github.com/timch/fab-deck-optimizer/internal/card"

// Role is what a card does on a given turn cycle.
type Role uint8

const (
	Pitch Role = iota
	Attack
	Defend
)

// Play is the chosen partition for a hand: one role per card, plus the
// resulting damage dealt and damage prevented.
type Play struct {
	Roles     []Role
	Dealt     int
	Prevented int
}

// Value returns the total value of the play (damage dealt + damage prevented).
func (p Play) Value() int { return p.Dealt + p.Prevented }

// Best returns the optimal Play for the given hand against an opponent that
// will attack for incomingDamage on their next turn.
//
// Cards in the hand are partitioned into three roles:
//   - Pitch: contributes its Pitch value as resources
//   - Attack: consumes Cost resources, contributes Attack to damage dealt
//   - Defend: held in hand, contributes Defend to damage prevented (capped
//     at incomingDamage; excess block is wasted)
//
// The optimizer brute-forces all 3^N partitions, which is fine for typical
// hand sizes (4 cards = 81 combos).
func Best(hand []card.Card, incomingDamage int) Play {
	n := len(hand)
	best := Play{Roles: make([]Role, n)}
	roles := make([]Role, n)

	var recurse func(i int)
	recurse = func(i int) {
		if i == n {
			score(hand, roles, incomingDamage, &best)
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

func score(hand []card.Card, roles []Role, incoming int, best *Play) {
	var resources, costs, dealt, defense int
	for i, c := range hand {
		switch roles[i] {
		case Pitch:
			resources += c.Pitch
		case Attack:
			costs += c.Cost
			dealt += c.Attack
		case Defend:
			defense += c.Defend
		}
	}
	if resources < costs {
		return
	}
	prevented := defense
	if prevented > incoming {
		prevented = incoming
	}
	v := dealt + prevented
	if v > best.Dealt+best.Prevented {
		best.Dealt = dealt
		best.Prevented = prevented
		copy(best.Roles, roles)
	}
}
