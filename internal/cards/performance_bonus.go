// Performance Bonus — Generic Action - Attack.
// Text: "When this hits, create a Gold token. If this was played from arsenal, it gets
// **Go again**."

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

func performanceBonusOnHit(s *sim.TurnState, l sim.Logger, self *sim.CardState, _ *sim.OnHitHandler) {
	s.CreateGold(1)
	self.LogRider(l, 0, "On-hit created a gold token")
}

func performanceBonusPlay(s *sim.TurnState, l sim.Logger, self *sim.CardState) {
	self.GrantGoAgainIfFromArsenal()
	n := self.DealEffectiveAttack(s)
	self.Log(l, n)
	self.RegisterOnHit(performanceBonusOnHit)
}

func (PerformanceBonusRed) Play(s *sim.TurnState, l sim.Logger, self *sim.CardState) {
	performanceBonusPlay(s, l, self)
}

func (PerformanceBonusYellow) Play(s *sim.TurnState, l sim.Logger, self *sim.CardState) {
	performanceBonusPlay(s, l, self)
}

func (PerformanceBonusBlue) Play(s *sim.TurnState, l sim.Logger, self *sim.CardState) {
	performanceBonusPlay(s, l, self)
}
