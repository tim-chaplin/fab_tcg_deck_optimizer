// Rune Flash — Runeblade Action - Attack. Printed cost 3, Defense 3. Go again.
// Printed pitch variants: Red 1, Yellow 2, Blue 3.
// Printed power: Red 4, Yellow 3, Blue 2.
// Text: "Rune Flash costs {r} less to play for each Runechant you control."
//
// Variable cost: Cost reads s.Runechants() to return max(0, printed - Runechants).
// Standard sim.VariableCost wiring (docs/dev-standards.md).

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

const runeFlashPrintedCost = 3

func runeFlashCost(s *sim.TurnState) int {
	eff := runeFlashPrintedCost - s.Runechants()
	if eff < 0 {
		return 0
	}
	return eff
}

func (RuneFlashRed) Cost(s *sim.TurnState) int { return runeFlashCost(s) }
func (RuneFlashRed) MinCost() int              { return 0 }
func (RuneFlashRed) MaxCost() int              { return runeFlashPrintedCost }

func (RuneFlashRed) Play(s *sim.TurnState, l sim.Logger, self *sim.CardState) {
	n := self.DealEffectiveAttack(s)
	self.Log(l, n)
}

func (RuneFlashYellow) Cost(s *sim.TurnState) int { return runeFlashCost(s) }
func (RuneFlashYellow) MinCost() int              { return 0 }
func (RuneFlashYellow) MaxCost() int              { return runeFlashPrintedCost }

func (RuneFlashYellow) Play(s *sim.TurnState, l sim.Logger, self *sim.CardState) {
	n := self.DealEffectiveAttack(s)
	self.Log(l, n)
}

func (RuneFlashBlue) Cost(s *sim.TurnState) int { return runeFlashCost(s) }
func (RuneFlashBlue) MinCost() int              { return 0 }
func (RuneFlashBlue) MaxCost() int              { return runeFlashPrintedCost }

func (RuneFlashBlue) Play(s *sim.TurnState, l sim.Logger, self *sim.CardState) {
	n := self.DealEffectiveAttack(s)
	self.Log(l, n)
}
