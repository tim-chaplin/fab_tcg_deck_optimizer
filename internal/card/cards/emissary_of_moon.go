// Emissary of Moon — Generic Action - Attack. Cost 0, Pitch 1, Power 4, Defense 2. Only
// printed in Red.
//
// Text: "When this attacks, you may put a card from your hand on the bottom of your deck.
// If you do, draw a card."

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

func (EmissaryOfMoonRed) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	if len(ge.HeldHand()) == 0 {
		return
	}
	cycled := ge.PopHandAt(0)
	ge.AppendToDeck(cycled)
	ge.DrawOne()
	l.AppendPostTriggerf(self.Card.DisplayName(), 0, "Cycled %s to deck bottom and drew", cycled.DisplayName())
}
