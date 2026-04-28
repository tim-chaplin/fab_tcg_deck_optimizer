// Springboard Somersault — Generic Defense Reaction. Cost 0, Pitch 2, Defense 2. Only printed in
// Yellow.
// Text: "If Springboard Somersault is played from arsenal, it gains +2{d}."
//
// +2{d} when played from arsenal via sim.ArsenalDefenseBonus (docs/dev-standards.md).

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

type SpringboardSomersaultYellow struct{}

func (SpringboardSomersaultYellow) ID() ids.CardID          { return ids.SpringboardSomersaultYellow }
func (SpringboardSomersaultYellow) Name() string            { return "Springboard Somersault" }
func (SpringboardSomersaultYellow) Cost(*sim.TurnState) int { return 0 }
func (SpringboardSomersaultYellow) Pitch() int              { return 2 }
func (SpringboardSomersaultYellow) Attack() int             { return 0 }
func (SpringboardSomersaultYellow) Defense() int            { return 2 }
func (SpringboardSomersaultYellow) Types() card.TypeSet     { return defenseReactionTypes }
func (SpringboardSomersaultYellow) GoAgain() bool           { return false }
func (SpringboardSomersaultYellow) Play(s *sim.TurnState, self *sim.CardState) {
	s.ApplyAndLogEffectiveDefense(self)
}
func (SpringboardSomersaultYellow) ArsenalDefenseBonus() int { return 2 }
