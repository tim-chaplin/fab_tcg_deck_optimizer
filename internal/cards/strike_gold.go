// Strike Gold — Generic Action - Attack.
// Text: "When this hits, create a Gold token."

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
)

func strikeGoldOnHit(g card.GameEngine, l card.Logger, self *card.CardState, _ *card.OnHitHandler) {
	g.CreateGold(1)
	l.AppendPostTrigger(self.Card.DisplayName(), "On-hit created a gold token", 0)
}

func strikeGoldPlay(g card.GameEngine, l card.Logger, self *card.CardState) {
	self.RegisterOnHit(strikeGoldOnHit)
}

func (StrikeGoldRed) Play(g card.GameEngine, l card.Logger, self *card.CardState) {
	strikeGoldPlay(g, l, self)
}

func (StrikeGoldYellow) Play(g card.GameEngine, l card.Logger, self *card.CardState) {
	strikeGoldPlay(g, l, self)
}

func (StrikeGoldBlue) Play(g card.GameEngine, l card.Logger, self *card.CardState) {
	strikeGoldPlay(g, l, self)
}
