package turntests

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/cards"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/heroes"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
)

// Tests that Performance Bonus's on-hit Gold-create rider lands a Gold token via the
// chain runner. Solo Blue printing in hand: cost 0, power 1 sits in the LikelyToHit
// window, on-hit fires.
func TestPerformanceBonus_OnHitCreatesGold(t *testing.T) {
	d := sim.New(heroes.Viserai{}, nil, fillerDeck())
	hand := []sim.Card{cards.PerformanceBonusBlue{}}
	summary := d.EvalOneTurnForTesting(sim.Matchup{IncomingDamage: 0}, sim.TurnState{}, hand)
	state := summary.State
	if summary.Value != 1 {
		t.Fatalf("Value = %d, want 1 (PB Blue power 1 hits)\nBestLine: %s",
			summary.Value, formatBestLine(summary.BestLine))
	}
	if state.Gold() != 1 {
		t.Fatalf("Gold = %d, want 1 (on-hit token)", state.Gold())
	}
	_ = testutils.RedAttack{}
}
