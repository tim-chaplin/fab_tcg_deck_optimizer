// Nimby — Generic Action - Attack. Cost 0. Printed power: Red 3, Yellow 2, Blue 1. Printed pitch
// variants: Red 1, Yellow 2, Blue 3. Defense 2.
//
// Text: "When this attacks, you may search your deck for a Nimblism, reveal it, put it into your
// hand, then shuffle."

package unplayable

import (
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
)

func (NimbyRed) Play(g card.GameEngine, l card.Logger, self *card.CardState) {}

func (NimbyYellow) Play(g card.GameEngine, l card.Logger, self *card.CardState) {}

func (NimbyBlue) Play(g card.GameEngine, l card.Logger, self *card.CardState) {}
