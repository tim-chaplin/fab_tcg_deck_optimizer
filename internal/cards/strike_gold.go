// Strike Gold — Generic Action - Attack.
// Text: "When this hits, create a Gold token."

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

func strikeGoldOnHit(s sim.GameEngine, l sim.Logger, self *sim.CardState, _ *sim.OnHitHandler) {
	s.CreateGold(1)
	l.AppendPostTrigger(self.Card.DisplayName(), "On-hit created a gold token", 0)
}

func strikeGoldPlay(s sim.GameEngine, l sim.Logger, self *sim.CardState) {
	self.RegisterOnHit(strikeGoldOnHit)
}

func (StrikeGoldRed) Play(s sim.GameEngine, l sim.Logger, self *sim.CardState) {
	strikeGoldPlay(s, l, self)
}

func (StrikeGoldYellow) Play(s sim.GameEngine, l sim.Logger, self *sim.CardState) {
	strikeGoldPlay(s, l, self)
}

func (StrikeGoldBlue) Play(s sim.GameEngine, l sim.Logger, self *sim.CardState) {
	strikeGoldPlay(s, l, self)
}
