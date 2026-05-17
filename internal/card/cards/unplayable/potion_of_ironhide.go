// Potion of Ironhide — Generic Action - Item. Cost 0. Printed pitch variants: Blue 3.
//
// Text: "**Instant** - Destroy Potion of Ironhide: Attack action cards you own gain +1{d} this
// turn."

package unplayable

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

func (PotionOfIronhideBlue) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {}
