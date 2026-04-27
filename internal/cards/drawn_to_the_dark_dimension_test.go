package cards

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

// Compile-time: the Drawn to the Dark Dimension variants must implement card.VariableCost,
// otherwise the solver can't pre-screen with MinCost / MaxCost bounds.
var (
	_ card.VariableCost = DrawnToTheDarkDimensionRed{}
	_ card.VariableCost = DrawnToTheDarkDimensionYellow{}
	_ card.VariableCost = DrawnToTheDarkDimensionBlue{}
)

func TestDrawnToTheDarkDimension_CostBounds(t *testing.T) {
	cases := []card.Card{
		DrawnToTheDarkDimensionRed{},
		DrawnToTheDarkDimensionYellow{},
		DrawnToTheDarkDimensionBlue{},
	}
	for _, c := range cases {
		vc, ok := c.(card.VariableCost)
		if !ok {
			t.Fatalf("%s: does not implement card.VariableCost", c.Name())
		}
		if vc.MaxCost() != 2 {
			t.Errorf("%s: MaxCost() = %d, want 2", c.Name(), vc.MaxCost())
		}
		if vc.MinCost() != 0 {
			t.Errorf("%s: MinCost() = %d, want 0", c.Name(), vc.MinCost())
		}
		if c.Cost(&card.TurnState{}) != 2 {
			t.Errorf("%s: Cost(zeroState) = %d, want 2", c.Name(), c.Cost(&card.TurnState{}))
		}
		if c.Cost(&card.TurnState{Runechants: 5}) != 0 {
			t.Errorf("%s: Cost(Runechants=5) = %d, want 0", c.Name(), c.Cost(&card.TurnState{Runechants: 5}))
		}
	}
}
