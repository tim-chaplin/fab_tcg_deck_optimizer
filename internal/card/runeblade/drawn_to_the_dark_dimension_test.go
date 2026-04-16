package runeblade

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

// Compile-time: the Drawn to the Dark Dimension variants must implement DiscountPerRunechant,
// otherwise the solver will treat their cost as 0 instead of the printed 2.
var (
	_ card.DiscountPerRunechant = DrawnToTheDarkDimensionRed{}
	_ card.DiscountPerRunechant = DrawnToTheDarkDimensionYellow{}
	_ card.DiscountPerRunechant = DrawnToTheDarkDimensionBlue{}
)

func TestDrawnToTheDarkDimension_PrintedCost(t *testing.T) {
	// All three variants share printed cost 2; Cost() stays 0 for the partition-level minimum.
	cases := []card.Card{
		DrawnToTheDarkDimensionRed{},
		DrawnToTheDarkDimensionYellow{},
		DrawnToTheDarkDimensionBlue{},
	}
	for _, c := range cases {
		if c.Cost() != 0 {
			t.Errorf("%s: Cost() = %d, want 0 (minimum-cost floor)", c.Name(), c.Cost())
		}
		d, ok := c.(card.DiscountPerRunechant)
		if !ok {
			t.Fatalf("%s: does not implement DiscountPerRunechant", c.Name())
		}
		if d.PrintedCost() != 2 {
			t.Errorf("%s: PrintedCost() = %d, want 2", c.Name(), d.PrintedCost())
		}
	}
}
