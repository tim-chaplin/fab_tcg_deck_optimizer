// Count Your Blessings — Generic Instant. Cost 2. Printed pitch variants: Red 1, Yellow 2, Blue 3.
//
// Text: "Gain X{h}, where X is 3 plus the number of Count Your Blessings in your graveyard."

package unplayable

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

func (CountYourBlessingsRed) Play(s *sim.TurnState, l sim.Logger, self *sim.CardState) {
	l.Log(self, 0)
}

func (CountYourBlessingsYellow) Play(s *sim.TurnState, l sim.Logger, self *sim.CardState) {
	l.Log(self, 0)
}

func (CountYourBlessingsBlue) Play(s *sim.TurnState, l sim.Logger, self *sim.CardState) {
	l.Log(self, 0)
}
