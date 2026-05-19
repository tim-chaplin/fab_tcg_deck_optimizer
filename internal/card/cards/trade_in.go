// Trade In — Generic Action - Attack. Cost 0. Printed power: Red 3, Yellow 2, Blue 1.
// Printed pitch variants: Red 1, Yellow 2, Blue 3. Defense 2.
//
// Text: "When this attacks, you may discard a card. If you do, draw a card. If this was
// played from arsenal, it gains **go again**."

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

func tradeInPlay(ge card.GameEngine, l card.Logger, self *card.CardState) {
	self.GrantGoAgainIfFromArsenal()
	discarded, ok := ge.Discard()
	if !ok {
		return
	}
	ge.DrawOne()
	l.AppendPostTriggerf(self.Card.DisplayName(), 0, "Discarded %s and drew", discarded.DisplayName())
}

func (TradeInRed) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	tradeInPlay(ge, l, self)
}

func (TradeInYellow) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	tradeInPlay(ge, l, self)
}

func (TradeInBlue) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	tradeInPlay(ge, l, self)
}
