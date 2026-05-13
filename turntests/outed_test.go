package turntests

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/cards"

	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/gameengine"
)

// Tests that Outed doesn't apply the marked-defender bonus when OpponentMarked is false.
func TestOuted_NoMarkUnbuffed(t *testing.T) {
	self := &card.CardState{Card: cards.OutedRed{}}
	s := gameengine.New()
	s.ResolveChainStep(s.Logger(), self)
	if self.BonusAttack != 0 {
		t.Errorf("BonusAttack = %d, want 0 (mark off)", self.BonusAttack)
	}
}

// Tests that Outed self-buffs +1{p} when the opposing hero is marked at Play time.
func TestOuted_MarkedDefenderAddsOne(t *testing.T) {
	self := &card.CardState{Card: cards.OutedRed{}}
	s := &gameengine.GameEngine{GameState: gameengine.GameStateBuilder().SetOpponentMarked(true).Build()}
	s.ResolveChainStep(s.Logger(), self)
	if s.Value() != 4 {
		t.Errorf("Play() Value = %d, want 4 (3 printed + 1 marked-defender)", s.Value())
	}
	if self.BonusAttack != 1 {
		t.Errorf("BonusAttack = %d, want 1", self.BonusAttack)
	}
}
