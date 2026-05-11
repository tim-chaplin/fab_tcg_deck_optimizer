// Potion of Strength — Generic Action - Item. Cost 0. Printed pitch variants: Blue 3.
//
// Text: "**Action** - Destroy this: Your next attack this turn gains +2{p}. **Go again**"

package unplayable

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

func (PotionOfStrengthBlue) Play(s sim.GameEngine, l sim.Logger, self *sim.CardState) {}
