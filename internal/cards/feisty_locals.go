// Feisty Locals — Generic Action - Attack. Cost 0. Printed power: Red 3, Yellow 2, Blue 1. Printed
// pitch variants: Red 1, Yellow 2, Blue 3. Defense 2.
//
// Text: "If this is defended by an action card, this gets +2{p}."
//
// Conservative model: the +2{p} bonus is dropped — assume the defender doesn't block with
// an action card.

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
)

func (FeistyLocalsRed) Play(g card.GameEngine, l card.Logger, self *card.CardState) {}

func (FeistyLocalsYellow) Play(g card.GameEngine, l card.Logger, self *card.CardState) {}

func (FeistyLocalsBlue) Play(g card.GameEngine, l card.Logger, self *card.CardState) {}
