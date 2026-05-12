package turntests

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/cards"
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
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
		s := &sim.TurnState{}
		sim.ResolveChainStep(s, s.Logger(), &card.CardState{Card: c})
		got := s.Value()
		if got != 1 {
			t.Errorf("%s: Play() = %d, want 1 (created Runechant credits +1)", c.Name(), got)
		}
		if s.Runechants() != 1 {
			t.Errorf("%s: Runechants = %d, want 1 after Play", c.Name(), s.Runechants())
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
