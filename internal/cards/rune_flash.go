// Rune Flash — Runeblade Action - Attack. Printed cost 3, Defense 3. Go again.
// Printed pitch variants: Red 1, Yellow 2, Blue 3.
// Printed power: Red 4, Yellow 3, Blue 2.
// Text: "Rune Flash costs {r} less to play for each Runechant you control."
//
// Variable cost: Cost reads s.Runechants to return max(0, printed - Runechants).
// Standard sim.VariableCost wiring (docs/dev-standards.md).

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

var runeFlashTypes = card.NewTypeSet(card.TypeRuneblade, card.TypeAction, card.TypeAttack)

const runeFlashPrintedCost = 3

func runeFlashCost(s *sim.TurnState) int {
	eff := runeFlashPrintedCost - s.Runechants
	if eff < 0 {
		return 0
	}
	return eff
}

type RuneFlashRed struct{}

func (RuneFlashRed) ID() ids.CardID                             { return ids.RuneFlashRed }
func (RuneFlashRed) Name() string                               { return "Rune Flash" }
func (RuneFlashRed) Cost(s *sim.TurnState) int                  { return runeFlashCost(s) }
func (RuneFlashRed) MinCost() int                               { return 0 }
func (RuneFlashRed) MaxCost() int                               { return runeFlashPrintedCost }
func (RuneFlashRed) Pitch() int                                 { return 1 }
func (RuneFlashRed) Attack() int                                { return 4 }
func (RuneFlashRed) Defense() int                               { return 3 }
func (RuneFlashRed) Types() card.TypeSet                        { return runeFlashTypes }
func (RuneFlashRed) GoAgain() bool                              { return true }
func (RuneFlashRed) Play(s *sim.TurnState, self *sim.CardState) { s.ApplyAndLogEffectiveAttack(self) }

type RuneFlashYellow struct{}

func (RuneFlashYellow) ID() ids.CardID            { return ids.RuneFlashYellow }
func (RuneFlashYellow) Name() string              { return "Rune Flash" }
func (RuneFlashYellow) Cost(s *sim.TurnState) int { return runeFlashCost(s) }
func (RuneFlashYellow) MinCost() int              { return 0 }
func (RuneFlashYellow) MaxCost() int              { return runeFlashPrintedCost }
func (RuneFlashYellow) Pitch() int                { return 2 }
func (RuneFlashYellow) Attack() int               { return 3 }
func (RuneFlashYellow) Defense() int              { return 3 }
func (RuneFlashYellow) Types() card.TypeSet       { return runeFlashTypes }
func (RuneFlashYellow) GoAgain() bool             { return true }
func (RuneFlashYellow) Play(s *sim.TurnState, self *sim.CardState) {
	s.ApplyAndLogEffectiveAttack(self)
}

type RuneFlashBlue struct{}

func (RuneFlashBlue) ID() ids.CardID            { return ids.RuneFlashBlue }
func (RuneFlashBlue) Name() string              { return "Rune Flash" }
func (RuneFlashBlue) Cost(s *sim.TurnState) int { return runeFlashCost(s) }
func (RuneFlashBlue) MinCost() int              { return 0 }
func (RuneFlashBlue) MaxCost() int              { return runeFlashPrintedCost }
func (RuneFlashBlue) Pitch() int                { return 3 }
func (RuneFlashBlue) Attack() int               { return 2 }
func (RuneFlashBlue) Defense() int              { return 3 }
func (RuneFlashBlue) Types() card.TypeSet       { return runeFlashTypes }
func (RuneFlashBlue) GoAgain() bool             { return true }
func (RuneFlashBlue) Play(s *sim.TurnState, self *sim.CardState) {
	s.ApplyAndLogEffectiveAttack(self)
}
