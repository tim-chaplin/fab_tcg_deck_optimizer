// Lead the Charge — Generic Action. Cost 0. Printed pitch variants: Red 1, Yellow 2, Blue 3.
// Defense 2.
//
// Text: "The next time you play an action card with cost 0 or greater this turn, gain 1 action
// point. **Go again**"

package notimplemented

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

// not implemented: action point grant

func (LeadTheChargeRed) Play(s *sim.TurnState, l sim.Logger, self *sim.CardState) { l.Log(self, 0) }

// not implemented: action point grant

func (LeadTheChargeYellow) Play(s *sim.TurnState, l sim.Logger, self *sim.CardState) { l.Log(self, 0) }

// not implemented: action point grant

func (LeadTheChargeBlue) Play(s *sim.TurnState, l sim.Logger, self *sim.CardState) { l.Log(self, 0) }
