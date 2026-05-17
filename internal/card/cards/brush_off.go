// Brush Off — Generic Instant. Cost 0. Printed pitch variants: Red 1, Yellow 2, Blue 3.
//
// Text (Red): "The next time you would be dealt 3 or less damage this turn, prevent it."
// Yellow caps at 2, Blue at 1. The "or less" gate is dropped — prevention always fires
// up to the cap.

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

func (BrushOffRed) DefensiveInstant()                                            {}
func (BrushOffRed) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {}

func (BrushOffYellow) DefensiveInstant()                                            {}
func (BrushOffYellow) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {}

func (BrushOffBlue) DefensiveInstant()                                            {}
func (BrushOffBlue) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {}
