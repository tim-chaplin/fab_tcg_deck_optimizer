package turntests

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/cards"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/heroes"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/deck"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/gameengine"
)

// Tests that Starting Stake creates a Gold token via the deck-eval path when no Gold
// is in play and the hand has nothing more profitable to do — solo Starting Stake in
// hand picks the create line over the Held alternative.
func TestStartingStake_CreatesGoldViaChain(t *testing.T) {
	d := deck.New(heroes.Viserai{}, nil, fillerDeck())
	hand := []deck.Card{cards.StartingStakeYellow{}}
	state := sim.EvalOneTurnForTesting(d, sim.Matchup{IncomingDamage: 0}, gameengine.Spec{}, hand)
	if state.GoldCount() != 1 {
		t.Fatalf("Gold = %d, want 1 (Starting Stake creates one)\nBestLine: %s",
			state.GoldCount(), formatBestLine(state.BestLine))
	}
}

// Tests Starting Stake's "if you control no Gold tokens" gate: with prior Gold in play, Play
// is a no-op and doesn't stack a second token. Drives Play directly to isolate the gate.
func TestStartingStake_NoOpWhenGoldInPlay(t *testing.T) {
	s := gameengine.NewFromCards(nil, nil)
	s.CreateGold(2)
	s.ResolveChainStep(s.Logger(), &card.CardState{Card: cards.StartingStakeYellow{}})
	if s.GoldCount() != 2 {
		t.Fatalf("Gold = %d, want 2 (already had Gold, Starting Stake is a no-op)", s.GoldCount())
	}
	_ = testutils.RedAttack{}
}
