package cards

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

// Tests that Fate Foreseen blocks for printed defense and credits Opt 1 on top.
func TestFateForeseen_BlocksAndCreditsOpt1(t *testing.T) {
	cases := []struct {
		c     sim.Card
		block int
	}{
		{FateForeseenRed{}, 4},
		{FateForeseenYellow{}, 3},
		{FateForeseenBlue{}, 2},
	}
	for _, tc := range cases {
		s := sim.TurnState{IncomingDamage: 10}
		tc.c.Play(&s, &sim.CardState{Card: tc.c})
		want := tc.block + sim.OptValue
		if got := s.Value; got != want {
			t.Errorf("%s: Play(IncomingDamage=10) Value = %d, want %d (block %d + Opt 1)",
				tc.c.Name(), got, want, tc.block)
		}
	}
}
