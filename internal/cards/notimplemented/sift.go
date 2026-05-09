// Sift — Generic Action. Cost 0. Printed pitch variants: Red 1, Yellow 2, Blue 3. Defense 3.
//
// Text: "Put up to 4 cards from your hand on the bottom of your deck, then draw that many cards.
// **Go again**"

package notimplemented

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

// not implemented: hand cycling

func (SiftRed) Play(s *sim.TurnState, self *sim.CardState) { s.Log(self, 0) }

// not implemented: hand cycling

func (SiftYellow) Play(s *sim.TurnState, self *sim.CardState) { s.Log(self, 0) }

// not implemented: hand cycling

func (SiftBlue) Play(s *sim.TurnState, self *sim.CardState) { s.Log(self, 0) }
