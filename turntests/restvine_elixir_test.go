package turntests

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/v2/card/cards/notimplemented"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/gameengine"
)

// TestRestvineElixir_NoAttackReturnsZero: no qualifying next attack card → +3 rider fizzles.
func TestRestvineElixir_NoAttackReturnsZero(t *testing.T) {
	ge := gameengine.New()
	ge.ResolveChainStep(ge.Logger(), &card.CardState{Card: notimplemented.RestvineElixirRed{}})
	if got := ge.Value(); got != 0 {
		t.Errorf("Play() = %d, want 0", got)
	}
}

// TestRestvineElixir_NonAttackInRemainingFizzles: non-attack action fails the predicate.
func TestRestvineElixir_NonAttackInRemainingFizzles(t *testing.T) {
	ge := &gameengine.GameEngine{GameState: gameengine.GameStateBuilder().SetCardsRemaining([]*card.CardState{{Card: testutils.GenericAction()}}).Build()}
	ge.ResolveChainStep(ge.Logger(), &card.CardState{Card: notimplemented.RestvineElixirRed{}})
	if got := ge.Value(); got != 0 {
		t.Errorf("Play() = %d, want 0 (non-attack skipped)", got)
	}
}

// TestRestvineElixir_NextAttackGrantsBonusAttack: first attack-action picks up +3 on its
// BonusAttack. Granter returns 0; the +3 attributes to the target.
func TestRestvineElixir_NextAttackGrantsBonusAttack(t *testing.T) {
	target := &card.CardState{Card: testutils.GenericAttack(0, 0)}
	ge := &gameengine.GameEngine{GameState: gameengine.GameStateBuilder().SetCardsRemaining([]*card.CardState{target}).Build()}
	ge.ResolveChainStep(ge.Logger(), &card.CardState{Card: notimplemented.RestvineElixirRed{}})
	if got := ge.Value(); got != 0 {
		t.Errorf("Play() = %d, want 0 (granter returns 0; +N rides on target'ge BonusAttack)", got)
	}
	if target.BonusAttack != 3 {
		t.Errorf("target BonusAttack = %d, want 3", target.BonusAttack)
	}
}
