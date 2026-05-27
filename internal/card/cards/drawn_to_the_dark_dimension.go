// Drawn to the Dark Dimension — Runeblade Action - Attack. Printed cost 2, Defense 3.
// Printed pitch variants: Red 1, Yellow 2, Blue 3.
// Printed power: Red 3, Yellow 2, Blue 1.
// Text: "Drawn to the Dark Dimension costs {r} less to play for each Runechant you control.
// Draw a card."

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

const drawnToTheDarkDimensionPrintedCost = 2

// drawnToTheDarkDimensionEffectiveCost is max(0, printed - RunechantCount).
func drawnToTheDarkDimensionEffectiveCost(ge card.GameEngine) int {
	eff := drawnToTheDarkDimensionPrintedCost - ge.RunechantCount()
	if eff < 0 {
		return 0
	}
	return eff
}

func (DrawnToTheDarkDimensionRed) EffectiveCost(ge card.GameEngine) int {
	return drawnToTheDarkDimensionEffectiveCost(ge)
}
func (DrawnToTheDarkDimensionRed) MinCost() int { return 0 }

func (c DrawnToTheDarkDimensionRed) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	ge.DrawOne()
}

func (DrawnToTheDarkDimensionYellow) EffectiveCost(ge card.GameEngine) int {
	return drawnToTheDarkDimensionEffectiveCost(ge)
}
func (DrawnToTheDarkDimensionYellow) MinCost() int { return 0 }

func (c DrawnToTheDarkDimensionYellow) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	ge.DrawOne()
}

func (DrawnToTheDarkDimensionBlue) EffectiveCost(ge card.GameEngine) int {
	return drawnToTheDarkDimensionEffectiveCost(ge)
}
func (DrawnToTheDarkDimensionBlue) MinCost() int { return 0 }

func (c DrawnToTheDarkDimensionBlue) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	ge.DrawOne()
}
