// Starting Stake — Generic Action. Cost 0, Pitch 2, Defense 3. Yellow only.
// Text: "If you control no Gold tokens, create a Gold token."

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

func (StartingStakeYellow) Play(s *sim.TurnState, l sim.Logger, self *sim.CardState) {
	if s.Gold() == 0 {
		s.CreateGold(1)
		l.AppendPostTrigger(self.Card.DisplayName(), "Created a gold token", 0)
	}
}
