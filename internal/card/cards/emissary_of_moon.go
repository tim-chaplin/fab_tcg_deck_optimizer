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
	if !ge.MoveFromHandToBottomOfDeck(self.Card.DisplayName()) {
		return
	}
	ge.DrawOne()
	l.AppendPostTriggerf(self.Card.DisplayName(), 0, "Drew a card")
}
