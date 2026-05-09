// Performance Bonus — Generic Action - Attack.
// Text: "When this hits, create a Gold token. If this was played from arsenal, it gets
// **Go again**."

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

func performanceBonusOnHit(s *sim.TurnState, self *sim.CardState, _ *sim.OnHitHandler) {
	s.CreateGold(1)
	s.LogRider(self, 0, "On-hit created a gold token")
}

func performanceBonusPlay(s *sim.TurnState, self *sim.CardState) {
	self.GrantGoAgainIfFromArsenal()
	n := self.DealEffectiveAttack(s)
	s.Log(self, n)
	self.RegisterOnHit(performanceBonusOnHit)
}

func (PerformanceBonusRed) Play(s *sim.TurnState, self *sim.CardState) {
	performanceBonusPlay(s, self)
}

func (PerformanceBonusYellow) Play(s *sim.TurnState, self *sim.CardState) {
	performanceBonusPlay(s, self)
}

func (PerformanceBonusBlue) Play(s *sim.TurnState, self *sim.CardState) {
	performanceBonusPlay(s, self)
}
