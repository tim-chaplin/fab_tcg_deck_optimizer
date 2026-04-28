// Money Where Ya Mouth Is — Generic Action. Cost 1. Printed pitch variants: Red 1, Yellow 2, Blue
// 3. Defense 2.
//
// Text: "Your next attack this turn gets +N{p} and "When this attacks a hero, you may **wager** a
// Gold token with them."" (Red N=3, Yellow N=2, Blue N=1.)
//
// Scans TurnState.CardsRemaining for the first matching attack action card and credits the
// +N{p} bonus assuming it will be played; if none is scheduled after this card, the bonus
// fizzles.

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

var moneyWhereYaMouthIsTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction)

type MoneyWhereYaMouthIsRed struct{}

func (MoneyWhereYaMouthIsRed) ID() ids.CardID          { return ids.MoneyWhereYaMouthIsRed }
func (MoneyWhereYaMouthIsRed) Name() string            { return "Money Where Ya Mouth Is" }
func (MoneyWhereYaMouthIsRed) Cost(*sim.TurnState) int { return 1 }
func (MoneyWhereYaMouthIsRed) Pitch() int              { return 1 }
func (MoneyWhereYaMouthIsRed) Attack() int             { return 0 }
func (MoneyWhereYaMouthIsRed) Defense() int            { return 2 }
func (MoneyWhereYaMouthIsRed) Types() card.TypeSet     { return moneyWhereYaMouthIsTypes }
func (MoneyWhereYaMouthIsRed) GoAgain() bool           { return true }

// not implemented: gold tokens
func (MoneyWhereYaMouthIsRed) NotImplemented() {}
func (MoneyWhereYaMouthIsRed) Play(s *sim.TurnState, self *sim.CardState) {
	grantNextAttackActionBonus(s, 3)
	s.ApplyAndLogEffectiveAttack(self)
}

type MoneyWhereYaMouthIsYellow struct{}

func (MoneyWhereYaMouthIsYellow) ID() ids.CardID          { return ids.MoneyWhereYaMouthIsYellow }
func (MoneyWhereYaMouthIsYellow) Name() string            { return "Money Where Ya Mouth Is" }
func (MoneyWhereYaMouthIsYellow) Cost(*sim.TurnState) int { return 1 }
func (MoneyWhereYaMouthIsYellow) Pitch() int              { return 2 }
func (MoneyWhereYaMouthIsYellow) Attack() int             { return 0 }
func (MoneyWhereYaMouthIsYellow) Defense() int            { return 2 }
func (MoneyWhereYaMouthIsYellow) Types() card.TypeSet     { return moneyWhereYaMouthIsTypes }
func (MoneyWhereYaMouthIsYellow) GoAgain() bool           { return true }

// not implemented: gold tokens
func (MoneyWhereYaMouthIsYellow) NotImplemented() {}
func (MoneyWhereYaMouthIsYellow) Play(s *sim.TurnState, self *sim.CardState) {
	grantNextAttackActionBonus(s, 2)
	s.ApplyAndLogEffectiveAttack(self)
}

type MoneyWhereYaMouthIsBlue struct{}

func (MoneyWhereYaMouthIsBlue) ID() ids.CardID          { return ids.MoneyWhereYaMouthIsBlue }
func (MoneyWhereYaMouthIsBlue) Name() string            { return "Money Where Ya Mouth Is" }
func (MoneyWhereYaMouthIsBlue) Cost(*sim.TurnState) int { return 1 }
func (MoneyWhereYaMouthIsBlue) Pitch() int              { return 3 }
func (MoneyWhereYaMouthIsBlue) Attack() int             { return 0 }
func (MoneyWhereYaMouthIsBlue) Defense() int            { return 2 }
func (MoneyWhereYaMouthIsBlue) Types() card.TypeSet     { return moneyWhereYaMouthIsTypes }
func (MoneyWhereYaMouthIsBlue) GoAgain() bool           { return true }

// not implemented: gold tokens
func (MoneyWhereYaMouthIsBlue) NotImplemented() {}
func (MoneyWhereYaMouthIsBlue) Play(s *sim.TurnState, self *sim.CardState) {
	grantNextAttackActionBonus(s, 1)
	s.ApplyAndLogEffectiveAttack(self)
}
