// Pursue to the Pits of Despair — Generic Action - Attack. Cost 1, Pitch 1, Power 5, Defense 3.
// Only printed in Red.
//
// Text: "When this hits a hero, **mark** them."

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
)

func (PursueToThePitsOfDespairRed) Play(g card.GameEngine, l card.Logger, self *card.CardState) {
	self.RegisterOnHit(markOpponentOnHit)
}
