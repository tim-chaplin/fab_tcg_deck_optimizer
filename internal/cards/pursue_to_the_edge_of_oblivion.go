// Pursue to the Edge of Oblivion — Generic Action - Attack. Cost 0, Pitch 1, Power 4, Defense 3.
// Only printed in Red.
//
// Text: "When this hits a hero, **mark** them."

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

func (PursueToTheEdgeOfOblivionRed) Play(s sim.GameEngine, l sim.Logger, self *sim.CardState) {
	self.RegisterOnHit(markOpponentOnHit)
}
