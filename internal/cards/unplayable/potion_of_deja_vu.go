// Potion of Déjà Vu — Generic Action - Item. Cost 0. Printed pitch variants: Blue 3.
//
// Text: "**Instant** - Destroy Potion of Déjà Vu: Put all cards from your pitch zone on top of your
// deck in any order."

package unplayable

import (
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
)

func (PotionOfDejaVuBlue) Play(s card.GameEngine, l card.Logger, self *card.CardState) {}
