// Emissary of Wind — Generic Action - Attack. Cost 0, Pitch 1, Power 4, Defense 2. Only
// printed in Red.
//
// Text: "When this attacks, you may put a card from your hand on the bottom of your deck.
// If you do, this gets **go again**."

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

func (EmissaryOfWindRed) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	if !ge.DiscardToBottomOfDeck(self.Card.DisplayName()) {
		return
	}
	self.GrantedGoAgain = true
	l.AppendPostTriggerf(self.Card.DisplayName(), 0, "+go again")
}
