// Healing Potion — Generic Action - Item. Cost 0. Printed pitch variants: Blue 3.
//
// Text: "**Action** - Destroy this: Gain 2{h}. **Go again**"

package unplayable

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

func (HealingPotionBlue) Play(s *sim.TurnState, l sim.Logger, self *sim.CardState) { l.Log(self, 0) }
