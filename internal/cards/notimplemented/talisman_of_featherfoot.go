// Talisman of Featherfoot — Generic Action - Item. Cost 0. Printed pitch variants: Yellow 2.
//
// Text: "**Go again** When an attack you control gains exactly +1{p} from an effect during the
// reaction step, destroy Talisman of Featherfoot and the attack gains **go again**."

package notimplemented

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

// not implemented: self-destroys when an attack gains exactly +1{p} in the reaction step →
// grants go again

func (TalismanOfFeatherfootYellow) Play(s *sim.TurnState, self *sim.CardState) { s.Log(self, 0) }
