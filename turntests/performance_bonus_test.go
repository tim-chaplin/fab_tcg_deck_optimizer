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

// Tests that Performance Bonus's on-hit Gold-create rider lands a Gold token via the
// attack-turn runner. Solo Blue printing in hand: cost 0, power 1 sits in the LikelyToHit
// window, on-hit fires.
func TestPerformanceBonus_OnHitCreatesGold(t *testing.T) {
	d := deck.New(heroes.Viserai, nil, nil)
	hand := []card.Card{cards.PerformanceBonusBlue{}}
	summary := sim.EvalOneTurnForTesting(d, gameengine.GameStateBuilder().SetIncomingPhysicalDamage(0).Build(), hand)
	if summary.Value != 1 {
		t.Fatalf("Value = %d, want 1 (PB Blue power 1 hits)\nBestLine: %s",
			summary.Value, formatBestLine(summary.BestLine))
	}
	if summary.State.GoldCount() != 1 {
		t.Fatalf("Gold = %d, want 1 (on-hit token)", summary.State.GoldCount())
	}
	_ = testutils.FakeRedAttack()
}
