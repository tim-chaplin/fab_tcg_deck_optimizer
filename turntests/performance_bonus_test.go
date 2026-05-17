package turntests

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card/cards"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/deck"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/hero/heroes"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
)

// Tests that Performance Bonus's on-hit Gold-create rider lands a Gold token via the
// chain runner. Solo Blue printing in hand: cost 0, power 1 sits in the LikelyToHit
// window, on-hit fires.
func TestPerformanceBonus_OnHitCreatesGold(t *testing.T) {
	d := deck.New(heroes.Viserai{}, nil, fillerDeck())
	hand := []deck.Card{cards.PerformanceBonusBlue{}}
	gs, extras := sim.EvalOneTurnForTesting(d, sim.Matchup{IncomingDamage: 0}, nil, hand)
	if extras.Value != 1 {
		t.Fatalf("Value = %d, want 1 (PB Blue power 1 hits)\nBestLine: %s",
			extras.Value, formatBestLine(extras.BestLine))
	}
	if gs.GoldCount() != 1 {
		t.Fatalf("Gold = %d, want 1 (on-hit token)", gs.GoldCount())
	}
	_ = testutils.RedAttack{}
}
