package runeblade

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

// TestRunebloodIncantation_BlockCoversIncomingReturnsN pins the N-per-variant payoff when the
// partition's block total covers incoming damage — the aura survives to tick its verse counters.
func TestRunebloodIncantation_BlockCoversIncomingReturnsN(t *testing.T) {
	cases := []struct {
		c card.Card
		n int
	}{
		{RunebloodIncantationRed{}, 3},
		{RunebloodIncantationYellow{}, 2},
		{RunebloodIncantationBlue{}, 1},
	}
	for _, tc := range cases {
		s := card.TurnState{IncomingDamage: 3, BlockTotal: 3}
		if got := tc.c.Play(&s); got != tc.n {
			t.Errorf("%s: Play() = %d, want %d (block == incoming)", tc.c.Name(), got, tc.n)
		}
	}
}

// TestRunebloodIncantation_BlockShortReturnsZero pins the collapse-to-0 case when incoming
// damage gets through.
func TestRunebloodIncantation_BlockShortReturnsZero(t *testing.T) {
	cases := []card.Card{
		RunebloodIncantationRed{},
		RunebloodIncantationYellow{},
		RunebloodIncantationBlue{},
	}
	for _, c := range cases {
		s := card.TurnState{IncomingDamage: 3, BlockTotal: 2}
		if got := c.Play(&s); got != 0 {
			t.Errorf("%s: Play() = %d, want 0 (block < incoming)", c.Name(), got)
		}
	}
}
