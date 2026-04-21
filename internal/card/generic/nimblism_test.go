package generic

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

// TestNimblism_NoAttackReturnsZero: no qualifying next attack card → +3 rider fizzles.
func TestNimblism_NoAttackReturnsZero(t *testing.T) {
	s := card.TurnState{}
	for _, c := range []card.Card{NimblismRed{}, NimblismYellow{}, NimblismBlue{}} {
		if got := c.Play(&s, &card.CardState{}); got != 0 {
			t.Errorf("%s: Play() = %d, want 0", c.Name(), got)
		}
	}
}

// TestNimblism_HighCostFilteredOut: a cost-2 attack is seen but the cost<=1 filter rejects it.
func TestNimblism_HighCostFilteredOut(t *testing.T) {
	s := card.TurnState{CardsRemaining: []*card.CardState{{Card: stubGenericAttack(2, 0)}}}
	if got := (NimblismRed{}).Play(&s, &card.CardState{}); got != 0 {
		t.Errorf("Play() = %d, want 0 (cost 2 > 1)", got)
	}
}

// TestNimblism_LowCostReturnsBonus: first cost<=1 attack triggers the per-variant bonus
// (Red +3, Yellow +2, Blue +1).
func TestNimblism_LowCostReturnsBonus(t *testing.T) {
	s := card.TurnState{CardsRemaining: []*card.CardState{{Card: stubGenericAttack(1, 0)}}}
	cases := []struct {
		c    card.Card
		want int
	}{
		{NimblismRed{}, 3},
		{NimblismYellow{}, 2},
		{NimblismBlue{}, 1},
	}
	for _, tc := range cases {
		if got := tc.c.Play(&s, &card.CardState{}); got != tc.want {
			t.Errorf("%s: Play() = %d, want %d", tc.c.Name(), got, tc.want)
		}
	}
}
