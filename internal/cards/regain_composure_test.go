package cards

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

// TestRegainComposure_NoAttackReturnsZero: no qualifying next attack card → +1 rider fizzles.
func TestRegainComposure_NoAttackReturnsZero(t *testing.T) {
	s := card.TurnState{}
	(RegainComposureBlue{}).Play(&s, &card.CardState{Card: RegainComposureBlue{}})
	if got := s.Value; got != 0 {
		t.Errorf("Play() = %d, want 0", got)
	}
}

// TestRegainComposure_NonAttackInRemainingFizzles: non-attack action fails the predicate.
func TestRegainComposure_NonAttackInRemainingFizzles(t *testing.T) {
	s := card.TurnState{CardsRemaining: []*card.CardState{{Card: stubGenericAction()}}}
	(RegainComposureBlue{}).Play(&s, &card.CardState{Card: RegainComposureBlue{}})
	if got := s.Value; got != 0 {
		t.Errorf("Play() = %d, want 0 (non-attack skipped)", got)
	}
}

// TestRegainComposure_NextAttackGrantsBonusAttack: first attack-action picks up +1 on its
// BonusAttack so EffectiveAttack folds it into LikelyToHit. Granter returns 0.
func TestRegainComposure_NextAttackGrantsBonusAttack(t *testing.T) {
	target := &card.CardState{Card: stubGenericAttack(0, 0)}
	s := card.TurnState{CardsRemaining: []*card.CardState{target}}
	(RegainComposureBlue{}).Play(&s, &card.CardState{Card: RegainComposureBlue{}})
	if got := s.Value; got != 0 {
		t.Errorf("Play() = %d, want 0 (granter returns 0; +N rides on target's BonusAttack)", got)
	}
	if target.BonusAttack != 1 {
		t.Errorf("target BonusAttack = %d, want 1", target.BonusAttack)
	}
}
