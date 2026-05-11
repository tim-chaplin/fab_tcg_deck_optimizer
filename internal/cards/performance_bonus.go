// Performance Bonus — Generic Action - Attack.
// Text: "When this hits, create a Gold token. If this was played from arsenal, it gets
// **Go again**."

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
)

func performanceBonusOnHit(s card.GameEngine, l card.Logger, self *card.CardState, _ *card.OnHitHandler) {
	s.CreateGold(1)
	l.AppendPostTrigger(self.Card.DisplayName(), "On-hit created a gold token", 0)
}

func performanceBonusPlay(s card.GameEngine, l card.Logger, self *card.CardState) {
	self.GrantGoAgainIfFromArsenal()
	self.RegisterOnHit(performanceBonusOnHit)
}

func (PerformanceBonusRed) Play(s card.GameEngine, l card.Logger, self *card.CardState) {
	performanceBonusPlay(s, l, self)
}

func (PerformanceBonusYellow) Play(s card.GameEngine, l card.Logger, self *card.CardState) {
	performanceBonusPlay(s, l, self)
}

func (PerformanceBonusBlue) Play(s card.GameEngine, l card.Logger, self *card.CardState) {
	performanceBonusPlay(s, l, self)
}
