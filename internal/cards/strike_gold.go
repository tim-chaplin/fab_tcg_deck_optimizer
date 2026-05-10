// Strike Gold — Generic Action - Attack.
// Text: "When this hits, create a Gold token."

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

func strikeGoldOnHit(s *sim.TurnState, l sim.Logger, self *sim.CardState, _ *sim.OnHitHandler) {
	s.CreateGold(1)
	l.LogRider(self, 0, "On-hit created a gold token")
}

func strikeGoldPlay(s *sim.TurnState, l sim.Logger, self *sim.CardState) {
	n := self.DealEffectiveAttack(s)
	l.Log(self, n)
	self.RegisterOnHit(strikeGoldOnHit)
}

func (StrikeGoldRed) Play(s *sim.TurnState, l sim.Logger, self *sim.CardState) {
	strikeGoldPlay(s, l, self)
}

func (StrikeGoldYellow) Play(s *sim.TurnState, l sim.Logger, self *sim.CardState) {
	strikeGoldPlay(s, l, self)
}

func (StrikeGoldBlue) Play(s *sim.TurnState, l sim.Logger, self *sim.CardState) {
	strikeGoldPlay(s, l, self)
}
