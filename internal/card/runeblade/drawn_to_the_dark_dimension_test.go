package runeblade

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

// Compile-time: the Drawn to the Dark Dimension variants must implement DiscountPerRunechant,
// otherwise the solver will not apply the Runechant discount at chain time.
var (
	_ card.DiscountPerRunechant = DrawnToTheDarkDimensionRed{}
	_ card.DiscountPerRunechant = DrawnToTheDarkDimensionYellow{}
	_ card.DiscountPerRunechant = DrawnToTheDarkDimensionBlue{}
)

func TestDrawnToTheDarkDimension_PrintedCost(t *testing.T) {
	cases := []card.Card{
		DrawnToTheDarkDimensionRed{},
		DrawnToTheDarkDimensionYellow{},
		DrawnToTheDarkDimensionBlue{},
	}
	for _, c := range cases {
		d, ok := c.(card.DiscountPerRunechant)
		if !ok {
			t.Fatalf("%s: does not implement DiscountPerRunechant", c.Name())
		}
		if d.PrintedCost() != 2 {
			t.Errorf("%s: PrintedCost() = %d, want 2", c.Name(), d.PrintedCost())
		}
		if c.Cost() != 2 {
			t.Errorf("%s: Cost() = %d, want 2", c.Name(), c.Cost())
		}
	}
}
