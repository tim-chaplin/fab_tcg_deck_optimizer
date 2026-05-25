// Clap 'Em in Irons — Generic Action - Item. Cost 0. Printed pitch variants: Blue 3.
//
// Text: "**Go again** When this enters the arena, {t} target Pirate hero or ally. It can't {u}
// while this is in the arena. At the start of your turn, destroy this."

package unplayable

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

func (ClapEmInIronsBlue) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {}
