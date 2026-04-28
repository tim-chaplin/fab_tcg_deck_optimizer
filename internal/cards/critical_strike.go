// Critical Strike — Generic Action - Attack. Cost 1. Printed power: Red 5, Yellow 4, Blue 3.
// Printed pitch variants: Red 1, Yellow 2, Blue 3. Defense 3.

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

var criticalStrikeTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeAttack)

type CriticalStrikeRed struct{}

func (CriticalStrikeRed) ID() ids.CardID          { return ids.CriticalStrikeRed }
func (CriticalStrikeRed) Name() string            { return "Critical Strike" }
func (CriticalStrikeRed) Cost(*sim.TurnState) int { return 1 }
func (CriticalStrikeRed) Pitch() int              { return 1 }
func (CriticalStrikeRed) Attack() int             { return 5 }
func (CriticalStrikeRed) Defense() int            { return 3 }
func (CriticalStrikeRed) Types() card.TypeSet     { return criticalStrikeTypes }
func (CriticalStrikeRed) GoAgain() bool           { return false }
func (c CriticalStrikeRed) Play(s *sim.TurnState, self *sim.CardState) {
	s.ApplyAndLogEffectiveAttack(self)
}

type CriticalStrikeYellow struct{}

func (CriticalStrikeYellow) ID() ids.CardID          { return ids.CriticalStrikeYellow }
func (CriticalStrikeYellow) Name() string            { return "Critical Strike" }
func (CriticalStrikeYellow) Cost(*sim.TurnState) int { return 1 }
func (CriticalStrikeYellow) Pitch() int              { return 2 }
func (CriticalStrikeYellow) Attack() int             { return 4 }
func (CriticalStrikeYellow) Defense() int            { return 3 }
func (CriticalStrikeYellow) Types() card.TypeSet     { return criticalStrikeTypes }
func (CriticalStrikeYellow) GoAgain() bool           { return false }
func (c CriticalStrikeYellow) Play(s *sim.TurnState, self *sim.CardState) {
	s.ApplyAndLogEffectiveAttack(self)
}

type CriticalStrikeBlue struct{}

func (CriticalStrikeBlue) ID() ids.CardID          { return ids.CriticalStrikeBlue }
func (CriticalStrikeBlue) Name() string            { return "Critical Strike" }
func (CriticalStrikeBlue) Cost(*sim.TurnState) int { return 1 }
func (CriticalStrikeBlue) Pitch() int              { return 3 }
func (CriticalStrikeBlue) Attack() int             { return 3 }
func (CriticalStrikeBlue) Defense() int            { return 3 }
func (CriticalStrikeBlue) Types() card.TypeSet     { return criticalStrikeTypes }
func (CriticalStrikeBlue) GoAgain() bool           { return false }
func (c CriticalStrikeBlue) Play(s *sim.TurnState, self *sim.CardState) {
	s.ApplyAndLogEffectiveAttack(self)
}
