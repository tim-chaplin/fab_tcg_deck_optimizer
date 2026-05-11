// Titanium Bauble — Generic Resource. Cost 0. Printed pitch variants: Blue 3. Defense 3.

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

func (TitaniumBaubleBlue) Play(s *sim.TurnState, l sim.Logger, self *sim.CardState) { self.Log(l, 0) }
