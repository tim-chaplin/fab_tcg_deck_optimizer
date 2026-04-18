package generic

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

// TestCaptainsCall_NoAttackReturnsZero: no qualifying next attack card → +2 rider fizzles.
func TestCaptainsCall_NoAttackReturnsZero(t *testing.T) {
	s := card.TurnState{}
	for _, c := range []card.Card{CaptainsCallRed{}, CaptainsCallYellow{}, CaptainsCallBlue{}} {
		if got := c.Play(&s); got != 0 {
			t.Errorf("%s: Play() = %d, want 0", c.Name(), got)
		}
	}
}

// TestCaptainsCall_HighCostFilteredOut: a cost-3 attack is seen but the cost<=2 filter rejects it.
func TestCaptainsCall_HighCostFilteredOut(t *testing.T) {
	s := card.TurnState{CardsRemaining: []*card.PlayedCard{{Card: stubGenericAttack(3, 0)}}}
	if got := (CaptainsCallRed{}).Play(&s); got != 0 {
		t.Errorf("Play() = %d, want 0 (cost 3 > 2)", got)
	}
}

// TestCaptainsCall_LowCostReturnsBonus: first cost<=2 attack triggers +2.
func TestCaptainsCall_LowCostReturnsBonus(t *testing.T) {
	s := card.TurnState{CardsRemaining: []*card.PlayedCard{{Card: stubGenericAttack(2, 0)}}}
	for _, c := range []card.Card{CaptainsCallRed{}, CaptainsCallYellow{}, CaptainsCallBlue{}} {
		if got := c.Play(&s); got != 2 {
			t.Errorf("%s: Play() = %d, want 2", c.Name(), got)
		}
	}
}
