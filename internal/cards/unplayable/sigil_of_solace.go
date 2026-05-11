// Sigil of Solace — Generic Instant. Cost 0. Printed pitch variants: Red 1, Yellow 2, Blue 3.
//
// Text: "Gain 3{h}"

package unplayable

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

func (SigilOfSolaceRed) Play(s *sim.TurnState, l sim.Logger, self *sim.CardState) {}

func (SigilOfSolaceYellow) Play(s *sim.TurnState, l sim.Logger, self *sim.CardState) {}

func (SigilOfSolaceBlue) Play(s *sim.TurnState, l sim.Logger, self *sim.CardState) {}
