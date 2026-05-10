// Fool's Gold — Generic Resource. Cost 0. Printed pitch variants: Yellow 2.
//
// Text: "When this is discarded, create a Gold token."

package notimplemented

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

// not implemented: discard trigger creates a Gold token

func (FoolsGoldYellow) Play(s *sim.TurnState, l sim.Logger, self *sim.CardState) { l.Log(self, 0) }
