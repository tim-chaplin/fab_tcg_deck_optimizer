package turntests

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/cards"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/heroes"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/deck"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/gameengine"
)

// Tests that Performance Bonus's on-hit Gold-create rider lands a Gold token via the
// chain runner. Solo Blue printing in hand: cost 0, power 1 sits in the LikelyToHit
// window, on-hit fires.
func TestPerformanceBonus_OnHitCreatesGold(t *testing.T) {
	d := deck.New(heroes.Viserai{}, nil, fillerDeck())
	hand := []deck.Card{cards.PerformanceBonusBlue{}}
	state := sim.EvalOneTurnForTesting(d, sim.Matchup{IncomingDamage: 0}, gameengine.Spec{}, hand)
	if state.Value != 1 {
		t.Fatalf("Value = %d, want 1 (PB Blue power 1 hits)\nBestLine: %s",
			state.Value, formatBestLine(state.BestLine))
	}
	if state.GoldCount() != 1 {
		t.Fatalf("Gold = %d, want 1 (on-hit token)", state.GoldCount())
	}
	_ = testutils.RedAttack{}
}
