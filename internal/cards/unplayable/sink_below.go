// Sink Below — Generic Defense Reaction. Cost 0. Printed pitch variants: Red 1, Yellow 2,
// Blue 3. Printed defense: Red 4, Yellow 3, Blue 2.
//
// Text: "You may put a card from your hand on the bottom of your deck. If you do, draw a card."

package unplayable

import (
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
)

func (SinkBelowRed) Play(g card.GameEngine, l card.Logger, self *card.CardState) {}

func (SinkBelowYellow) Play(g card.GameEngine, l card.Logger, self *card.CardState) {}

func (SinkBelowBlue) Play(g card.GameEngine, l card.Logger, self *card.CardState) {}
