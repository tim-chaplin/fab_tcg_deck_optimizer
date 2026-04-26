// Money or Your Life? — Generic Action - Attack. Cost 3. Printed power: Red 6, Yellow 5, Blue 4.
// Printed pitch variants: Red 1, Yellow 2, Blue 3. Defense 2.
//
// Text: "When this hits a hero, deal 2 damage to them unless they give you a Gold token they
// control. If you are a Thief, repeat this process once."

package generic

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var moneyOrYourLifeTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeAttack)

type MoneyOrYourLifeRed struct{}

func (MoneyOrYourLifeRed) ID() card.ID              { return card.MoneyOrYourLifeRed }
func (MoneyOrYourLifeRed) Name() string             { return "Money or Your Life?" }
func (MoneyOrYourLifeRed) Cost(*card.TurnState) int { return 3 }
func (MoneyOrYourLifeRed) Pitch() int               { return 1 }
func (MoneyOrYourLifeRed) Attack() int              { return 6 }
func (MoneyOrYourLifeRed) Defense() int             { return 2 }
func (MoneyOrYourLifeRed) Types() card.TypeSet      { return moneyOrYourLifeTypes }
func (MoneyOrYourLifeRed) GoAgain() bool            { return false }

// not implemented: gold tokens
func (MoneyOrYourLifeRed) NotImplemented() {}
func (c MoneyOrYourLifeRed) Play(s *card.TurnState, self *card.CardState) {
	s.ApplyAndLogEffectiveAttackPlus(self, moneyOrYourLifeDamage(c.Attack(), self)-self.Card.Attack())
}

type MoneyOrYourLifeYellow struct{}

func (MoneyOrYourLifeYellow) ID() card.ID              { return card.MoneyOrYourLifeYellow }
func (MoneyOrYourLifeYellow) Name() string             { return "Money or Your Life?" }
func (MoneyOrYourLifeYellow) Cost(*card.TurnState) int { return 3 }
func (MoneyOrYourLifeYellow) Pitch() int               { return 2 }
func (MoneyOrYourLifeYellow) Attack() int              { return 5 }
func (MoneyOrYourLifeYellow) Defense() int             { return 2 }
func (MoneyOrYourLifeYellow) Types() card.TypeSet      { return moneyOrYourLifeTypes }
func (MoneyOrYourLifeYellow) GoAgain() bool            { return false }

// not implemented: gold tokens
func (MoneyOrYourLifeYellow) NotImplemented() {}
func (c MoneyOrYourLifeYellow) Play(s *card.TurnState, self *card.CardState) {
	s.ApplyAndLogEffectiveAttackPlus(self, moneyOrYourLifeDamage(c.Attack(), self)-self.Card.Attack())
}

type MoneyOrYourLifeBlue struct{}

func (MoneyOrYourLifeBlue) ID() card.ID              { return card.MoneyOrYourLifeBlue }
func (MoneyOrYourLifeBlue) Name() string             { return "Money or Your Life?" }
func (MoneyOrYourLifeBlue) Cost(*card.TurnState) int { return 3 }
func (MoneyOrYourLifeBlue) Pitch() int               { return 3 }
func (MoneyOrYourLifeBlue) Attack() int              { return 4 }
func (MoneyOrYourLifeBlue) Defense() int             { return 2 }
func (MoneyOrYourLifeBlue) Types() card.TypeSet      { return moneyOrYourLifeTypes }
func (MoneyOrYourLifeBlue) GoAgain() bool            { return false }

// not implemented: gold tokens
func (MoneyOrYourLifeBlue) NotImplemented() {}
func (c MoneyOrYourLifeBlue) Play(s *card.TurnState, self *card.CardState) {
	s.ApplyAndLogEffectiveAttackPlus(self, moneyOrYourLifeDamage(c.Attack(), self)-self.Card.Attack())
}

// moneyOrYourLifeDamage is a breadcrumb for the on-hit "deal 2 damage unless they give Gold"
// rider — Gold tokens aren't tracked, Thief-repeat isn't modelled (see TODO.md).
func moneyOrYourLifeDamage(attack int, self *card.CardState) int {
	if card.LikelyToHit(self) {
		// TODO: model on-hit Gold-trade rider.
	}
	return attack
}
