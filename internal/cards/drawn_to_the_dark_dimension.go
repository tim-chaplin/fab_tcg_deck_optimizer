// Drawn to the Dark Dimension — Runeblade Action - Attack. Printed cost 2, Defense 3.
// Printed pitch variants: Red 1, Yellow 2, Blue 3.
// Printed power: Red 3, Yellow 2, Blue 1.
// Text: "Drawn to the Dark Dimension costs {r} less to play for each Runechant you control.
// Draw a card."
//
// Cost returns max(0, printed - RunechantCount).

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
)

const drawnToTheDarkDimensionPrintedCost = 2

func drawnToTheDarkDimensionCost(g card.GameEngine) int {
	eff := drawnToTheDarkDimensionPrintedCost - g.RunechantCount()
	if eff < 0 {
		return 0
	}
	return eff
}

func (DrawnToTheDarkDimensionRed) Cost(g card.GameEngine) int { return drawnToTheDarkDimensionCost(g) }
func (DrawnToTheDarkDimensionRed) MinCost() int               { return 0 }
func (DrawnToTheDarkDimensionRed) MaxCost() int               { return drawnToTheDarkDimensionPrintedCost }

func (c DrawnToTheDarkDimensionRed) Play(g card.GameEngine, l card.Logger, self *card.CardState) {
	g.DrawOne()
}

func (DrawnToTheDarkDimensionYellow) Cost(g card.GameEngine) int {
	return drawnToTheDarkDimensionCost(g)
}
func (DrawnToTheDarkDimensionYellow) MinCost() int { return 0 }
func (DrawnToTheDarkDimensionYellow) MaxCost() int { return drawnToTheDarkDimensionPrintedCost }

func (c DrawnToTheDarkDimensionYellow) Play(g card.GameEngine, l card.Logger, self *card.CardState) {
	g.DrawOne()
}

func (DrawnToTheDarkDimensionBlue) Cost(g card.GameEngine) int { return drawnToTheDarkDimensionCost(g) }
func (DrawnToTheDarkDimensionBlue) MinCost() int               { return 0 }
func (DrawnToTheDarkDimensionBlue) MaxCost() int               { return drawnToTheDarkDimensionPrintedCost }

func (c DrawnToTheDarkDimensionBlue) Play(g card.GameEngine, l card.Logger, self *card.CardState) {
	g.DrawOne()
}
