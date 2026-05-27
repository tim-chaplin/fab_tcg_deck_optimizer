// Emissary of Tides — Generic Action - Attack. Cost 0, Pitch 1, Power 4, Defense 2. Only
// printed in Red.
//
// Text: "When this attacks, you may put a card from your hand on the bottom of your deck.
// If you do, this gets +2{p}."

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

func (EmissaryOfTidesRed) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	if !ge.MoveFromHandToBottomOfDeck(self.Card.DisplayName()) {
		return
	}
	self.BonusAttack += 2
	l.AppendPostTriggerf(self.Card.DisplayName(), 0, "+2{p}")
}
