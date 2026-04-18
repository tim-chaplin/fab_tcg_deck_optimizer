package generic

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

// TestSloggism_NoAttackReturnsZero: no qualifying next attack card → +6 rider fizzles.
func TestSloggism_NoAttackReturnsZero(t *testing.T) {
	s := card.TurnState{}
	for _, c := range []card.Card{SloggismRed{}, SloggismYellow{}, SloggismBlue{}} {
		if got := c.Play(&s); got != 0 {
			t.Errorf("%s: Play() = %d, want 0", c.Name(), got)
		}
	}
}

// TestSloggism_LowCostFilteredOut: a cost-1 attack is seen but the cost>=2 filter rejects it.
func TestSloggism_LowCostFilteredOut(t *testing.T) {
	s := card.TurnState{CardsRemaining: []*card.PlayedCard{{Card: stubGenericAttack(1, 0)}}}
	if got := (SloggismRed{}).Play(&s); got != 0 {
		t.Errorf("Play() = %d, want 0 (cost 1 < 2)", got)
	}
}

// TestSloggism_HighCostReturnsBonus: first cost>=2 attack triggers +6.
func TestSloggism_HighCostReturnsBonus(t *testing.T) {
	s := card.TurnState{CardsRemaining: []*card.PlayedCard{{Card: stubGenericAttack(2, 0)}}}
	for _, c := range []card.Card{SloggismRed{}, SloggismYellow{}, SloggismBlue{}} {
		if got := c.Play(&s); got != 6 {
			t.Errorf("%s: Play() = %d, want 6", c.Name(), got)
		}
	}
}
