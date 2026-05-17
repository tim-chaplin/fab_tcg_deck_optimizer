// Potion of Strength — Generic Action - Item. Cost 0. Printed pitch variants: Blue 3.
//
// Text: "**Action** - Destroy this: Your next attack this turn gains +2{p}. **Go again**"

package unplayable

import (
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
)

func (PotionOfStrengthBlue) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {}
