// Strike Gold — Generic Action - Attack.
// Text: "When this hits, create a Gold token."

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
)

func strikeGoldOnHit(s card.GameEngine, l card.Logger, self *card.CardState, _ *card.OnHitHandler) {
	s.CreateGold(1)
	l.AppendPostTrigger(self.Card.DisplayName(), "On-hit created a gold token", 0)
}

func strikeGoldPlay(s card.GameEngine, l card.Logger, self *card.CardState) {
	self.RegisterOnHit(strikeGoldOnHit)
}

func (StrikeGoldRed) Play(s card.GameEngine, l card.Logger, self *card.CardState) {
	strikeGoldPlay(s, l, self)
}

func (StrikeGoldYellow) Play(s card.GameEngine, l card.Logger, self *card.CardState) {
	strikeGoldPlay(s, l, self)
}

func (StrikeGoldBlue) Play(s card.GameEngine, l card.Logger, self *card.CardState) {
	strikeGoldPlay(s, l, self)
}
