// Rune Flash — Runeblade Action - Attack. Printed cost 3, Defense 3. Go again.
// Printed pitch variants: Red 1, Yellow 2, Blue 3.
// Printed power: Red 4, Yellow 3, Blue 2.
// Text: "Rune Flash costs {r} less to play for each Runechant you control."
//
// Variable cost: Cost reads g.RunechantCount() to return max(0, printed - Runechants).
// Standard card.VariableCost wiring (docs/dev-standards.md).

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
)

const runeFlashPrintedCost = 3

func runeFlashCost(g card.GameEngine) int {
	eff := runeFlashPrintedCost - g.RunechantCount()
	if eff < 0 {
		return 0
	}
	return eff
}

func (RuneFlashRed) Cost(g card.GameEngine) int { return runeFlashCost(g) }
func (RuneFlashRed) MinCost() int               { return 0 }
func (RuneFlashRed) MaxCost() int               { return runeFlashPrintedCost }

func (RuneFlashRed) Play(g card.GameEngine, l card.Logger, self *card.CardState) {
}

func (RuneFlashYellow) Cost(g card.GameEngine) int { return runeFlashCost(g) }
func (RuneFlashYellow) MinCost() int               { return 0 }
func (RuneFlashYellow) MaxCost() int               { return runeFlashPrintedCost }

func (RuneFlashYellow) Play(g card.GameEngine, l card.Logger, self *card.CardState) {
}

func (RuneFlashBlue) Cost(g card.GameEngine) int { return runeFlashCost(g) }
func (RuneFlashBlue) MinCost() int               { return 0 }
func (RuneFlashBlue) MaxCost() int               { return runeFlashPrintedCost }

func (RuneFlashBlue) Play(g card.GameEngine, l card.Logger, self *card.CardState) {
}
