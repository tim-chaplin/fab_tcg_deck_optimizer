package cards

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

// Tests that each printing prevents its full Defense() amount (4/3/2).
func TestOasisRespite_PreventsByPrinting(t *testing.T) {
	cases := []struct {
		card sim.Card
		want int
	}{
		{OasisRespiteRed{}, 4},
		{OasisRespiteYellow{}, 3},
		{OasisRespiteBlue{}, 2},
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
