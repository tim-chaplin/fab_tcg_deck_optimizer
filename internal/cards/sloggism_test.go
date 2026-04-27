package cards

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

// TestSloggism_NoAttackReturnsZero: no qualifying next attack card → +6 rider fizzles.
func TestSloggism_NoAttackReturnsZero(t *testing.T) {
	s := card.TurnState{}
	for _, c := range []card.Card{SloggismRed{}, SloggismYellow{}, SloggismBlue{}} {
		c.Play(&s, &card.CardState{Card: c})
		if got := s.Value; got != 0 {
			t.Errorf("%s: Play() = %d, want 0", c.Name(), got)
		}
	}
}

// TestSloggism_LowCostFilteredOut: a cost-1 attack is seen but the cost>=2 filter rejects it.
func TestSloggism_LowCostFilteredOut(t *testing.T) {
	s := card.TurnState{CardsRemaining: []*card.CardState{{Card: stubGenericAttack(1, 0)}}}
	(SloggismRed{}).Play(&s, &card.CardState{Card: SloggismRed{}})
	if got := s.Value; got != 0 {
		t.Errorf("Play() = %d, want 0 (cost 1 < 2)", got)
	}
}

// TestSloggism_HighCostReturnsBonus: first cost>=2 attack triggers the per-variant bonus
// (Red +6, Yellow +5, Blue +4).
func TestSloggism_HighCostReturnsBonus(t *testing.T) {
	cases := []struct {
		c    card.Card
		want int
	}{
		{SloggismRed{}, 6},
		{SloggismYellow{}, 5},
		{SloggismBlue{}, 4},
	}
	for _, tc := range cases {
		target := &card.CardState{Card: stubGenericAttack(2, 0)}
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
