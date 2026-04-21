package generic

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

// TestRestvineElixir_NoAttackReturnsZero: no qualifying next attack card → +3 rider fizzles.
func TestRestvineElixir_NoAttackReturnsZero(t *testing.T) {
	s := card.TurnState{}
	if got := (RestvineElixirRed{}).Play(&s, nil); got != 0 {
		t.Errorf("Play() = %d, want 0", got)
	}
}

// TestRestvineElixir_NonAttackInRemainingFizzles: non-attack action fails the predicate.
func TestRestvineElixir_NonAttackInRemainingFizzles(t *testing.T) {
	s := card.TurnState{CardsRemaining: []*card.PlayedCard{{Card: stubGenericAction()}}}
	if got := (RestvineElixirRed{}).Play(&s, nil); got != 0 {
		t.Errorf("Play() = %d, want 0 (non-attack skipped)", got)
	}
}

// TestRestvineElixir_NextAttackReturnsBonus: first attack-action triggers +3.
func TestRestvineElixir_NextAttackReturnsBonus(t *testing.T) {
	s := card.TurnState{CardsRemaining: []*card.PlayedCard{{Card: stubGenericAttack(0, 0)}}}
	if got := (RestvineElixirRed{}).Play(&s, nil); got != 3 {
		t.Errorf("Play() = %d, want 3", got)
	}
}
