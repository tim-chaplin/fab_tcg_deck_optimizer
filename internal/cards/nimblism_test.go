package cards

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
)

// TestNimblism_NoAttackReturnsZero: no qualifying next attack card → +3 rider fizzles.
func TestNimblism_NoAttackReturnsZero(t *testing.T) {
	s := card.TurnState{}
	for _, c := range []card.Card{NimblismRed{}, NimblismYellow{}, NimblismBlue{}} {
		c.Play(&s, &card.CardState{Card: c})
		if got := s.Value; got != 0 {
			t.Errorf("%s: Play() = %d, want 0", c.Name(), got)
		}
	}
}

// TestNimblism_HighCostFilteredOut: a cost-2 attack is seen but the cost<=1 filter rejects it.
func TestNimblism_HighCostFilteredOut(t *testing.T) {
	s := card.TurnState{CardsRemaining: []*card.CardState{{Card: testutils.GenericAttack(2, 0)}}}
	(NimblismRed{}).Play(&s, &card.CardState{Card: NimblismRed{}})
	if got := s.Value; got != 0 {
		t.Errorf("Play() = %d, want 0 (cost 2 > 1)", got)
	}
}

// TestNimblism_LowCostReturnsBonus: first cost<=1 attack triggers the per-variant bonus
// (Red +3, Yellow +2, Blue +1).
func TestNimblism_LowCostReturnsBonus(t *testing.T) {
	cases := []struct {
		c    card.Card
		want int
	}{
		{NimblismRed{}, 3},
		{NimblismYellow{}, 2},
		{NimblismBlue{}, 1},
	}
	for _, tc := range cases {
		target := &card.CardState{Card: testutils.GenericAttack(1, 0)}
		s := card.TurnState{CardsRemaining: []*card.CardState{target}}
		tc.c.Play(&s, &card.CardState{Card: tc.c})
		if got := s.Value; got != 0 {
			t.Errorf("%s: Play() = %d, want 0 (granter returns 0; +N rides on target's BonusAttack)", tc.c.Name(), got)
		}
		if target.BonusAttack != tc.want {
			t.Errorf("%s: target BonusAttack = %d, want %d", tc.c.Name(), target.BonusAttack, tc.want)
		}
	}
}
