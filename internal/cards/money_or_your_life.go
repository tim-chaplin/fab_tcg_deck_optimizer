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

	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

func moneyOrYourLifeOnHit(s *sim.TurnState, l sim.Logger, self *sim.CardState, _ *sim.OnHitHandler) {
	n := 2
	if sim.CurrentHero != nil && sim.CurrentHero.Types().Has(card.TypeThief) {
		n = 4
	}
	s.AddValue(n)
	self.LogRiderf(l, n, "On-hit dealt %d (opponent surrendered no Gold)", n)
}

func moneyOrYourLifePlay(s *sim.TurnState, l sim.Logger, self *sim.CardState) {
	n := self.DealEffectiveAttack(s)
	self.Log(l, n)
	self.RegisterOnHit(moneyOrYourLifeOnHit)
}

func (MoneyOrYourLifeRed) Play(s *sim.TurnState, l sim.Logger, self *sim.CardState) {
	moneyOrYourLifePlay(s, l, self)
}

func (MoneyOrYourLifeYellow) Play(s *sim.TurnState, l sim.Logger, self *sim.CardState) {
	moneyOrYourLifePlay(s, l, self)
}

func (MoneyOrYourLifeBlue) Play(s *sim.TurnState, l sim.Logger, self *sim.CardState) {
	moneyOrYourLifePlay(s, l, self)
}
