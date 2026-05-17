package turntests

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card/cards"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/gameengine"
)

// Compile-time: all three Reduce variants must implement card.VariableCost.
var (
	_ card.VariableCost = cards.ReduceToRunechantRed{}
	_ card.VariableCost = cards.ReduceToRunechantYellow{}
	_ card.VariableCost = cards.ReduceToRunechantBlue{}
)

func TestReduceToRunechant_PlayCreditsCreatedToken(t *testing.T) {
	cases := []card.Card{
		cards.ReduceToRunechantRed{},
		cards.ReduceToRunechantYellow{},
		cards.ReduceToRunechantBlue{},
	}
	for _, c := range cases {
		ge := gameengine.New()
		ge.ResolveChainStep(ge.Logger(), &card.CardState{Card: c})
		got := ge.Value()
		if got != 1 {
			t.Errorf("%s: Play() = %d, want 1 (created Runechant credits +1)", c.Name(), got)
		}
		if ge.RunechantCount() != 1 {
			t.Errorf("%s: Runechants = %d, want 1 after Play", c.Name(), ge.RunechantCount())
		}
	}
}

func TestReduceToRunechant_CostBounds(t *testing.T) {
	cases := []card.Card{
		cards.ReduceToRunechantRed{},
		cards.ReduceToRunechantYellow{},
		cards.ReduceToRunechantBlue{},
	}
	for _, c := range cases {
		vc, ok := c.(card.VariableCost)
		if !ok {
			t.Fatalf("%s: does not implement card.VariableCost", c.Name())
		}
		if vc.MaxCost() != 1 {
			t.Errorf("%s: MaxCost() = %d, want 1", c.Name(), vc.MaxCost())
		}
		if vc.MinCost() != 0 {
			t.Errorf("%s: MinCost() = %d, want 0", c.Name(), vc.MinCost())
		}
	}
}
