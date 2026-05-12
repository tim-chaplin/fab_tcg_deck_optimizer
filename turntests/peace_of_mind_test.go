package turntests

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/cards"
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
)

// Tests that each printing prevents its full Defense() amount (4/3/2).
func TestPeaceOfMind_PreventsByPrinting(t *testing.T) {
	cases := []struct {
		card card.Card
		want int
	}{
		{cards.PeaceOfMindRed{}, 4},
		{cards.PeaceOfMindYellow{}, 3},
		{cards.PeaceOfMindBlue{}, 2},
	}
	for _, tc := range cases {
		s := sim.NewTurnStateFromSpec(sim.TurnStateSpec{IncomingDamage: 10})
		self := &card.CardState{Card: tc.card}
		sim.ResolveChainStep(&s, s.Logger(), self)
		if s.Value() != tc.want {
			t.Errorf("%s: Value = %d, want %d", tc.card.Name(), s.Value(), tc.want)
		}
	}
}

// Tests that each printing creates one Ponder token on play.
func TestPeaceOfMind_CreatesPonder(t *testing.T) {
	for _, c := range []card.Card{cards.PeaceOfMindRed{}, cards.PeaceOfMindYellow{}, cards.PeaceOfMindBlue{}} {
		s := sim.NewTurnStateFromSpec(sim.TurnStateSpec{IncomingDamage: 10})
		self := &card.CardState{Card: c}
		sim.ResolveChainStep(&s, s.Logger(), self)
		if got := s.Ponders(); got != 1 {
			t.Errorf("%s: Ponders = %d, want 1", c.Name(), got)
		}
		if !s.AuraCreated() {
			t.Errorf("%s: AuraCreated = false, want true after creating a Ponder", c.Name())
		}
	}
}
