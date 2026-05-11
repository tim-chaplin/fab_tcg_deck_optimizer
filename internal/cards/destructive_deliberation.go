// Destructive Deliberation — "When this hits a hero, create a Ponder token."

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
)

func destructiveDeliberationPlay(s card.GameEngine, l card.Logger, self *card.CardState) {
	self.RegisterOnHit(destructiveDeliberationOnHit)
}

func destructiveDeliberationOnHit(s card.GameEngine, l card.Logger, self *card.CardState, _ *card.OnHitHandler) {
	s.CreatePonder(1)
	l.AppendPostTrigger(self.Card.DisplayName(), "On-hit created a ponder", 0)
}

func (DestructiveDeliberationRed) Play(s card.GameEngine, l card.Logger, self *card.CardState) {
	destructiveDeliberationPlay(s, l, self)
}

func (DestructiveDeliberationYellow) Play(s card.GameEngine, l card.Logger, self *card.CardState) {
	destructiveDeliberationPlay(s, l, self)
}

func (DestructiveDeliberationBlue) Play(s card.GameEngine, l card.Logger, self *card.CardState) {
	destructiveDeliberationPlay(s, l, self)
}
