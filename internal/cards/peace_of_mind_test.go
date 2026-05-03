package cards

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

// Tests that each printing prevents its full Defense() amount (4/3/2). The Ponder token
// rider is intentionally not credited — token economies aren't modelled.
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
		tc.card.Play(&s, self)
		if s.Value != tc.want {
			t.Errorf("%s: Value = %d, want %d", tc.card.Name(), s.Value, tc.want)
		}
	}
}

// Tests that the printed cost is 2 across all printings.
func TestPeaceOfMind_CostIsTwo(t *testing.T) {
	cases := []sim.Card{PeaceOfMindRed{}, PeaceOfMindYellow{}, PeaceOfMindBlue{}}
	for _, c := range cases {
		if got := c.Cost(&sim.TurnState{}); got != 2 {
			t.Errorf("%s: Cost = %d, want 2", c.Name(), got)
		}
	}
}
