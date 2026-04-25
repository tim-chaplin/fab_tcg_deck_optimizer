package generic

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

// TestSapwoodElixir_NoAttackReturnsZero: no qualifying next attack card → +3 rider fizzles.
func TestSapwoodElixir_NoAttackReturnsZero(t *testing.T) {
	s := card.TurnState{}
	if got := (SapwoodElixirRed{}).Play(&s, &card.CardState{}); got != 0 {
		t.Errorf("Play() = %d, want 0", got)
	}
}

// TestSapwoodElixir_NonAttackInRemainingFizzles: non-attack action fails the predicate.
func TestSapwoodElixir_NonAttackInRemainingFizzles(t *testing.T) {
	s := card.TurnState{CardsRemaining: []*card.CardState{{Card: stubGenericAction()}}}
	if got := (SapwoodElixirRed{}).Play(&s, &card.CardState{}); got != 0 {
		t.Errorf("Play() = %d, want 0 (non-attack skipped)", got)
	}
}

// TestSapwoodElixir_NextAttackGrantsBonusAttack: first attack-action picks up +3 on its
// BonusAttack. Granter returns 0; the +3 attributes to the target.
func TestSapwoodElixir_NextAttackGrantsBonusAttack(t *testing.T) {
	target := &card.CardState{Card: stubGenericAttack(0, 0)}
	s := card.TurnState{CardsRemaining: []*card.CardState{target}}
	if got := (SapwoodElixirRed{}).Play(&s, &card.CardState{}); got != 0 {
		t.Errorf("Play() = %d, want 0 (granter returns 0; +N rides on target's BonusAttack)", got)
	}
	if target.BonusAttack != 3 {
		t.Errorf("target BonusAttack = %d, want 3", target.BonusAttack)
	}
}
