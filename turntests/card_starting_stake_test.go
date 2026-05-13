package turntests

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/cards"

	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/gameengine"
)

// Tests that Starting Stake creates a Gold token when none are in play.
func TestStartingStake_CreatesGoldWhenNone(t *testing.T) {
	ge := gameengine.New()
	ge.ResolveChainStep(ge.Logger(), &card.CardState{Card: cards.StartingStakeYellow{}})
	if ge.GoldCount() != 1 {
		t.Fatalf("Gold = %d, want 1 (no prior tokens, creates one)", ge.GoldCount())
	}
}

// Tests that Starting Stake does NOT create a Gold token when one is already in play.
func TestStartingStake_NoOpWhenGoldExists(t *testing.T) {
	ge := gameengine.New()
	ge.CreateGold(2)
	ge.ResolveChainStep(ge.Logger(), &card.CardState{Card: cards.StartingStakeYellow{}})
	if ge.GoldCount() != 2 {
		t.Fatalf("Gold = %d, want 2 (already had Gold, no create)", ge.GoldCount())
	}
}
