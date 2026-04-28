// Wounded Bull — Generic Action - Attack. Cost 3. Printed power: Red 7, Yellow 6, Blue 5. Printed
// pitch variants: Red 1, Yellow 2, Blue 3. Defense 2.
//
// Text: "When you play this, if you have less {h} than an opposing hero, this gains +1{p}."

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

var woundedBullTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeAttack)

// woundedBullBonus returns the +1{p} power buff when the current hero opts into
// LowerHealthWanter, else 0.
func woundedBullBonus() int {
	if sim.HeroWantsLowerHealth() {
		return 1
	}
	return 0
}

type WoundedBullRed struct{}

func (WoundedBullRed) ID() ids.CardID          { return ids.WoundedBullRed }
func (WoundedBullRed) Name() string            { return "Wounded Bull" }
func (WoundedBullRed) Cost(*sim.TurnState) int { return 3 }
func (WoundedBullRed) Pitch() int              { return 1 }
func (WoundedBullRed) Attack() int             { return 7 }
func (WoundedBullRed) Defense() int            { return 2 }
func (WoundedBullRed) Types() card.TypeSet     { return woundedBullTypes }
func (WoundedBullRed) GoAgain() bool           { return false }
func (WoundedBullRed) Play(s *sim.TurnState, self *sim.CardState) {
	self.BonusAttack += woundedBullBonus()
	s.ApplyAndLogEffectiveAttack(self)
}

type WoundedBullYellow struct{}

func (WoundedBullYellow) ID() ids.CardID          { return ids.WoundedBullYellow }
func (WoundedBullYellow) Name() string            { return "Wounded Bull" }
func (WoundedBullYellow) Cost(*sim.TurnState) int { return 3 }
func (WoundedBullYellow) Pitch() int              { return 2 }
func (WoundedBullYellow) Attack() int             { return 6 }
func (WoundedBullYellow) Defense() int            { return 2 }
func (WoundedBullYellow) Types() card.TypeSet     { return woundedBullTypes }
func (WoundedBullYellow) GoAgain() bool           { return false }
func (WoundedBullYellow) Play(s *sim.TurnState, self *sim.CardState) {
	self.BonusAttack += woundedBullBonus()
	s.ApplyAndLogEffectiveAttack(self)
}

type WoundedBullBlue struct{}

func (WoundedBullBlue) ID() ids.CardID          { return ids.WoundedBullBlue }
func (WoundedBullBlue) Name() string            { return "Wounded Bull" }
func (WoundedBullBlue) Cost(*sim.TurnState) int { return 3 }
func (WoundedBullBlue) Pitch() int              { return 3 }
func (WoundedBullBlue) Attack() int             { return 5 }
func (WoundedBullBlue) Defense() int            { return 2 }
func (WoundedBullBlue) Types() card.TypeSet     { return woundedBullTypes }
func (WoundedBullBlue) GoAgain() bool           { return false }
func (WoundedBullBlue) Play(s *sim.TurnState, self *sim.CardState) {
	self.BonusAttack += woundedBullBonus()
	s.ApplyAndLogEffectiveAttack(self)
}
