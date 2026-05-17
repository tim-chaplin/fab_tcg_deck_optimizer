// Money or Your Life? — Generic Action - Attack. Cost 3.
// Text: "When this hits a hero, deal 2 damage to them unless they give you a Gold token
// they control. If you are a Thief, repeat this process once."
//
// Opponent items aren't modelled, so the "give a Gold token" branch never fires and the
// rider resolves as +2 on hit. The Thief repeat doubles the rider to +4 when the current
// hero's Types include TypeThief.

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

func moneyOrYourLifeOnHit(ge card.GameEngine, l card.Logger, self *card.CardState, _ *card.OnHitHandler) {
	n := 2
	if ge.CurrentHeroClass() == card.TypeThief {
		n = 4
	}
	ge.AddValue(n)
	l.AppendPostTriggerf(self.Card.DisplayName(), n, "On-hit dealt %d (opponent surrendered no Gold)", n)
}

func moneyOrYourLifePlay(ge card.GameEngine, l card.Logger, self *card.CardState) {
	self.RegisterOnHit(moneyOrYourLifeOnHit)
}

func (MoneyOrYourLifeRed) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	moneyOrYourLifePlay(ge, l, self)
}

func (MoneyOrYourLifeYellow) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	moneyOrYourLifePlay(ge, l, self)
}

func (MoneyOrYourLifeBlue) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	moneyOrYourLifePlay(ge, l, self)
}
