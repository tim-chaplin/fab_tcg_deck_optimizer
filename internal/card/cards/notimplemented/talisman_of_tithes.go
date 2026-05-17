// Talisman of Tithes — Generic Action - Item. Cost 0. Printed pitch variants: Blue 3.
//
// Text: "**Go again** If an opponent would draw 1 or more cards during your action phase, instead
// destroy Talisman of Tithes and they draw that many cards minus 1."

package notimplemented

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

// not implemented: self-destroys on an opposing draw during your action phase → opponent draws
// minus 1

func (TalismanOfTithesBlue) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {}
