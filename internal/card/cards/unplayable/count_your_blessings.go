// Count Your Blessings — Generic Instant. Cost 2. Printed pitch variants: Red 1, Yellow 2, Blue 3.
//
// Text: "Gain X{h}, where X is 3 plus the number of Count Your Blessings in your graveyard."

package unplayable

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

func (CountYourBlessingsRed) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
}

func (CountYourBlessingsYellow) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
}

func (CountYourBlessingsBlue) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
}
