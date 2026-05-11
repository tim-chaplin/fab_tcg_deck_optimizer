// Cracked Bauble — Generic Resource. Cost 0. Printed pitch variants: Yellow 2.
//
// Text: "*(A player may add any number of Cracked Baubles to their card-pool in sealed deck or
// booster draft formats.)*"

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

func (CrackedBaubleYellow) Play(s *sim.TurnState, l sim.Logger, self *sim.CardState) { self.Log(l, 0) }
