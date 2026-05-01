// Unmovable — Generic Defense Reaction. Cost 3.
// Printed pitch variants: Red 1, Yellow 2, Blue 3.
// Printed defense: Red 7, Yellow 6, Blue 5.
// Text: "If Unmovable is played from arsenal, it gains +1{d}."
//
// +1{d} when played from arsenal via sim.ArsenalDefenseBonus (docs/dev-standards.md).

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

type UnmovableRed struct{}

func (UnmovableRed) ID() ids.CardID          { return ids.UnmovableRed }
func (UnmovableRed) Name() string            { return "Unmovable" }
func (UnmovableRed) Cost(*sim.TurnState) int { return 3 }
func (UnmovableRed) Pitch() int              { return 1 }
func (UnmovableRed) Attack() int             { return 0 }
func (UnmovableRed) Defense() int            { return 7 }
func (UnmovableRed) Types() card.TypeSet     { return defenseReactionTypes }
func (UnmovableRed) GoAgain() bool           { return false }
func (UnmovableRed) Play(s *sim.TurnState, self *sim.CardState) {
	s.LogChain(self, s.ApplyDefenseValue(self.EffectiveDefense()))
}
func (UnmovableRed) ArsenalDefenseBonus() int { return 1 }

type UnmovableYellow struct{}

func (UnmovableYellow) ID() ids.CardID          { return ids.UnmovableYellow }
func (UnmovableYellow) Name() string            { return "Unmovable" }
func (UnmovableYellow) Cost(*sim.TurnState) int { return 3 }
func (UnmovableYellow) Pitch() int              { return 2 }
func (UnmovableYellow) Attack() int             { return 0 }
func (UnmovableYellow) Defense() int            { return 6 }
func (UnmovableYellow) Types() card.TypeSet     { return defenseReactionTypes }
func (UnmovableYellow) GoAgain() bool           { return false }
func (UnmovableYellow) Play(s *sim.TurnState, self *sim.CardState) {
	s.LogChain(self, s.ApplyDefenseValue(self.EffectiveDefense()))
}
func (UnmovableYellow) ArsenalDefenseBonus() int { return 1 }

type UnmovableBlue struct{}

func (UnmovableBlue) ID() ids.CardID          { return ids.UnmovableBlue }
func (UnmovableBlue) Name() string            { return "Unmovable" }
func (UnmovableBlue) Cost(*sim.TurnState) int { return 3 }
func (UnmovableBlue) Pitch() int              { return 3 }
func (UnmovableBlue) Attack() int             { return 0 }
func (UnmovableBlue) Defense() int            { return 5 }
func (UnmovableBlue) Types() card.TypeSet     { return defenseReactionTypes }
func (UnmovableBlue) GoAgain() bool           { return false }
func (UnmovableBlue) Play(s *sim.TurnState, self *sim.CardState) {
	s.LogChain(self, s.ApplyDefenseValue(self.EffectiveDefense()))
}
func (UnmovableBlue) ArsenalDefenseBonus() int { return 1 }
