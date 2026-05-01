// Singeing Steelblade — Runeblade Action - Attack. Cost 1, Defense 3, Arcane 1.
// Printed power: Red 4, Yellow 3, Blue 2.
// Text: "When you attack with Singeing Steelblade, deal 1 arcane damage to target hero."
//
// The printed 1 arcane is added to combat damage (both hit the same target). Play also sets
// ArcaneDamageDealt so same-turn triggers keyed on that flag fire.

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

var singeingSteelbladeTypes = card.NewTypeSet(card.TypeRuneblade, card.TypeAction, card.TypeAttack)

type SingeingSteelbladeRed struct{}

func (SingeingSteelbladeRed) ID() ids.CardID          { return ids.SingeingSteelbladeRed }
func (SingeingSteelbladeRed) Name() string            { return "Singeing Steelblade" }
func (SingeingSteelbladeRed) Cost(*sim.TurnState) int { return 1 }
func (SingeingSteelbladeRed) Pitch() int              { return 1 }
func (SingeingSteelbladeRed) Attack() int             { return 4 }
func (SingeingSteelbladeRed) Defense() int            { return 3 }
func (SingeingSteelbladeRed) Types() card.TypeSet     { return singeingSteelbladeTypes }
func (SingeingSteelbladeRed) GoAgain() bool           { return false }
func (SingeingSteelbladeRed) Play(s *sim.TurnState, self *sim.CardState) {
	n := self.DealEffectiveAttack(s)
	s.Log(self, n)
	s.AddValue(s.DealArcaneDamage(1))
	s.LogRider(self, 1, "Dealt 1 arcane damage")
}

type SingeingSteelbladeYellow struct{}

func (SingeingSteelbladeYellow) ID() ids.CardID          { return ids.SingeingSteelbladeYellow }
func (SingeingSteelbladeYellow) Name() string            { return "Singeing Steelblade" }
func (SingeingSteelbladeYellow) Cost(*sim.TurnState) int { return 1 }
func (SingeingSteelbladeYellow) Pitch() int              { return 2 }
func (SingeingSteelbladeYellow) Attack() int             { return 3 }
func (SingeingSteelbladeYellow) Defense() int            { return 3 }
func (SingeingSteelbladeYellow) Types() card.TypeSet     { return singeingSteelbladeTypes }
func (SingeingSteelbladeYellow) GoAgain() bool           { return false }
func (SingeingSteelbladeYellow) Play(s *sim.TurnState, self *sim.CardState) {
	n := self.DealEffectiveAttack(s)
	s.Log(self, n)
	s.AddValue(s.DealArcaneDamage(1))
	s.LogRider(self, 1, "Dealt 1 arcane damage")
}

type SingeingSteelbladeBlue struct{}

func (SingeingSteelbladeBlue) ID() ids.CardID          { return ids.SingeingSteelbladeBlue }
func (SingeingSteelbladeBlue) Name() string            { return "Singeing Steelblade" }
func (SingeingSteelbladeBlue) Cost(*sim.TurnState) int { return 1 }
func (SingeingSteelbladeBlue) Pitch() int              { return 3 }
func (SingeingSteelbladeBlue) Attack() int             { return 2 }
func (SingeingSteelbladeBlue) Defense() int            { return 3 }
func (SingeingSteelbladeBlue) Types() card.TypeSet     { return singeingSteelbladeTypes }
func (SingeingSteelbladeBlue) GoAgain() bool           { return false }
func (SingeingSteelbladeBlue) Play(s *sim.TurnState, self *sim.CardState) {
	n := self.DealEffectiveAttack(s)
	s.Log(self, n)
	s.AddValue(s.DealArcaneDamage(1))
	s.LogRider(self, 1, "Dealt 1 arcane damage")
}
