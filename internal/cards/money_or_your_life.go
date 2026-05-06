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
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

var moneyOrYourLifeTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeAttack)

func moneyOrYourLifeOnHit(s *sim.TurnState, self *sim.CardState, _ *sim.OnHitHandler) {
	n := 2
	if sim.CurrentHero != nil && sim.CurrentHero.Types().Has(card.TypeThief) {
		n = 4
	}
	s.AddValue(n)
	s.LogRiderf(self, n, "On-hit dealt %d (opponent surrendered no Gold)", n)
}

func moneyOrYourLifePlay(s *sim.TurnState, self *sim.CardState) {
	n := self.DealEffectiveAttack(s)
	s.Log(self, n)
	self.RegisterOnHit(moneyOrYourLifeOnHit)
}

type MoneyOrYourLifeRed struct{}

func (MoneyOrYourLifeRed) ID() ids.CardID                             { return ids.MoneyOrYourLifeRed }
func (MoneyOrYourLifeRed) Name() string                               { return "Money or Your Life?" }
func (MoneyOrYourLifeRed) Cost(*sim.TurnState) int                    { return 3 }
func (MoneyOrYourLifeRed) Pitch() int                                 { return 1 }
func (MoneyOrYourLifeRed) Attack() int                                { return 6 }
func (MoneyOrYourLifeRed) Defense() int                               { return 2 }
func (MoneyOrYourLifeRed) Types() card.TypeSet                        { return moneyOrYourLifeTypes }
func (MoneyOrYourLifeRed) GoAgain() bool                              { return false }
func (MoneyOrYourLifeRed) Play(s *sim.TurnState, self *sim.CardState) { moneyOrYourLifePlay(s, self) }

type MoneyOrYourLifeYellow struct{}

func (MoneyOrYourLifeYellow) ID() ids.CardID          { return ids.MoneyOrYourLifeYellow }
func (MoneyOrYourLifeYellow) Name() string            { return "Money or Your Life?" }
func (MoneyOrYourLifeYellow) Cost(*sim.TurnState) int { return 3 }
func (MoneyOrYourLifeYellow) Pitch() int              { return 2 }
func (MoneyOrYourLifeYellow) Attack() int             { return 5 }
func (MoneyOrYourLifeYellow) Defense() int            { return 2 }
func (MoneyOrYourLifeYellow) Types() card.TypeSet     { return moneyOrYourLifeTypes }
func (MoneyOrYourLifeYellow) GoAgain() bool           { return false }
func (MoneyOrYourLifeYellow) Play(s *sim.TurnState, self *sim.CardState) {
	moneyOrYourLifePlay(s, self)
}

type MoneyOrYourLifeBlue struct{}

func (MoneyOrYourLifeBlue) ID() ids.CardID                             { return ids.MoneyOrYourLifeBlue }
func (MoneyOrYourLifeBlue) Name() string                               { return "Money or Your Life?" }
func (MoneyOrYourLifeBlue) Cost(*sim.TurnState) int                    { return 3 }
func (MoneyOrYourLifeBlue) Pitch() int                                 { return 3 }
func (MoneyOrYourLifeBlue) Attack() int                                { return 4 }
func (MoneyOrYourLifeBlue) Defense() int                               { return 2 }
func (MoneyOrYourLifeBlue) Types() card.TypeSet                        { return moneyOrYourLifeTypes }
func (MoneyOrYourLifeBlue) GoAgain() bool                              { return false }
func (MoneyOrYourLifeBlue) Play(s *sim.TurnState, self *sim.CardState) { moneyOrYourLifePlay(s, self) }
