// Starting Stake — Generic Action. Cost 0, Pitch 2, Defense 3. Yellow only.
// Text: "If you control no Gold tokens, create a Gold token."

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

func (StartingStakeYellow) Play(s *sim.TurnState, self *sim.CardState) {
	if s.Gold() == 0 {
		s.CreateGold(1)
		s.LogRider(self, 0, "Created a gold token")
	}
	s.Log(self, 0)
}
