// Peace of Mind — Generic Instant. Cost 2. Printed pitch variants: Red 1, Yellow 2, Blue 3.
//
// Text (Red): "The next time you would be dealt {p} damage, prevent 4 of that damage.
// Create a Ponder token." Yellow caps prevention at 3, Blue at 2.

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

func peaceOfMindPlay(s *sim.TurnState, l sim.Logger, self *sim.CardState) {
	s.CreatePonder(1)
}

func (PeaceOfMindRed) DefensiveInstant() {}
func (PeaceOfMindRed) Play(s *sim.TurnState, l sim.Logger, self *sim.CardState) {
	peaceOfMindPlay(s, l, self)
}

func (PeaceOfMindYellow) DefensiveInstant() {}
func (PeaceOfMindYellow) Play(s *sim.TurnState, l sim.Logger, self *sim.CardState) {
	peaceOfMindPlay(s, l, self)
}

func (PeaceOfMindBlue) DefensiveInstant() {}
func (PeaceOfMindBlue) Play(s *sim.TurnState, l sim.Logger, self *sim.CardState) {
	peaceOfMindPlay(s, l, self)
}
