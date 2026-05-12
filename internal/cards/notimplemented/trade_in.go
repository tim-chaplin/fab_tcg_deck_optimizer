// Trade In — Generic Action - Attack. Cost 0. Printed power: Red 3, Yellow 2, Blue 1. Printed pitch
// variants: Red 1, Yellow 2, Blue 3. Defense 2.
//
// Text: "When this attacks, you may discard a card. If you do, draw a card. If this was played from
// arsenal, it gains **go again**."

package notimplemented

import (
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
)

// not implemented: discard-to-draw rider, arsenal-conditional go again

func (c TradeInRed) Play(g card.GameEngine, l card.Logger, self *card.CardState) {
}

// not implemented: discard-to-draw rider, arsenal-conditional go again

func (c TradeInYellow) Play(g card.GameEngine, l card.Logger, self *card.CardState) {
}

// not implemented: discard-to-draw rider, arsenal-conditional go again

func (c TradeInBlue) Play(g card.GameEngine, l card.Logger, self *card.CardState) {
}
