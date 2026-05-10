// Tit for Tat — Generic Action. Cost 0, Pitch 3, Defense 2. Only printed in Blue.
//
// Text: "{t} target hero. {u} another target hero. **Go again**"

package notimplemented

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

// not implemented: freeze/unfreeze

func (TitForTatBlue) Play(s *sim.TurnState, l sim.Logger, self *sim.CardState) { l.Log(self, 0) }
