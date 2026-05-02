package cards

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

// Tests that with no arcane incoming the default branch credits 1{h}.
func TestArcanePolarity_NoArcaneIncomingCreditsOne(t *testing.T) {
	cases := []sim.Card{
		ArcanePolarityRed{},
		ArcanePolarityYellow{},
		ArcanePolarityBlue{},
	}
	for _, c := range cases {
		var s sim.TurnState
		c.Play(&s, &sim.CardState{Card: c})
		if s.Value != 1 {
			t.Errorf("%s: Value = %d, want 1", c.Name(), s.Value)
		}
	}
}

// Tests that ArcaneIncomingDamage > 0 swaps to the per-pitch alternate gain.
func TestArcanePolarity_ArcaneIncomingCreditsLargeGain(t *testing.T) {
	cases := []struct {
		c    sim.Card
		gain int
	}{
		{ArcanePolarityRed{}, 4},
		{ArcanePolarityYellow{}, 3},
		{ArcanePolarityBlue{}, 2},
	}
	for _, tc := range cases {
		s := sim.TurnState{ArcaneIncomingDamage: 1}
		tc.c.Play(&s, &sim.CardState{Card: tc.c})
		if s.Value != tc.gain {
			t.Errorf("%s: Value = %d, want %d", tc.c.Name(), s.Value, tc.gain)
		}
	}
}
