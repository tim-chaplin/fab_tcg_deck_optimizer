// Destructive Deliberation — "When this hits a hero, create a Ponder token."

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
)

func destructiveDeliberationPlay(g card.GameEngine, l card.Logger, self *card.CardState) {
	self.RegisterOnHit(destructiveDeliberationOnHit)
}

func destructiveDeliberationOnHit(g card.GameEngine, l card.Logger, self *card.CardState, _ *card.OnHitHandler) {
	g.CreatePonders(1)
	l.AppendPostTrigger(self.Card.DisplayName(), "On-hit created a ponder", 0)
}

func (DestructiveDeliberationRed) Play(g card.GameEngine, l card.Logger, self *card.CardState) {
	destructiveDeliberationPlay(g, l, self)
}

func (DestructiveDeliberationYellow) Play(g card.GameEngine, l card.Logger, self *card.CardState) {
	destructiveDeliberationPlay(g, l, self)
}

func (DestructiveDeliberationBlue) Play(g card.GameEngine, l card.Logger, self *card.CardState) {
	destructiveDeliberationPlay(g, l, self)
}
