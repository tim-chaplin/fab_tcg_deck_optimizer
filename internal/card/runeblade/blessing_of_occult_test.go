package runeblade

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

// TestBlessingOfOccult_PushesTokensToNextTurn confirms Blessing of Occult routes its Runechants
// through DelayedRunechants rather than live Runechants — same-turn attacks in a chain won't
// consume them, but they end up in LeftoverRunechants for the next turn.
func TestBlessingOfOccult_PushesTokensToNextTurn(t *testing.T) {
	cases := []struct {
		c card.Card
		n int
	}{
		{BlessingOfOccultRed{}, 3},
		{BlessingOfOccultYellow{}, 2},
		{BlessingOfOccultBlue{}, 1},
	}
	for _, tc := range cases {
		var s card.TurnState
		if got := tc.c.Play(&s); got != tc.n {
			t.Errorf("%s: Play() = %d, want %d", tc.c.Name(), got, tc.n)
		}
		if s.Runechants != 0 {
			t.Errorf("%s: Runechants = %d, want 0 (tokens skip this turn)", tc.c.Name(), s.Runechants)
		}
		if s.DelayedRunechants != tc.n {
			t.Errorf("%s: DelayedRunechants = %d, want %d", tc.c.Name(), s.DelayedRunechants, tc.n)
		}
		if !s.AuraCreated {
			t.Errorf("%s: AuraCreated should still be set", tc.c.Name())
		}
	}
}
