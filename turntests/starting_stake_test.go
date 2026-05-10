package turntests

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/cards"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/heroes"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/deck"
)

// Tests that Starting Stake creates a Gold token via the deck-eval path when no Gold
// is in play and the hand has nothing more profitable to do — solo Starting Stake in
// hand picks the create line over the Held alternative.
func TestStartingStake_CreatesGoldViaChain(t *testing.T) {
	d := deck.New(heroes.Viserai{}, nil, fillerDeck())
	hand := []deck.Card{cards.StartingStakeYellow{}}
	state := sim.EvalOneTurnForTesting(d, sim.Matchup{IncomingDamage: 0}, sim.TurnState{}, hand)
	if state.Gold() != 1 {
		t.Fatalf("Gold = %d, want 1 (Starting Stake creates one)\nBestLine: %s",
			state.Gold(), formatBestLine(state.BestLine))
	}
}

// Tests Starting Stake's "if you control no Gold tokens" gate: with prior Gold in play, Play
// is a no-op and doesn't stack a second token. Drives Play directly to isolate the gate.
func TestStartingStake_NoOpWhenGoldInPlay(t *testing.T) {
	s := sim.NewTurnStateFromCards(nil, nil)
	s.CreateGold(2)
	(cards.StartingStakeYellow{}).Play(s, &sim.CardState{Card: cards.StartingStakeYellow{}})
	if s.Gold() != 2 {
		t.Fatalf("Gold = %d, want 2 (already had Gold, Starting Stake is a no-op)", s.Gold())
	}
	_ = testutils.RedAttack{}
}
