// Destructive Deliberation — "When this hits a hero, create a Ponder token."

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
)

func destructiveDeliberationPlay(ge card.GameEngine, l card.Logger, self *card.CardState) {
	self.RegisterOnHit(destructiveDeliberationOnHit)
}

func destructiveDeliberationOnHit(ge card.GameEngine, l card.Logger, self *card.CardState, _ *card.OnHitHandler) {
	ge.CreatePonders(1)
	l.AppendPostTrigger(self.Card.DisplayName(), "On-hit created a ponder", 0)
}

func (DestructiveDeliberationRed) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	destructiveDeliberationPlay(ge, l, self)
}

func (DestructiveDeliberationYellow) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	destructiveDeliberationPlay(ge, l, self)
}

func (DestructiveDeliberationBlue) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	destructiveDeliberationPlay(ge, l, self)
}
