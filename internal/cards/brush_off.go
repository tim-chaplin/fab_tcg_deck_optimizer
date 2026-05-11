// Brush Off — Generic Instant. Cost 0. Printed pitch variants: Red 1, Yellow 2, Blue 3.
//
// Text (Red): "The next time you would be dealt 3 or less damage this turn, prevent it."
// Yellow caps at 2, Blue at 1. The "or less" gate is dropped — prevention always fires
// up to the cap.

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

func brushOffPlay(s *sim.TurnState, l sim.Logger, self *sim.CardState) {
	n := self.DealEffectiveDefense(s)
	self.Log(l, n)
}

func (BrushOffRed) DefensiveInstant() {}
func (BrushOffRed) Play(s *sim.TurnState, l sim.Logger, self *sim.CardState) {
	brushOffPlay(s, l, self)
}

func (BrushOffYellow) DefensiveInstant() {}
func (BrushOffYellow) Play(s *sim.TurnState, l sim.Logger, self *sim.CardState) {
	brushOffPlay(s, l, self)
}

func (BrushOffBlue) DefensiveInstant() {}
func (BrushOffBlue) Play(s *sim.TurnState, l sim.Logger, self *sim.CardState) {
	brushOffPlay(s, l, self)
}
