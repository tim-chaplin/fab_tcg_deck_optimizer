package e2etest

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/cards"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/heroes"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
)

// Tests that Starting Stake creates a Gold token via the deck-eval path when no Gold
// is in play and the hand has nothing more profitable to do — solo Starting Stake in
// hand picks the create line over the Held alternative.
func TestStartingStake_CreatesGoldViaChain(t *testing.T) {
	d := sim.New(heroes.Viserai{}, nil, fillerDeck())
	hand := []sim.Card{cards.StartingStakeYellow{}}
	state := d.EvalOneTurnForTesting(sim.Matchup{IncomingDamage: 0}, nil, hand)
	if state.Gold() != 1 {
		t.Fatalf("Gold = %d, want 1 (Starting Stake creates one)\nBestLine: %s",
			state.Gold(), formatBestLine(state.BestLine))
	}
}

// Tests the no-op branch: with prior Gold already in play, Starting Stake doesn't
// stack a second token. Asserts via Play() directly because the optimizer might
// also fund the Gold ability spend in the same chain (which would also leave Gold
// at 0) — we want to isolate the "if you control no Gold tokens, create one" gate.
func TestStartingStake_NoOpWhenGoldInPlay(t *testing.T) {
	s := sim.NewTurnState(nil, nil)
	s.CreateGold(2)
	(cards.StartingStakeYellow{}).Play(s, &sim.CardState{Card: cards.StartingStakeYellow{}})
	if s.Gold() != 2 {
		t.Fatalf("Gold = %d, want 2 (already had Gold, Starting Stake is a no-op)", s.Gold())
	}
	_ = testutils.RedAttack{}
}
