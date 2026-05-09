// On a Knife Edge — Generic Action. Cost 0, Pitch 2, Defense 2. Only printed in Yellow.
//
// Text: "Your next sword attack this turn gains **go again**. **Go again**"

package notimplemented

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

// not implemented: next-sword-attack go-again grant (weapon chain not scanned)

func (OnAKnifeEdgeYellow) Play(s *sim.TurnState, self *sim.CardState) { s.Log(self, 0) }
