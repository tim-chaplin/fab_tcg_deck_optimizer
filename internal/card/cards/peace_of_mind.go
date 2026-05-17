// Peace of Mind — Generic Instant. Cost 2. Printed pitch variants: Red 1, Yellow 2, Blue 3.
//
// Text (Red): "The next time you would be dealt {p} damage, prevent 4 of that damage.
// Create a Ponder token." Yellow caps prevention at 3, Blue at 2.

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

func peaceOfMindPlay(ge card.GameEngine, l card.Logger, self *card.CardState) {
	ge.CreatePonders(1)
}

func (PeaceOfMindRed) DefensiveInstant() {}
func (PeaceOfMindRed) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	peaceOfMindPlay(ge, l, self)
}

func (PeaceOfMindYellow) DefensiveInstant() {}
func (PeaceOfMindYellow) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	peaceOfMindPlay(ge, l, self)
}

func (PeaceOfMindBlue) DefensiveInstant() {}
func (PeaceOfMindBlue) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	peaceOfMindPlay(ge, l, self)
}
