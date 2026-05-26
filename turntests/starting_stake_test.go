package turntests

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card/cards"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/deck"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/gameengine"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/hero/heroes"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
)

// Tests that Starting Stake creates a Gold token via the deck-eval path when no Gold
// is in play and the hand has nothing more profitable to do — solo Starting Stake in
// hand picks the create line over the Held alternative.
func TestStartingStake_CreatesGoldViaAttackTurn(t *testing.T) {
	d := deck.New(heroes.Viserai, nil, nil)
	hand := []card.Card{cards.StartingStakeYellow{}}
	summary := sim.EvalOneTurnForTesting(d, gameengine.GameStateBuilder().SetIncomingDamage(0).Build(), hand)
	if summary.State.GoldCount() != 1 {
		t.Fatalf("Gold = %d, want 1 (Starting Stake creates one)\nBestLine: %s",
			summary.State.GoldCount(), formatBestLine(summary.BestLine))
	}
}

// Tests Starting Stake's "if you control no Gold tokens" gate: with prior Gold in play, Play
// is a no-op and doesn't stack a second token. Drives Play directly to isolate the gate.
func TestStartingStake_NoOpWhenGoldInPlay(t *testing.T) {
	ge := gameengine.New()
	ge.CreateGold(2)
	ge.ResolveAttackStep(ge.Logger(), &card.CardState{Card: cards.StartingStakeYellow{}})
	if ge.GoldCount() != 2 {
		t.Fatalf("Gold = %d, want 2 (already had Gold, Starting Stake is a no-op)", ge.GoldCount())
	}
	_ = testutils.FakeRedAttack()
}
