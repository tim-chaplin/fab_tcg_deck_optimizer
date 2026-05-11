// Put in Context — Generic Defense Reaction. Cost 0, Pitch 3, Defense 3. Only printed in Blue.
// Text: "This can only defend an attack with 3 or less base {p}."

package notimplemented

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

// not implemented: base-power cap on what this can defend is ignored; treated as legal vs every
// attack

func (PutInContextBlue) Play(s *sim.TurnState, l sim.Logger, self *sim.CardState) {
	n := self.DealEffectiveDefense(s)
	self.Log(l, n)
}
