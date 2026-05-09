// Belittle — Generic Action - Attack. Cost 1. Printed power: Red 3, Yellow 2, Blue 1. Printed pitch
// variants: Red 1, Yellow 2, Blue 3. Defense 2.
//
// Text: "As an additional cost to play Belittle, you may reveal an attack action card with 3 or
// less base {p} from your hand. If you do, search your deck for a card named Minnowism, reveal it,
// put it into your hand, then shuffle your deck. **Go again**"

package unplayable

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

func (BelittleRed) Play(s *sim.TurnState, self *sim.CardState) { s.Log(self, 0) }

func (BelittleYellow) Play(s *sim.TurnState, self *sim.CardState) { s.Log(self, 0) }

func (BelittleBlue) Play(s *sim.TurnState, self *sim.CardState) { s.Log(self, 0) }
