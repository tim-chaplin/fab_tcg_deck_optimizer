// Talisman of Warfare — Generic Action - Item. Cost 0. Printed pitch variants: Yellow 2.
//
// Text: "**Go again** When a source you control deals exactly 2 damage to an opposing hero, destroy
// Talisman of Warfare and all cards in all arsenals."

package notimplemented

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

// not implemented: self-destroys + wipes all arsenals on a 2-damage hit

func (TalismanOfWarfareYellow) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
}
