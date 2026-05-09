// Muscle Mutt — Generic Action - Attack. Cost 3, Pitch 2, Power 6, Defense 2. Only printed in
// Yellow.

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

func (c MuscleMuttYellow) Play(s *sim.TurnState, self *sim.CardState) {
	n := self.DealEffectiveAttack(s)
	s.Log(self, n)
}
