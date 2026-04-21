package generic

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

// TestCaptainsCall_NoAttackReturnsZero: no qualifying next attack card → +2 rider fizzles.
func TestCaptainsCall_NoAttackReturnsZero(t *testing.T) {
	s := card.TurnState{}
	for _, c := range []card.Card{CaptainsCallRed{}, CaptainsCallYellow{}, CaptainsCallBlue{}} {
		if got := c.Play(&s, &card.CardState{}); got != 0 {
			t.Errorf("%s: Play() = %d, want 0", c.Name(), got)
		}
	}
}

// TestCaptainsCall_HighCostFilteredOut: a cost-3 attack is past Red's cost<=2 gate.
func TestCaptainsCall_HighCostFilteredOut(t *testing.T) {
	s := card.TurnState{CardsRemaining: []*card.CardState{{Card: stubGenericAttack(3, 0)}}}
	if got := (CaptainsCallRed{}).Play(&s, &card.CardState{}); got != 0 {
		t.Errorf("Play() = %d, want 0 (cost 3 > 2)", got)
	}
}

// TestCaptainsCall_CostThresholdPerVariant: the +2 bonus is flat across variants, but each
// variant has its own cost threshold: Red cost<=2, Yellow cost<=1, Blue cost==0. A cost-2 attack
// only triggers Red; a cost-1 attack triggers Red and Yellow; a cost-0 attack triggers all three.
func TestCaptainsCall_CostThresholdPerVariant(t *testing.T) {
	cases := []struct {
		name    string
		cost    int
		red     int
		yellow  int
		blue    int
	}{
		{"cost 2", 2, 2, 0, 0},
		{"cost 1", 1, 2, 2, 0},
		{"cost 0", 0, 2, 2, 2},
	}
	for _, tc := range cases {
		s := card.TurnState{CardsRemaining: []*card.CardState{{Card: stubGenericAttack(tc.cost, 0)}}}
		if got := (CaptainsCallRed{}).Play(&s, &card.CardState{}); got != tc.red {
			t.Errorf("%s Red: Play() = %d, want %d", tc.name, got, tc.red)
		}
		if got := (CaptainsCallYellow{}).Play(&s, &card.CardState{}); got != tc.yellow {
			t.Errorf("%s Yellow: Play() = %d, want %d", tc.name, got, tc.yellow)
		}
		if got := (CaptainsCallBlue{}).Play(&s, &card.CardState{}); got != tc.blue {
			t.Errorf("%s Blue: Play() = %d, want %d", tc.name, got, tc.blue)
		}
	}
}
