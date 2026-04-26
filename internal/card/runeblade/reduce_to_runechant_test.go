package runeblade

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

// Compile-time: all three Reduce variants must implement card.VariableCost so the solver can
// pre-screen with MinCost / MaxCost bounds before running chain permutations.
var (
	_ card.VariableCost = ReduceToRunechantRed{}
	_ card.VariableCost = ReduceToRunechantYellow{}
	_ card.VariableCost = ReduceToRunechantBlue{}
)

func TestReduceToRunechant_PlayCreditsCreatedToken(t *testing.T) {
	cases := []card.Card{
		ReduceToRunechantRed{},
		ReduceToRunechantYellow{},
		ReduceToRunechantBlue{},
	}
	for _, c := range cases {
		s := &card.TurnState{}
		c.Play(s, &card.CardState{Card: c})
		got := s.Value
		if got != 1 {
			t.Errorf("%s: Play() = %d, want 1 (created Runechant credits +1)", c.Name(), got)
		}
		if s.Runechants != 1 {
			t.Errorf("%s: Runechants = %d, want 1 after Play", c.Name(), s.Runechants)
		}
	}
}

func TestReduceToRunechant_CostBounds(t *testing.T) {
	cases := []card.Card{
		ReduceToRunechantRed{},
		ReduceToRunechantYellow{},
		ReduceToRunechantBlue{},
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
