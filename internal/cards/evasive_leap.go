// Evasive Leap — Generic Defense Reaction. Cost 0.
// Printed pitch variants: Red 1, Yellow 2, Blue 3.
// Printed defense: Red 3, Yellow 2, Blue 1.

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

func (EvasiveLeapRed) Play(s *sim.TurnState, l sim.Logger, self *sim.CardState) {
}

func (EvasiveLeapYellow) Play(s *sim.TurnState, l sim.Logger, self *sim.CardState) {
}

func (EvasiveLeapBlue) Play(s *sim.TurnState, l sim.Logger, self *sim.CardState) {
}
