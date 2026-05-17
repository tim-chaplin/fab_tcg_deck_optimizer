package turntests

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/v2/card/cards"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/gameengine"
)

// Tests that the on-hit 1-damage rider credits +1 on a likely-hit attack.
func TestBlowForABlow_LikelyHitCreditsPing(t *testing.T) {
	ge := gameengine.New()
	c := cards.BlowForABlowRed{}
	cs := &card.CardState{Card: c}
	ge.ResolveChainStep(ge.Logger(), cs)
	testutils.FireOnHitIfLikely(ge, ge.Logger(), cs)
	if got := ge.Value(); got != 4+1 {
		t.Errorf("Play() = %d, want 5 (4 likely to hit + 1 ping)", got)
	}
}
