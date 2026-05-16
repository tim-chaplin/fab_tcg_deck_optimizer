package turntests

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/cards"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/gameengine"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/token"
)

// Compile-time: the Drawn to the Dark Dimension variants must implement card.VariableCost.
var (
	_ card.VariableCost = cards.DrawnToTheDarkDimensionRed{}
	_ card.VariableCost = cards.DrawnToTheDarkDimensionYellow{}
	_ card.VariableCost = cards.DrawnToTheDarkDimensionBlue{}
)

func TestDrawnToTheDarkDimension_CostBounds(t *testing.T) {
	cases := []card.Card{
		cards.DrawnToTheDarkDimensionRed{},
		cards.DrawnToTheDarkDimensionYellow{},
		cards.DrawnToTheDarkDimensionBlue{},
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
		if c.Cost(gameengine.New()) != 2 {
			t.Errorf("%s: Cost(zeroState) = %d, want 2", c.Name(), c.Cost(gameengine.New()))
		}
		withRune := gameengine.New()
		withRune.CreateAura(token.NewRunechant(5))
		if c.Cost(withRune) != 0 {
			t.Errorf("%s: Cost(Runechants=5) = %d, want 0", c.Name(), c.Cost(withRune))
		}
	}
}
