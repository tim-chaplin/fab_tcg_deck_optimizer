// Energy Potion — Generic Action - Item. Cost 0. Printed pitch variants: Blue 3.
//
// Text: "**Instant** - Destroy this: Gain {r}{r}"

package unplayable

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

func (EnergyPotionBlue) Play(s *sim.TurnState, self *sim.CardState) { s.Log(self, 0) }
