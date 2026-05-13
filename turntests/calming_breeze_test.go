package turntests

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/cards"

	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/gameengine"
)

// Tests that Play credits the flat 3-damage prevention.
func TestCalmingBreeze_PreventsFlat3(t *testing.T) {
	s := &gameengine.GameEngine{GameState: gameengine.GameStateBuilder().SetIncomingDamage(5).Build()}
	self := &card.CardState{Card: cards.CalmingBreezeRed{}}
	s.ResolveChainStep(s.Logger(), self)
	if s.Value() != 3 {
		t.Errorf("Value = %d, want 3", s.Value())
	}
	if s.IncomingDamage() != 2 {
		t.Errorf("IncomingDamage = %d, want 2", s.IncomingDamage())
	}
}
