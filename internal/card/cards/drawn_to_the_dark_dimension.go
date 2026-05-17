// Drawn to the Dark Dimension — Runeblade Action - Attack. Printed cost 2, Defense 3.
// Printed pitch variants: Red 1, Yellow 2, Blue 3.
// Printed power: Red 3, Yellow 2, Blue 1.
// Text: "Drawn to the Dark Dimension costs {r} less to play for each Runechant you control.
// Draw a card."
//
// Cost returns max(0, printed - RunechantCount).

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

const drawnToTheDarkDimensionPrintedCost = 2

func drawnToTheDarkDimensionCost(ge card.GameEngine) int {
	eff := drawnToTheDarkDimensionPrintedCost - ge.RunechantCount()
	if eff < 0 {
		return 0
	}
	return eff
}

func (DrawnToTheDarkDimensionRed) Cost(ge card.GameEngine) int {
	return drawnToTheDarkDimensionCost(ge)
}
func (DrawnToTheDarkDimensionRed) MinCost() int { return 0 }
func (DrawnToTheDarkDimensionRed) MaxCost() int { return drawnToTheDarkDimensionPrintedCost }

func (c DrawnToTheDarkDimensionRed) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	ge.DrawOne()
}

func (DrawnToTheDarkDimensionYellow) Cost(ge card.GameEngine) int {
	return drawnToTheDarkDimensionCost(ge)
}
func (DrawnToTheDarkDimensionYellow) MinCost() int { return 0 }
func (DrawnToTheDarkDimensionYellow) MaxCost() int { return drawnToTheDarkDimensionPrintedCost }

func (c DrawnToTheDarkDimensionYellow) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	ge.DrawOne()
}

func (DrawnToTheDarkDimensionBlue) Cost(ge card.GameEngine) int {
	return drawnToTheDarkDimensionCost(ge)
}
func (DrawnToTheDarkDimensionBlue) MinCost() int { return 0 }
func (DrawnToTheDarkDimensionBlue) MaxCost() int { return drawnToTheDarkDimensionPrintedCost }

func (c DrawnToTheDarkDimensionBlue) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	ge.DrawOne()
}
