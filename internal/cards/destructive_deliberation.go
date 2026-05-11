// Destructive Deliberation — "When this hits a hero, create a Ponder token."

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

func destructiveDeliberationPlay(s *sim.TurnState, l sim.Logger, self *sim.CardState) {
	n := self.DealEffectiveAttack(s)
	self.Log(l, n)
	self.RegisterOnHit(destructiveDeliberationOnHit)
}

func destructiveDeliberationOnHit(s *sim.TurnState, l sim.Logger, self *sim.CardState, _ *sim.OnHitHandler) {
	s.CreatePonder(1)
	self.LogRider(l, 0, "On-hit created a ponder")
}

func (DestructiveDeliberationRed) Play(s *sim.TurnState, l sim.Logger, self *sim.CardState) {
	destructiveDeliberationPlay(s, l, self)
}

func (DestructiveDeliberationYellow) Play(s *sim.TurnState, l sim.Logger, self *sim.CardState) {
	destructiveDeliberationPlay(s, l, self)
}

func (DestructiveDeliberationBlue) Play(s *sim.TurnState, l sim.Logger, self *sim.CardState) {
	destructiveDeliberationPlay(s, l, self)
}
