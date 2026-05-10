// Ransack and Raze — Generic Action. Cost X, Pitch 3, Defense 3. Only printed in Blue.
//
// Text: "Destroy target landmark with cost X. Create X Gold tokens. **Go again**"
//
// X cost treated as 0.

package notimplemented

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

// not implemented: gold tokens, landmarks

func (RansackAndRazeBlue) Play(s *sim.TurnState, l sim.Logger, self *sim.CardState) { l.Log(self, 0) }
