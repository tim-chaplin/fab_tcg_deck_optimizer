// Wage Gold — Generic Action - Attack. Cost 3.
// Text: "**Universal** When this attacks a hero, you may **wager** a Gold token with them."
//
// Types(g) ORs in card.NewTypeSet(g.CurrentHeroClass()) so the Universal keyword grants
// the active hero's class to the type-line for class-gated triggers. The "may" wager
// opts in only when likely-to-hit; the win (a Gold token) resolves on hit.

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
)

func wageGoldOnHit(g card.GameEngine, l card.Logger, self *card.CardState, _ *card.OnHitHandler) {
	g.CreateGold(1)
	l.AppendPostTrigger(self.Card.DisplayName(), "On-hit won wager", 0)
}

func wageGoldPlay(g card.GameEngine, l card.Logger, self *card.CardState) {
	self.RegisterOnHit(wageGoldOnHit)
}

func (WageGoldRed) Play(g card.GameEngine, l card.Logger, self *card.CardState) {
	wageGoldPlay(g, l, self)
}

func (WageGoldYellow) Play(g card.GameEngine, l card.Logger, self *card.CardState) {
	wageGoldPlay(g, l, self)
}

func (WageGoldBlue) Play(g card.GameEngine, l card.Logger, self *card.CardState) {
	wageGoldPlay(g, l, self)
}
