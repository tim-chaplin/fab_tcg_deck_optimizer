package generic

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

// TestClearwaterElixir_NoAttackReturnsZero: no qualifying next attack card → +3 rider fizzles.
func TestClearwaterElixir_NoAttackReturnsZero(t *testing.T) {
	s := card.TurnState{}
	if got := (ClearwaterElixirRed{}).Play(&s); got != 0 {
		t.Errorf("Play() = %d, want 0", got)
	}
}

// TestClearwaterElixir_NonAttackInRemainingFizzles: non-attack action fails the predicate.
func TestClearwaterElixir_NonAttackInRemainingFizzles(t *testing.T) {
	s := card.TurnState{CardsRemaining: []*card.PlayedCard{{Card: stubGenericAction()}}}
	if got := (ClearwaterElixirRed{}).Play(&s); got != 0 {
		t.Errorf("Play() = %d, want 0 (non-attack skipped)", got)
	}
}

// TestClearwaterElixir_NextAttackReturnsBonus: first attack-action triggers +3.
func TestClearwaterElixir_NextAttackReturnsBonus(t *testing.T) {
	s := card.TurnState{CardsRemaining: []*card.PlayedCard{{Card: stubGenericAttack(0, 0)}}}
	if got := (ClearwaterElixirRed{}).Play(&s); got != 3 {
		t.Errorf("Play() = %d, want 3", got)
	}
}
