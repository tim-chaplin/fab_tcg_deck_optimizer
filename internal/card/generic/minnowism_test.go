package generic

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

// TestMinnowism_NoAttackReturnsZero: no qualifying next attack card → +3 rider fizzles.
func TestMinnowism_NoAttackReturnsZero(t *testing.T) {
	s := card.TurnState{}
	for _, c := range []card.Card{MinnowismRed{}, MinnowismYellow{}, MinnowismBlue{}} {
		if got := c.Play(&s); got != 0 {
			t.Errorf("%s: Play() = %d, want 0", c.Name(), got)
		}
	}
}

// TestMinnowism_HighPowerFilteredOut: a power-4 attack is seen but the power<=3 filter rejects it,
// so the rider fizzles without falling through to a later match.
func TestMinnowism_HighPowerFilteredOut(t *testing.T) {
	s := card.TurnState{CardsRemaining: []*card.PlayedCard{{Card: stubGenericAttack(0, 4)}}}
	if got := (MinnowismRed{}).Play(&s); got != 0 {
		t.Errorf("Play() = %d, want 0 (power 4 > 3)", got)
	}
}

// TestMinnowism_LowPowerReturnsBonus: first power<=3 attack triggers +3.
func TestMinnowism_LowPowerReturnsBonus(t *testing.T) {
	s := card.TurnState{CardsRemaining: []*card.PlayedCard{{Card: stubGenericAttack(0, 3)}}}
	for _, c := range []card.Card{MinnowismRed{}, MinnowismYellow{}, MinnowismBlue{}} {
		if got := c.Play(&s); got != 3 {
			t.Errorf("%s: Play() = %d, want 3", c.Name(), got)
		}
	}
}
