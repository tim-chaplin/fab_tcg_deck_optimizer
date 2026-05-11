// Sigil of Solace — Generic Instant. Cost 0. Printed pitch variants: Red 1, Yellow 2, Blue 3.
//
// Text: "Gain 3{h}"

package unplayable

import (
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
)

func (SigilOfSolaceRed) Play(s card.GameEngine, l card.Logger, self *card.CardState) {}

func (SigilOfSolaceYellow) Play(s card.GameEngine, l card.Logger, self *card.CardState) {}

func (SigilOfSolaceBlue) Play(s card.GameEngine, l card.Logger, self *card.CardState) {}
