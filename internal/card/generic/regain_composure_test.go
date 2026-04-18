package generic

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

// TestRegainComposure_NoAttackReturnsZero: no qualifying next attack card → +1 rider fizzles.
func TestRegainComposure_NoAttackReturnsZero(t *testing.T) {
	s := card.TurnState{}
	if got := (RegainComposureBlue{}).Play(&s); got != 0 {
		t.Errorf("Play() = %d, want 0", got)
	}
}

// TestRegainComposure_NonAttackInRemainingFizzles: non-attack action fails the predicate.
func TestRegainComposure_NonAttackInRemainingFizzles(t *testing.T) {
	s := card.TurnState{CardsRemaining: []*card.PlayedCard{{Card: stubGenericAction()}}}
	if got := (RegainComposureBlue{}).Play(&s); got != 0 {
		t.Errorf("Play() = %d, want 0 (non-attack skipped)", got)
	}
}

// TestRegainComposure_NextAttackReturnsBonus: first attack-action triggers +1.
func TestRegainComposure_NextAttackReturnsBonus(t *testing.T) {
	s := card.TurnState{CardsRemaining: []*card.PlayedCard{{Card: stubGenericAttack(0, 0)}}}
	if got := (RegainComposureBlue{}).Play(&s); got != 1 {
		t.Errorf("Play() = %d, want 1", got)
	}
}
