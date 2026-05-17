// Rune Flash — Runeblade Action - Attack. Printed cost 3, Defense 3. Go again.
// Printed pitch variants: Red 1, Yellow 2, Blue 3.
// Printed power: Red 4, Yellow 3, Blue 2.
// Text: "Rune Flash costs {r} less to play for each Runechant you control."
//
// Cost returns max(0, printed - RunechantCount).

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

const runeFlashPrintedCost = 3

func runeFlashCost(ge card.GameEngine) int {
	eff := runeFlashPrintedCost - ge.RunechantCount()
	if eff < 0 {
		return 0
	}
	return eff
}

func (RuneFlashRed) Cost(ge card.GameEngine) int { return runeFlashCost(ge) }
func (RuneFlashRed) MinCost() int                { return 0 }
func (RuneFlashRed) MaxCost() int                { return runeFlashPrintedCost }

func (RuneFlashRed) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
}

func (RuneFlashYellow) Cost(ge card.GameEngine) int { return runeFlashCost(ge) }
func (RuneFlashYellow) MinCost() int                { return 0 }
func (RuneFlashYellow) MaxCost() int                { return runeFlashPrintedCost }

func (RuneFlashYellow) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
}

func (RuneFlashBlue) Cost(ge card.GameEngine) int { return runeFlashCost(ge) }
func (RuneFlashBlue) MinCost() int                { return 0 }
func (RuneFlashBlue) MaxCost() int                { return runeFlashPrintedCost }

func (RuneFlashBlue) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
}
