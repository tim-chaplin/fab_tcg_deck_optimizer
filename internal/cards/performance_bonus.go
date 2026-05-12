// Performance Bonus — Generic Action - Attack.
// Text: "When this hits, create a Gold token. If this was played from arsenal, it gets
// **Go again**."

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
)

func performanceBonusOnHit(g card.GameEngine, l card.Logger, self *card.CardState, _ *card.OnHitHandler) {
	g.CreateGold(1)
	l.AppendPostTrigger(self.Card.DisplayName(), "On-hit created a gold token", 0)
}

func performanceBonusPlay(g card.GameEngine, l card.Logger, self *card.CardState) {
	self.GrantGoAgainIfFromArsenal()
	self.RegisterOnHit(performanceBonusOnHit)
}

func (PerformanceBonusRed) Play(g card.GameEngine, l card.Logger, self *card.CardState) {
	performanceBonusPlay(g, l, self)
}

func (PerformanceBonusYellow) Play(g card.GameEngine, l card.Logger, self *card.CardState) {
	performanceBonusPlay(g, l, self)
}

func (PerformanceBonusBlue) Play(g card.GameEngine, l card.Logger, self *card.CardState) {
	performanceBonusPlay(g, l, self)
}
