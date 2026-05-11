package cards

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
)

// Tests that Fate Foreseen blocks for printed defense and emits an Opt 1 log entry.
func TestFateForeseen_BlocksAndCallsOpt1(t *testing.T) {
	cases := []struct {
		c     sim.Card
		block int
	}{
		{FateForeseenRed{}, 4},
		{FateForeseenYellow{}, 3},
		{FateForeseenBlue{}, 2},
	}
	defer testutils.SwapCurrentHero(testutils.Hero{})()

	for _, tc := range cases {
		top := testutils.NewStubCard("top")
		s := sim.NewTurnStateFromCards([]sim.Card{top}, nil)
		s.IncomingDamage = 10
		sim.ResolveChainStep(s, s.Logger(), &sim.CardState{Card: tc.c})
		if s.Value != tc.block {
			t.Errorf("%s: Play(IncomingDamage=10) Value = %d, want %d (block only)",
				tc.c.Name(), s.Value, tc.block)
		}
		if len(s.LogEntries()) != 2 {
			t.Errorf("%s: Log len = %d, want 2 (Opted... + chain step)", tc.c.Name(), len(s.LogEntries()))
			continue
		}
		// Play emits the Opted line during the DR resolution; the chain step is
		// auto-appended after Play returns, so the Opted entry lands first.
		want := "Opted [top], put [top] on top, put [] on bottom"
		if got := s.LogEntries()[0].Text; got != want {
			t.Errorf("%s: Opt log entry = %q, want %q", tc.c.Name(), got, want)
		}
	}
}
