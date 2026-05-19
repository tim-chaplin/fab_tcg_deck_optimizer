package turntests

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card/cards"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/gameengine"
)

// Tests that Play credits the flat 3-damage prevention.
func TestCalmingBreeze_PreventsFlat3(t *testing.T) {
	ge := &gameengine.GameEngine{GameState: gameengine.GameStateBuilder().SetIncomingDamage(5).Build()}
	pc := &card.CardState{Card: cards.CalmingBreezeRed{}}
	ge.ResolveChainStep(ge.Logger(), pc)
	if ge.Value() != 3 {
		t.Errorf("Value = %d, want 3", ge.Value())
	}
	if ge.RemainingUnblockedDamage() != 2 {
		t.Errorf("RemainingUnblockedDamage = %d, want 2", ge.RemainingUnblockedDamage())
	}
}
