// Drawn to the Dark Dimension — Runeblade Action - Attack. Printed cost 2, Defense 3.
// Printed pitch variants: Red 1, Yellow 2, Blue 3.
// Printed power: Red 3, Yellow 2, Blue 1.
// Text: "Drawn to the Dark Dimension costs {r} less to play for each Runechant you control.
// Draw a card."
//
// Cost reads s.Runechants() to return max(0, printed - Runechants) at play time; implements
// card.VariableCost with bounds [0, printed].
//
// The "Draw a card" rider fires unconditionally on play.

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
)

const drawnToTheDarkDimensionPrintedCost = 2

func drawnToTheDarkDimensionCost(s card.GameEngine) int {
	eff := drawnToTheDarkDimensionPrintedCost - s.Runechants()
	if eff < 0 {
		return 0
	}
	return eff
}

func (DrawnToTheDarkDimensionRed) Cost(s card.GameEngine) int { return drawnToTheDarkDimensionCost(s) }
func (DrawnToTheDarkDimensionRed) MinCost() int               { return 0 }
func (DrawnToTheDarkDimensionRed) MaxCost() int               { return drawnToTheDarkDimensionPrintedCost }

func (c DrawnToTheDarkDimensionRed) Play(s card.GameEngine, l card.Logger, self *card.CardState) {
	s.DrawOne()
}

func (DrawnToTheDarkDimensionYellow) Cost(s card.GameEngine) int {
	return drawnToTheDarkDimensionCost(s)
}
func (DrawnToTheDarkDimensionYellow) MinCost() int { return 0 }
func (DrawnToTheDarkDimensionYellow) MaxCost() int { return drawnToTheDarkDimensionPrintedCost }

func (c DrawnToTheDarkDimensionYellow) Play(s card.GameEngine, l card.Logger, self *card.CardState) {
	s.DrawOne()
}

func (DrawnToTheDarkDimensionBlue) Cost(s card.GameEngine) int { return drawnToTheDarkDimensionCost(s) }
func (DrawnToTheDarkDimensionBlue) MinCost() int               { return 0 }
func (DrawnToTheDarkDimensionBlue) MaxCost() int               { return drawnToTheDarkDimensionPrintedCost }

func (c DrawnToTheDarkDimensionBlue) Play(s card.GameEngine, l card.Logger, self *card.CardState) {
	s.DrawOne()
}
