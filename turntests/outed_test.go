package turntests

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card/cards"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/gameengine"
)

// Tests that Outed doesn't apply the marked-defender bonus when OpponentMarked is false.
func TestOuted_NoMarkUnbuffed(t *testing.T) {
	pc := &card.CardState{Card: cards.OutedRed{}}
	ge := gameengine.New()
	ge.ResolveChainStep(ge.Logger(), pc)
	if pc.BonusAttack != 0 {
		t.Errorf("BonusAttack = %d, want 0 (mark off)", pc.BonusAttack)
	}
}

// Tests that Outed self-buffs +1{p} when the opposing hero is marked at Play time.
func TestOuted_MarkedDefenderAddsOne(t *testing.T) {
	pc := &card.CardState{Card: cards.OutedRed{}}
	ge := &gameengine.GameEngine{GameState: gameengine.GameStateBuilder().SetOpponentMarked(true).Build()}
	ge.ResolveChainStep(ge.Logger(), pc)
	if ge.Value() != 4 {
		t.Errorf("Play() Value = %d, want 4 (3 printed + 1 marked-defender)", ge.Value())
	}
	if pc.BonusAttack != 1 {
		t.Errorf("BonusAttack = %d, want 1", pc.BonusAttack)
	}
}
