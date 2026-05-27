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
		vc, ok := c.(card.VariableCost)
		if !ok {
			t.Fatalf("%s: does not implement card.VariableCost", c.Name())
		}
		zeroState := gameengine.New()
		if got := vc.EffectiveCost(zeroState); got != c.Cost() {
			t.Errorf("%s: EffectiveCost(zeroState) = %d, want printed Cost()=%d", c.Name(), got, c.Cost())
		}
		withRune := &gameengine.GameEngine{GameState: gameengine.GameStateBuilder().
			AddAura(token.NewRunechant(5)).
			Build()}
		if got := vc.EffectiveCost(withRune); got != 0 {
			t.Errorf("%s: EffectiveCost(Runechants=5) = %d, want 0 (discount floors at 0)", c.Name(), got)
		}
	}
}
