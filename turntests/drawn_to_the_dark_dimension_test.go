package turntests

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card/cards"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/gameengine"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/token"
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
		if c.Cost() != 2 {
			t.Errorf("%s: Cost() = %d, want 2 (printed cost)", c.Name(), c.Cost())
		}
		vc, ok := c.(card.VariableCost)
		if !ok {
			t.Fatalf("%s: does not implement card.VariableCost", c.Name())
		}
		if vc.MinCost() != 0 {
			t.Errorf("%s: MinCost() = %d, want 0", c.Name(), vc.MinCost())
		}
		zeroState := gameengine.New()
		if vc.EffectiveCost(zeroState) != 2 {
			t.Errorf("%s: EffectiveCost(zeroState) = %d, want 2", c.Name(), vc.EffectiveCost(zeroState))
		}
		withRune := &gameengine.GameEngine{GameState: gameengine.GameStateBuilder().
			AddAura(token.NewRunechant(5)).
			Build()}
		if vc.EffectiveCost(withRune) != 0 {
			t.Errorf("%s: EffectiveCost(Runechants=5) = %d, want 0", c.Name(), vc.EffectiveCost(withRune))
		}
	}
}
