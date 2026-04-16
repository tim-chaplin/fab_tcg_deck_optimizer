package runeblade

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

func TestMaleficIncantation_VerseCounterValue(t *testing.T) {
	// Each variant creates N Runechants on play. Play returns N (each token credited +1 at
	// creation) and state.Runechants tracks the tokens for any downstream consume or carryover.
	cases := []struct {
		c card.Card
		n int
	}{
		{MaleficIncantationRed{}, 3},
		{MaleficIncantationYellow{}, 2},
		{MaleficIncantationBlue{}, 1},
	}
	for _, tc := range cases {
		var s card.TurnState
		if got := tc.c.Play(&s); got != tc.n {
			t.Errorf("%s: Play() = %d, want %d", tc.c.Name(), got, tc.n)
		}
		if s.Runechants != tc.n {
			t.Errorf("%s: Runechants = %d, want %d", tc.c.Name(), s.Runechants, tc.n)
		}
	}
}
