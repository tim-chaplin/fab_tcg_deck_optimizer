// Destructive Deliberation — "When this hits a hero, create a Ponder token."

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

func destructiveDeliberationPlay(s sim.GameEngine, l sim.Logger, self *sim.CardState) {
	self.RegisterOnHit(destructiveDeliberationOnHit)
}

func destructiveDeliberationOnHit(s sim.GameEngine, l sim.Logger, self *sim.CardState, _ *sim.OnHitHandler) {
	s.CreatePonder(1)
	l.AppendPostTrigger(self.Card.DisplayName(), "On-hit created a ponder", 0)
}

func (DestructiveDeliberationRed) Play(s sim.GameEngine, l sim.Logger, self *sim.CardState) {
	destructiveDeliberationPlay(s, l, self)
}

func (DestructiveDeliberationYellow) Play(s sim.GameEngine, l sim.Logger, self *sim.CardState) {
	destructiveDeliberationPlay(s, l, self)
}

func (DestructiveDeliberationBlue) Play(s sim.GameEngine, l sim.Logger, self *sim.CardState) {
	destructiveDeliberationPlay(s, l, self)
}
