package cards

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

// Tests that each printing prevents its full Defense() amount (4/3/2).
func TestPeaceOfMind_PreventsByPrinting(t *testing.T) {
	cases := []struct {
		card sim.Card
		want int
	}{
		{PeaceOfMindRed{}, 4},
		{PeaceOfMindYellow{}, 3},
		{PeaceOfMindBlue{}, 2},
	}
	for _, tc := range cases {
		s := sim.TurnState{IncomingDamage: 10}
		self := &sim.CardState{Card: tc.card}
		tc.card.Play(&s, s.Logger(), self)
		if s.Value != tc.want {
			t.Errorf("%s: Value = %d, want %d", tc.card.Name(), s.Value, tc.want)
		}
	}
}

// Tests that each printing creates one Ponder token on play.
func TestPeaceOfMind_CreatesPonder(t *testing.T) {
	for _, c := range []sim.Card{PeaceOfMindRed{}, PeaceOfMindYellow{}, PeaceOfMindBlue{}} {
		s := sim.TurnState{IncomingDamage: 10}
		self := &sim.CardState{Card: c}
		c.Play(&s, s.Logger(), self)
		if got := s.Ponders(); got != 1 {
			t.Errorf("%s: Ponders = %d, want 1", c.Name(), got)
		}
		if !s.AuraCreated {
			t.Errorf("%s: AuraCreated = false, want true after creating a Ponder", c.Name())
		}
	}
}
