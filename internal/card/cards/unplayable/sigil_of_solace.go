// Sigil of Solace — Generic Instant. Cost 0. Printed pitch variants: Red 1, Yellow 2, Blue 3.
//
// Text: "Gain 3{h}"

package unplayable

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

func (SigilOfSolaceRed) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {}

func (SigilOfSolaceYellow) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {}

func (SigilOfSolaceBlue) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {}
