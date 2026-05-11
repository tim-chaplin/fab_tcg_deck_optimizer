// Pursue to the Pits of Despair — Generic Action - Attack. Cost 1, Pitch 1, Power 5, Defense 3.
// Only printed in Red.
//
// Text: "When this hits a hero, **mark** them."

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

func (PursueToThePitsOfDespairRed) Play(s *sim.TurnState, l sim.Logger, self *sim.CardState) {
	n := self.DealEffectiveAttack(s)
	self.Log(l, n)
	self.RegisterOnHit(markOpponentOnHit)
}
