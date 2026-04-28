// Arcanic Crackle — Runeblade Action - Attack. Cost 0, Defense 3, Arcane 1.
// Printed power: Red 3, Yellow 2, Blue 1.
// Text: "Deal 1 arcane damage to target hero."
//
// The printed 1 arcane is added to combat damage (both hit the same target). Play also sets
// ArcaneDamageDealt so same-turn triggers reading "if you've dealt arcane damage this turn"
// fire.

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

var arcanicCrackleTypes = card.NewTypeSet(card.TypeRuneblade, card.TypeAction, card.TypeAttack)

type ArcanicCrackleRed struct{}

func (ArcanicCrackleRed) ID() ids.CardID          { return ids.ArcanicCrackleRed }
func (ArcanicCrackleRed) Name() string            { return "Arcanic Crackle" }
func (ArcanicCrackleRed) Cost(*sim.TurnState) int { return 0 }
func (ArcanicCrackleRed) Pitch() int              { return 1 }
func (ArcanicCrackleRed) Attack() int             { return 3 }
func (ArcanicCrackleRed) Defense() int            { return 3 }
func (ArcanicCrackleRed) Types() card.TypeSet     { return arcanicCrackleTypes }
func (ArcanicCrackleRed) GoAgain() bool           { return false }
func (ArcanicCrackleRed) Play(s *sim.TurnState, self *sim.CardState) {
	s.ApplyAndLogEffectiveAttack(self)
	s.DealAndLogArcaneDamage(self, 1)
}

type ArcanicCrackleYellow struct{}

func (ArcanicCrackleYellow) ID() ids.CardID          { return ids.ArcanicCrackleYellow }
func (ArcanicCrackleYellow) Name() string            { return "Arcanic Crackle" }
func (ArcanicCrackleYellow) Cost(*sim.TurnState) int { return 0 }
func (ArcanicCrackleYellow) Pitch() int              { return 2 }
func (ArcanicCrackleYellow) Attack() int             { return 2 }
func (ArcanicCrackleYellow) Defense() int            { return 3 }
func (ArcanicCrackleYellow) Types() card.TypeSet     { return arcanicCrackleTypes }
func (ArcanicCrackleYellow) GoAgain() bool           { return false }
func (ArcanicCrackleYellow) Play(s *sim.TurnState, self *sim.CardState) {
	s.ApplyAndLogEffectiveAttack(self)
	s.DealAndLogArcaneDamage(self, 1)
}

type ArcanicCrackleBlue struct{}

func (ArcanicCrackleBlue) ID() ids.CardID          { return ids.ArcanicCrackleBlue }
func (ArcanicCrackleBlue) Name() string            { return "Arcanic Crackle" }
func (ArcanicCrackleBlue) Cost(*sim.TurnState) int { return 0 }
func (ArcanicCrackleBlue) Pitch() int              { return 3 }
func (ArcanicCrackleBlue) Attack() int             { return 1 }
func (ArcanicCrackleBlue) Defense() int            { return 3 }
func (ArcanicCrackleBlue) Types() card.TypeSet     { return arcanicCrackleTypes }
func (ArcanicCrackleBlue) GoAgain() bool           { return false }
func (ArcanicCrackleBlue) Play(s *sim.TurnState, self *sim.CardState) {
	s.ApplyAndLogEffectiveAttack(self)
	s.DealAndLogArcaneDamage(self, 1)
}
