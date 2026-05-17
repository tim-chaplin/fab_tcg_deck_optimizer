// Barraging Brawnhide — Generic Action - Attack. Cost 3. Printed power: Red 7, Yellow 6, Blue 5.
// Printed pitch variants: Red 1, Yellow 2, Blue 3. Defense 2.
//
// Text: "While Barraging Brawnhide is defended by less than 2 non-equipment cards, it has +1{p}."
//
// Conservative model: the +1{p} bonus is dropped — assume the defender always blocks with
// 2+ non-equipment cards.

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

func (BarragingBrawnhideRed) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {}

func (BarragingBrawnhideYellow) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {}

func (BarragingBrawnhideBlue) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {}
