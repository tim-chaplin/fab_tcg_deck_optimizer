// Pursue to the Pits of Despair — Generic Action - Attack. Cost 1, Pitch 1, Power 5, Defense 3.
// Only printed in Red.
//
// Text: "When this hits a hero, **mark** them."

package generic

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var pursueToThePitsOfDespairTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeAttack)

type PursueToThePitsOfDespairRed struct{}

func (PursueToThePitsOfDespairRed) ID() card.ID                 { return card.PursueToThePitsOfDespairRed }
func (PursueToThePitsOfDespairRed) Name() string                { return "Pursue to the Pits of Despair" }
func (PursueToThePitsOfDespairRed) Cost(*card.TurnState) int                   { return 1 }
func (PursueToThePitsOfDespairRed) Pitch() int                  { return 1 }
func (PursueToThePitsOfDespairRed) Attack() int                 { return 5 }
func (PursueToThePitsOfDespairRed) Defense() int                { return 3 }
func (PursueToThePitsOfDespairRed) Types() card.TypeSet         { return pursueToThePitsOfDespairTypes }
func (PursueToThePitsOfDespairRed) GoAgain() bool               { return false }
// not implemented: on-hit mark
func (PursueToThePitsOfDespairRed) NotImplemented()             {}
func (c PursueToThePitsOfDespairRed) Play(s *card.TurnState, self *card.CardState) int { return pursueToThePitsOfDespairDamage(c.Attack(), self) }

// pursueToThePitsOfDespairDamage is a breadcrumb for the on-hit "mark the hero" rider — marks
// aren't tracked (see TODO.md).
func pursueToThePitsOfDespairDamage(attack int, self *card.CardState) int {
	if card.LikelyToHit(self) {
		// TODO: model on-hit mark rider.
	}
	return attack
}
