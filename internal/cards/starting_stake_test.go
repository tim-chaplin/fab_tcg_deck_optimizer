package cards

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
)

// Tests that Starting Stake creates a Gold token when none are in play.
func TestStartingStake_CreatesGoldWhenNone(t *testing.T) {
	s := sim.NewTurnStateFromCards(nil, nil)
	sim.ResolveChainStep(s, s.Logger(), &card.CardState{Card: StartingStakeYellow{}})
	if s.Gold() != 1 {
		t.Fatalf("Gold = %d, want 1 (no prior tokens, creates one)", s.Gold())
	}
}

// Tests that Starting Stake does NOT create a Gold token when one is already in play.
func TestStartingStake_NoOpWhenGoldExists(t *testing.T) {
	s := sim.NewTurnStateFromCards(nil, nil)
	s.CreateGold(2)
	sim.ResolveChainStep(s, s.Logger(), &card.CardState{Card: StartingStakeYellow{}})
	if s.Gold() != 2 {
		t.Fatalf("Gold = %d, want 2 (already had Gold, no create)", s.Gold())
	}
}
