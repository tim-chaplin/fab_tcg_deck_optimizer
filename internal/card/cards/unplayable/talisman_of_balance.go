// Talisman of Balance — Generic Action - Item. Cost 0. Printed pitch variants: Blue 3.
//
// Text: "**Go again** At the beginning of your end phase, if you have less cards in arsenal than an
// opposing hero, destroy Talisman of Balance and put the top card of your deck into an empty
// arsenal zone you control."

package unplayable

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

func (TalismanOfBalanceBlue) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
}
