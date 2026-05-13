// Peace of Mind — Generic Instant. Cost 2. Printed pitch variants: Red 1, Yellow 2, Blue 3.
//
// Text (Red): "The next time you would be dealt {p} damage, prevent 4 of that damage.
// Create a Ponder token." Yellow caps prevention at 3, Blue at 2.

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
)

func peaceOfMindPlay(g card.GameEngine, l card.Logger, self *card.CardState) {
	g.CreatePonders(1)
}

func (PeaceOfMindRed) DefensiveInstant() {}
func (PeaceOfMindRed) Play(g card.GameEngine, l card.Logger, self *card.CardState) {
	peaceOfMindPlay(g, l, self)
}

func (PeaceOfMindYellow) DefensiveInstant() {}
func (PeaceOfMindYellow) Play(g card.GameEngine, l card.Logger, self *card.CardState) {
	peaceOfMindPlay(g, l, self)
}

func (PeaceOfMindBlue) DefensiveInstant() {}
func (PeaceOfMindBlue) Play(g card.GameEngine, l card.Logger, self *card.CardState) {
	peaceOfMindPlay(g, l, self)
}
