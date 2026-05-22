// Talisman of Tithes — Generic Action - Item. Cost 0. Printed pitch variants: Blue 3.
//
// Text: "**Go again** If an opponent would draw 1 or more cards during your action phase, instead
// destroy Talisman of Tithes and they draw that many cards minus 1."

package unplayable

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

func (TalismanOfTithesBlue) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {}
