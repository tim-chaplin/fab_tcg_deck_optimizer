package runeblade

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

// Compile-time: all three Reduce variants must implement DiscountPerRunechant so the solver
// uses the per-runechant discount path instead of treating their Cost() as authoritative.
var (
	_ card.DiscountPerRunechant = ReduceToRunechantRed{}
	_ card.DiscountPerRunechant = ReduceToRunechantYellow{}
	_ card.DiscountPerRunechant = ReduceToRunechantBlue{}
)

func TestReduceToRunechant_PlayCreditsCreatedToken(t *testing.T) {
	cases := []card.Card{
		ReduceToRunechantRed{},
		ReduceToRunechantYellow{},
		ReduceToRunechantBlue{},
	}
	for _, c := range cases {
		s := &card.TurnState{}
		got := c.Play(s)
		if got != 1 {
			t.Errorf("%s: Play() = %d, want 1 (created Runechant credits +1)", c.Name(), got)
		}
		if s.Runechants != 1 {
			t.Errorf("%s: Runechants = %d, want 1 after Play", c.Name(), s.Runechants)
		}
	}
}

func TestReduceToRunechant_PrintedCost(t *testing.T) {
	cases := []card.DiscountPerRunechant{
		ReduceToRunechantRed{},
		ReduceToRunechantYellow{},
		ReduceToRunechantBlue{},
	}
	for _, c := range cases {
		if got := c.PrintedCost(); got != 1 {
			t.Errorf("%s: PrintedCost() = %d, want 1", c.(card.Card).Name(), got)
		}
	}
}
