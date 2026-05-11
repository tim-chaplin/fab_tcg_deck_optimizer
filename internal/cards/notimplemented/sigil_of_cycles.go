// Sigil of Cycles — Generic Action - Aura. Cost 0, Pitch 3, Defense 2. Only printed in Blue.
//
// Text: "**Go again** At the beginning of your action phase, destroy this. When this leaves the
// arena, discard a card then draw a card."

package notimplemented

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

// not implemented: start-of-action-phase self-destroy, leaves-arena discard/draw

func (SigilOfCyclesBlue) Play(s *sim.TurnState, l sim.Logger, self *sim.CardState) {}
