// Amulet of Ignition — Generic Action - Item. Cost 0. Printed pitch variants: Yellow 2.
//
// Text: "**Go again** **Instant** - Destroy Amulet of Ignition: The next ability you activate this
// turn costs {r} less. Activate this ability only if you haven't played a card or activated an
// ability this turn."

package unplayable

import (
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
)

func (AmuletOfIgnitionYellow) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
}
