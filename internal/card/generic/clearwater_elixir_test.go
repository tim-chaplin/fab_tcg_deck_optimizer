package generic

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

// TestClearwaterElixir_NoAttackReturnsZero: no qualifying next attack card → +3 rider fizzles.
func TestClearwaterElixir_NoAttackReturnsZero(t *testing.T) {
	s := card.TurnState{}
	if got := (ClearwaterElixirRed{}).Play(&s, &card.CardState{}); got != 0 {
		t.Errorf("Play() = %d, want 0", got)
	}
}

// TestClearwaterElixir_NonAttackInRemainingFizzles: non-attack action fails the predicate.
func TestClearwaterElixir_NonAttackInRemainingFizzles(t *testing.T) {
	s := card.TurnState{CardsRemaining: []*card.CardState{{Card: stubGenericAction()}}}
	if got := (ClearwaterElixirRed{}).Play(&s, &card.CardState{}); got != 0 {
		t.Errorf("Play() = %d, want 0 (non-attack skipped)", got)
	}
}

// TestClearwaterElixir_NextAttackGrantsBonusAttack: first attack-action picks up +3 on its
// BonusAttack so EffectiveAttack folds it into LikelyToHit and the solver routes the bonus
// to the buffed attack's chain slot. Granter returns 0 — the +3 attributes to the target.
func TestClearwaterElixir_NextAttackGrantsBonusAttack(t *testing.T) {
	target := &card.CardState{Card: stubGenericAttack(0, 0)}
	s := card.TurnState{CardsRemaining: []*card.CardState{target}}
	if got := (ClearwaterElixirRed{}).Play(&s, &card.CardState{}); got != 0 {
		t.Errorf("Play() = %d, want 0 (granter returns 0; +N rides on target's BonusAttack)", got)
	}
	if target.BonusAttack != 3 {
		t.Errorf("target BonusAttack = %d, want 3", target.BonusAttack)
	}
}
