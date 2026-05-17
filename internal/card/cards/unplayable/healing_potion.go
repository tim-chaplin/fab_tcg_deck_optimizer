// Healing Potion — Generic Action - Item. Cost 0. Printed pitch variants: Blue 3.
//
// Text: "**Action** - Destroy this: Gain 2{h}. **Go again**"

package unplayable

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

func (HealingPotionBlue) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {}
