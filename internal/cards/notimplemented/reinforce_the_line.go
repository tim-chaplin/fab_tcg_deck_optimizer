// Reinforce the Line — Generic Instant. Cost 0. Printed pitch variants: Red 1, Yellow 2, Blue 3.
//
// Text: "Target defending attack action card gains +4{d}."

package notimplemented

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

// not implemented: Instant +N{d} grant to a defending attack action card

func (ReinforceTheLineRed) Play(s *sim.TurnState, l sim.Logger, self *sim.CardState) { l.Log(self, 0) }

// not implemented: Instant +N{d} grant to a defending attack action card

func (ReinforceTheLineYellow) Play(s *sim.TurnState, l sim.Logger, self *sim.CardState) {
	l.Log(self, 0)
}

// not implemented: Instant +N{d} grant to a defending attack action card

func (ReinforceTheLineBlue) Play(s *sim.TurnState, l sim.Logger, self *sim.CardState) { l.Log(self, 0) }
