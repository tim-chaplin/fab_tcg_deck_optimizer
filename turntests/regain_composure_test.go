package turntests

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card/cards/notimplemented"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/gameengine"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
)

// TestRegainComposure_NoAttackReturnsZero: no qualifying next attack card → +1 rider fizzles.
func TestRegainComposure_NoAttackReturnsZero(t *testing.T) {
	ge := gameengine.New()
	ge.ResolveChainStep(ge.Logger(), &card.CardState{Card: notimplemented.RegainComposureBlue{}})
	if got := ge.Value(); got != 0 {
		t.Errorf("Play() = %d, want 0", got)
	}
}

// TestRegainComposure_NonAttackInRemainingFizzles: non-attack action fails the predicate.
func TestRegainComposure_NonAttackInRemainingFizzles(t *testing.T) {
	ge := &gameengine.GameEngine{GameState: gameengine.GameStateBuilder().SetCardsRemaining([]*card.CardState{{Card: testutils.GenericAction()}}).Build()}
	ge.ResolveChainStep(ge.Logger(), &card.CardState{Card: notimplemented.RegainComposureBlue{}})
	if got := ge.Value(); got != 0 {
		t.Errorf("Play() = %d, want 0 (non-attack skipped)", got)
	}
}

// TestRegainComposure_NextAttackGrantsBonusAttack: first attack-action picks up +1 on its
// BonusAttack so EffectiveAttack folds it into LikelyToHit. Granter returns 0.
func TestRegainComposure_NextAttackGrantsBonusAttack(t *testing.T) {
	target := &card.CardState{Card: testutils.GenericAttack(0, 0)}
	ge := &gameengine.GameEngine{GameState: gameengine.GameStateBuilder().SetCardsRemaining([]*card.CardState{target}).Build()}
	ge.ResolveChainStep(ge.Logger(), &card.CardState{Card: notimplemented.RegainComposureBlue{}})
	if got := ge.Value(); got != 0 {
		t.Errorf("Play() = %d, want 0 (granter returns 0; +N rides on target'ge BonusAttack)", got)
	}
	if target.BonusAttack != 1 {
		t.Errorf("target BonusAttack = %d, want 1", target.BonusAttack)
	}
}
