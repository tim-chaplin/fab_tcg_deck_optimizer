// Talisman of Recompense — Generic Action - Item. Cost 0. Printed pitch variants: Yellow 2.
//
// Text: "**Go again** Whenever you pitch a card, if you would gain exactly one {r}, instead destroy
// Talisman of Recompense and gain {r}{r}{r}."

package notimplemented

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

// not implemented: self-destroys on pitching a 1-resource card → gain {r}{r}{r} instead

func (TalismanOfRecompenseYellow) Play(s *sim.TurnState, self *sim.CardState) { s.Log(self, 0) }
