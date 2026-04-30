package cards

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

// Tests that On the Horizon blocks for its printed defense on each variant.
func TestOnTheHorizon_BlocksForPrintedDefense(t *testing.T) {
	cases := []struct {
		c    sim.Card
		want int
	}{
		{OnTheHorizonRed{}, 4},
		{OnTheHorizonYellow{}, 3},
		{OnTheHorizonBlue{}, 2},
	}
	for _, tc := range cases {
		s := sim.TurnState{IncomingDamage: 10}
		tc.c.Play(&s, &sim.CardState{Card: tc.c})
		if got := s.Value; got != tc.want {
			t.Errorf("%s: Play(IncomingDamage=10) Value = %d, want %d", tc.c.Name(), got, tc.want)
		}
	}
}
