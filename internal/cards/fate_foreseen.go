// Fate Foreseen — Generic Defense Reaction. Cost 0.
// Printed pitch variants: Red 1, Yellow 2, Blue 3.
// Printed defense: Red 4, Yellow 3, Blue 2.
// Text: "Opt 1"

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
)

func fateForeseenPlay(g card.GameEngine, l card.Logger, self *card.CardState) {
	g.Opt(l, 1)
}

func (FateForeseenRed) Play(g card.GameEngine, l card.Logger, self *card.CardState) {
	fateForeseenPlay(g, l, self)
}

func (FateForeseenYellow) Play(g card.GameEngine, l card.Logger, self *card.CardState) {
	fateForeseenPlay(g, l, self)
}

func (FateForeseenBlue) Play(g card.GameEngine, l card.Logger, self *card.CardState) {
	fateForeseenPlay(g, l, self)
}
