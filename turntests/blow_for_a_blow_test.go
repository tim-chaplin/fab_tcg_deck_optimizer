package turntests

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card/cards"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/deck"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/gameengine"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
)

// Tests that the on-hit 1-damage rider credits +1 on a likely-hit attack. Red lands its 4{p}
// inside the LikelyDamageHits window so the on-hit ping fires.
func TestBlowForABlow_LikelyHitCreditsPing(t *testing.T) {
	d := deck.New(testutils.Hero{Intel: 4}, nil, nil)
	hand := []card.Card{cards.BlowForABlowRed{}, testutils.FakeBlueResource()}
	summary := sim.EvalOneTurnForTesting(d, gameengine.GameStateBuilder().SetIncomingDamage(0).Build(), hand)
	if summary.Value != 4+1 {
		t.Errorf("Value = %d, want 5 (4 likely-hit + 1 ping)", summary.Value)
	}
}
