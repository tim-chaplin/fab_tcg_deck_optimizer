// Pursue to the Edge of Oblivion — Generic Action - Attack. Cost 0, Pitch 1, Power 4, Defense 3.
// Only printed in Red.
//
// Text: "When this hits a hero, **mark** them."

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
)

func (PursueToTheEdgeOfOblivionRed) Play(g card.GameEngine, l card.Logger, self *card.CardState) {
	self.RegisterOnHit(markOpponentOnHit)
}
