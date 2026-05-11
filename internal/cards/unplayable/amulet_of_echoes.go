// Amulet of Echoes — Generic Action - Item. Cost 0. Printed pitch variants: Blue 3.
//
// Text: "**Go again** **Instant** - Destroy Amulet of Echoes: Target hero discards 2 cards.
// Activate this ability only if they have played 2 or more cards with the same name this turn."

package unplayable

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

func (AmuletOfEchoesBlue) Play(s sim.GameEngine, l sim.Logger, self *sim.CardState) {}
