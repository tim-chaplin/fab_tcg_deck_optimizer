// Outed — Generic Action - Attack. Cost 0, Pitch 1, Power 3, Defense 0. Only printed in Red.
//
// Text: "If you are **marked**, you can't play this. If the defending hero is **marked**, this gets
// +1{p}. **Go again**"

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

func (OutedRed) Play(s *sim.TurnState, self *sim.CardState) {
	if s.OpponentMarked {
		self.BonusAttack++
	}
	n := self.DealEffectiveAttack(s)
	s.Log(self, n)
}
