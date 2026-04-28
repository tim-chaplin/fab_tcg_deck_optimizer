package cards

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

// Compile-time: all three Reduce variants must implement sim.VariableCost so the solver can
// pre-screen with MinCost / MaxCost bounds before running chain permutations.
var (
	_ sim.VariableCost = ReduceToRunechantRed{}
	_ sim.VariableCost = ReduceToRunechantYellow{}
	_ sim.VariableCost = ReduceToRunechantBlue{}
)

func TestReduceToRunechant_PlayCreditsCreatedToken(t *testing.T) {
	cases := []sim.Card{
		ReduceToRunechantRed{},
		ReduceToRunechantYellow{},
		ReduceToRunechantBlue{},
	}
	for _, c := range cases {
		s := &sim.TurnState{}
		c.Play(s, &sim.CardState{Card: c})
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
	cases := []sim.Card{
		ReduceToRunechantRed{},
		ReduceToRunechantYellow{},
		ReduceToRunechantBlue{},
	}
	for _, c := range cases {
		vc, ok := c.(sim.VariableCost)
		if !ok {
			t.Fatalf("%s: does not implement sim.VariableCost", c.Name())
		}
		if vc.MaxCost() != 1 {
			t.Errorf("%s: MaxCost() = %d, want 1", c.Name(), vc.MaxCost())
		}
		if vc.MinCost() != 0 {
			t.Errorf("%s: MinCost() = %d, want 0", c.Name(), vc.MinCost())
		}
	}
}
