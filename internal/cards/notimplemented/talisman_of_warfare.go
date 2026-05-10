// Talisman of Warfare — Generic Action - Item. Cost 0. Printed pitch variants: Yellow 2.
//
// Text: "**Go again** When a source you control deals exactly 2 damage to an opposing hero, destroy
// Talisman of Warfare and all cards in all arsenals."

package notimplemented

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

// not implemented: self-destroys + wipes all arsenals on a 2-damage hit

func (TalismanOfWarfareYellow) Play(s *sim.TurnState, l sim.Logger, self *sim.CardState) {
	l.Log(self, 0)
}
