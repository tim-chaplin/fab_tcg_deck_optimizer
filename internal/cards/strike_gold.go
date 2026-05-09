// Strike Gold — Generic Action - Attack.
// Text: "When this hits, create a Gold token."

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

func strikeGoldOnHit(s *sim.TurnState, self *sim.CardState, _ *sim.OnHitHandler) {
	s.CreateGold(1)
	s.LogRider(self, 0, "On-hit created a gold token")
}

func strikeGoldPlay(s *sim.TurnState, self *sim.CardState) {
	n := self.DealEffectiveAttack(s)
	s.Log(self, n)
	self.RegisterOnHit(strikeGoldOnHit)
}

func (StrikeGoldRed) Play(s *sim.TurnState, self *sim.CardState) { strikeGoldPlay(s, self) }

func (StrikeGoldYellow) Play(s *sim.TurnState, self *sim.CardState) { strikeGoldPlay(s, self) }

func (StrikeGoldBlue) Play(s *sim.TurnState, self *sim.CardState) { strikeGoldPlay(s, self) }
