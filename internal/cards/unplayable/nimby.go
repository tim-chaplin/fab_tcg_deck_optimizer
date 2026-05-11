// Nimby — Generic Action - Attack. Cost 0. Printed power: Red 3, Yellow 2, Blue 1. Printed pitch
// variants: Red 1, Yellow 2, Blue 3. Defense 2.
//
// Text: "When this attacks, you may search your deck for a Nimblism, reveal it, put it into your
// hand, then shuffle."

package unplayable

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

func (NimbyRed) Play(s *sim.TurnState, l sim.Logger, self *sim.CardState) {}

func (NimbyYellow) Play(s *sim.TurnState, l sim.Logger, self *sim.CardState) {}

func (NimbyBlue) Play(s *sim.TurnState, l sim.Logger, self *sim.CardState) {}
