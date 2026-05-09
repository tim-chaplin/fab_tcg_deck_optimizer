// Wage Gold — Generic Action - Attack. Cost 3.
// Text: "**Universal** When this attacks a hero, you may **wager** a Gold token with them."
//
// Types() ORs in sim.Universal() so the keyword applies. The "may" wager opts in only when
// likely-to-hit; the win (a Gold token) resolves on hit.

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

func wageGoldOnHit(s *sim.TurnState, self *sim.CardState, _ *sim.OnHitHandler) {
	s.CreateGold(1)
	s.LogRider(self, 0, "On-hit won wager")
}

func wageGoldPlay(s *sim.TurnState, self *sim.CardState) {
	n := self.DealEffectiveAttack(s)
	s.Log(self, n)
	self.RegisterOnHit(wageGoldOnHit)
}

func (WageGoldRed) Play(s *sim.TurnState, self *sim.CardState) { wageGoldPlay(s, self) }

func (WageGoldYellow) Play(s *sim.TurnState, self *sim.CardState) { wageGoldPlay(s, self) }

func (WageGoldBlue) Play(s *sim.TurnState, self *sim.CardState) { wageGoldPlay(s, self) }
