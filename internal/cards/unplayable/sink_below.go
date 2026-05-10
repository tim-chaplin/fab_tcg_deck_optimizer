// Sink Below — Generic Defense Reaction. Cost 0. Printed pitch variants: Red 1, Yellow 2,
// Blue 3. Printed defense: Red 4, Yellow 3, Blue 2.
//
// Text: "You may put a card from your hand on the bottom of your deck. If you do, draw a card."

package unplayable

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

func (SinkBelowRed) Play(s *sim.TurnState, l sim.Logger, self *sim.CardState) { l.Log(self, 0) }

func (SinkBelowYellow) Play(s *sim.TurnState, l sim.Logger, self *sim.CardState) { l.Log(self, 0) }

func (SinkBelowBlue) Play(s *sim.TurnState, l sim.Logger, self *sim.CardState) { l.Log(self, 0) }
