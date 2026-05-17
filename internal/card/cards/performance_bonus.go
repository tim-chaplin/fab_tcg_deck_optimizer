// Performance Bonus — Generic Action - Attack.
// Text: "When this hits, create a Gold token. If this was played from arsenal, it gets
// **Go again**."

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

func performanceBonusOnHit(ge card.GameEngine, l card.Logger, self *card.CardState, _ *card.OnHitHandler) {
	ge.CreateGold(1)
	l.AppendPostTrigger(self.Card.DisplayName(), "On-hit created a gold token", 0)
}

func performanceBonusPlay(ge card.GameEngine, l card.Logger, self *card.CardState) {
	self.GrantGoAgainIfFromArsenal()
	self.RegisterOnHit(performanceBonusOnHit)
}

func (PerformanceBonusRed) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	performanceBonusPlay(ge, l, self)
}

func (PerformanceBonusYellow) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	performanceBonusPlay(ge, l, self)
}

func (PerformanceBonusBlue) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	performanceBonusPlay(ge, l, self)
}
