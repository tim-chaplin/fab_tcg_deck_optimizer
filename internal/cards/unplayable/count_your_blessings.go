// Count Your Blessings — Generic Instant. Cost 2. Printed pitch variants: Red 1, Yellow 2, Blue 3.
//
// Text: "Gain X{h}, where X is 3 plus the number of Count Your Blessings in your graveyard."

package unplayable

import (
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
)

func (CountYourBlessingsRed) Play(s card.GameEngine, l card.Logger, self *card.CardState) {
}

func (CountYourBlessingsYellow) Play(s card.GameEngine, l card.Logger, self *card.CardState) {
}

func (CountYourBlessingsBlue) Play(s card.GameEngine, l card.Logger, self *card.CardState) {
}
