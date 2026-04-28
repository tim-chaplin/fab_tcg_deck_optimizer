// Evasive Leap — Generic Defense Reaction. Cost 0.
// Printed pitch variants: Red 1, Yellow 2, Blue 3.
// Printed defense: Red 3, Yellow 2, Blue 1.

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

type EvasiveLeapRed struct{}

func (EvasiveLeapRed) ID() ids.CardID          { return ids.EvasiveLeapRed }
func (EvasiveLeapRed) Name() string            { return "Evasive Leap" }
func (EvasiveLeapRed) Cost(*sim.TurnState) int { return 0 }
func (EvasiveLeapRed) Pitch() int              { return 1 }
func (EvasiveLeapRed) Attack() int             { return 0 }
func (EvasiveLeapRed) Defense() int            { return 3 }
func (EvasiveLeapRed) Types() card.TypeSet     { return defenseReactionTypes }
func (EvasiveLeapRed) GoAgain() bool           { return false }
func (EvasiveLeapRed) Play(s *sim.TurnState, self *sim.CardState) {
	s.ApplyAndLogEffectiveDefense(self)
}

type EvasiveLeapYellow struct{}

func (EvasiveLeapYellow) ID() ids.CardID          { return ids.EvasiveLeapYellow }
func (EvasiveLeapYellow) Name() string            { return "Evasive Leap" }
func (EvasiveLeapYellow) Cost(*sim.TurnState) int { return 0 }
func (EvasiveLeapYellow) Pitch() int              { return 2 }
func (EvasiveLeapYellow) Attack() int             { return 0 }
func (EvasiveLeapYellow) Defense() int            { return 2 }
func (EvasiveLeapYellow) Types() card.TypeSet     { return defenseReactionTypes }
func (EvasiveLeapYellow) GoAgain() bool           { return false }
func (EvasiveLeapYellow) Play(s *sim.TurnState, self *sim.CardState) {
	s.ApplyAndLogEffectiveDefense(self)
}

type EvasiveLeapBlue struct{}

func (EvasiveLeapBlue) ID() ids.CardID          { return ids.EvasiveLeapBlue }
func (EvasiveLeapBlue) Name() string            { return "Evasive Leap" }
func (EvasiveLeapBlue) Cost(*sim.TurnState) int { return 0 }
func (EvasiveLeapBlue) Pitch() int              { return 3 }
func (EvasiveLeapBlue) Attack() int             { return 0 }
func (EvasiveLeapBlue) Defense() int            { return 1 }
func (EvasiveLeapBlue) Types() card.TypeSet     { return defenseReactionTypes }
func (EvasiveLeapBlue) GoAgain() bool           { return false }
func (EvasiveLeapBlue) Play(s *sim.TurnState, self *sim.CardState) {
	s.ApplyAndLogEffectiveDefense(self)
}
