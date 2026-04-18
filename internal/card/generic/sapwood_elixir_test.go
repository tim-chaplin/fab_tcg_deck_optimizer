package generic

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

// TestSapwoodElixir_NoAttackReturnsZero: no qualifying next attack card → +3 rider fizzles.
func TestSapwoodElixir_NoAttackReturnsZero(t *testing.T) {
	s := card.TurnState{}
	if got := (SapwoodElixirRed{}).Play(&s); got != 0 {
		t.Errorf("Play() = %d, want 0", got)
	}
}

// TestSapwoodElixir_NonAttackInRemainingFizzles: non-attack action fails the predicate.
func TestSapwoodElixir_NonAttackInRemainingFizzles(t *testing.T) {
	s := card.TurnState{CardsRemaining: []*card.PlayedCard{{Card: stubGenericAction()}}}
	if got := (SapwoodElixirRed{}).Play(&s); got != 0 {
		t.Errorf("Play() = %d, want 0 (non-attack skipped)", got)
	}
}

// TestSapwoodElixir_NextAttackReturnsBonus: first attack-action triggers +3.
func TestSapwoodElixir_NextAttackReturnsBonus(t *testing.T) {
	s := card.TurnState{CardsRemaining: []*card.PlayedCard{{Card: stubGenericAttack(0, 0)}}}
	if got := (SapwoodElixirRed{}).Play(&s); got != 3 {
		t.Errorf("Play() = %d, want 3", got)
	}
}
