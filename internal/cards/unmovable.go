// Unmovable — Generic Defense Reaction. Cost 3.
// Printed pitch variants: Red 1, Yellow 2, Blue 3.
// Printed defense: Red 7, Yellow 6, Blue 5.
// Text: "If Unmovable is played from arsenal, it gains +1{d}."

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
)

func (UnmovableRed) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
}
func (UnmovableRed) ArsenalDefenseBonus() int { return 1 }

func (UnmovableYellow) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
}
func (UnmovableYellow) ArsenalDefenseBonus() int { return 1 }

func (UnmovableBlue) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
}
func (UnmovableBlue) ArsenalDefenseBonus() int { return 1 }
