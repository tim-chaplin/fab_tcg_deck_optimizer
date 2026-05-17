// Stony Woottonhog — Generic Action - Attack. Cost 2. Printed power: Red 6, Yellow 5, Blue 4.
// Printed pitch variants: Red 1, Yellow 2, Blue 3. Defense 2.
//
// Text: "While Stony Woottonhog is defended by less than 2 non-equipment cards, it has +1{p}."
//
// Conservative model: the +1{p} bonus is dropped — assume the defender always blocks with
// 2+ non-equipment cards.

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
)

func (StonyWoottonhogRed) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {}

func (StonyWoottonhogYellow) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {}

func (StonyWoottonhogBlue) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {}
