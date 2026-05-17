// Strike Gold — Generic Action - Attack.
// Text: "When this hits, create a Gold token."

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
)

func strikeGoldOnHit(ge card.GameEngine, l card.Logger, self *card.CardState, _ *card.OnHitHandler) {
	ge.CreateGold(1)
	l.AppendPostTrigger(self.Card.DisplayName(), "On-hit created a gold token", 0)
}

func strikeGoldPlay(ge card.GameEngine, l card.Logger, self *card.CardState) {
	self.RegisterOnHit(strikeGoldOnHit)
}

func (StrikeGoldRed) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	strikeGoldPlay(ge, l, self)
}

func (StrikeGoldYellow) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	strikeGoldPlay(ge, l, self)
}

func (StrikeGoldBlue) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	strikeGoldPlay(ge, l, self)
}
