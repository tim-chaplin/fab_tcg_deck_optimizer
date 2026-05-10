// Public Bounty — Generic Action. Cost 1. Printed pitch variants: Red 1, Yellow 2, Blue 3. Defense
// 2.
//
// Text: "**Mark** target opposing hero. The next time you attack a **marked** hero this turn, the
// attack gets +N{p}. **Go again**" (Red N=3, Yellow N=2, Blue N=1.)

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

func publicBountyPlay(s *sim.TurnState, l sim.Logger, self *sim.CardState, n int) {
	s.OpponentMarked = true
	GrantNextCardBonusAttack(s, n, IsAttack)
	l.Log(self, 0)
}

func (PublicBountyRed) Play(s *sim.TurnState, l sim.Logger, self *sim.CardState) {
	publicBountyPlay(s, l, self, 3)
}

func (PublicBountyYellow) Play(s *sim.TurnState, l sim.Logger, self *sim.CardState) {
	publicBountyPlay(s, l, self, 2)
}

func (PublicBountyBlue) Play(s *sim.TurnState, l sim.Logger, self *sim.CardState) {
	publicBountyPlay(s, l, self, 1)
}
