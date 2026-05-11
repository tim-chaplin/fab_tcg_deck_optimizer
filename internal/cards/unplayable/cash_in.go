// Cash In — Generic Action. Cost 4, Pitch 2, Defense 2. Only printed in Yellow.
//
// Text: "You may destroy 4 Coppers, 2 Silvers, or 1 Gold you control rather than pay Cash In's {r}
// cost. Draw 2 cards. **Go again**"

package unplayable

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

func (CashInYellow) Play(s *sim.TurnState, l sim.Logger, self *sim.CardState) {}
