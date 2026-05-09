// Destructive Deliberation — "When this hits a hero, create a Ponder token."

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

func destructiveDeliberationPlay(s *sim.TurnState, self *sim.CardState) {
	n := self.DealEffectiveAttack(s)
	s.Log(self, n)
	self.RegisterOnHit(destructiveDeliberationOnHit)
}

func destructiveDeliberationOnHit(s *sim.TurnState, self *sim.CardState, _ *sim.OnHitHandler) {
	s.CreatePonder(1)
	s.LogRider(self, 0, "On-hit created a ponder")
}

func (DestructiveDeliberationRed) Play(s *sim.TurnState, self *sim.CardState) {
	destructiveDeliberationPlay(s, self)
}

func (DestructiveDeliberationYellow) Play(s *sim.TurnState, self *sim.CardState) {
	destructiveDeliberationPlay(s, self)
}

func (DestructiveDeliberationBlue) Play(s *sim.TurnState, self *sim.CardState) {
	destructiveDeliberationPlay(s, self)
}
