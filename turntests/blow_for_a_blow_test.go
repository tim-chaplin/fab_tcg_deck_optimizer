package turntests

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/cards"
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
)

// Tests that the on-hit 1-damage rider credits +1 on a likely-hit attack.
func TestBlowForABlow_LikelyHitCreditsPing(t *testing.T) {
	var s sim.TurnState
	c := cards.BlowForABlowRed{}
	cs := &card.CardState{Card: c}
	sim.ResolveChainStep(&s, s.Logger(), cs)
	testutils.FireOnHitIfLikely(&s, s.Logger(), cs)
	if got := s.Value(); got != 4+1 {
		t.Errorf("Play() = %d, want 5 (4 likely to hit + 1 ping)", got)
	}
}
