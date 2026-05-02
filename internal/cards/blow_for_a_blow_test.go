package cards

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
)

// Tests that the on-hit 1-damage rider credits +1 on a likely-hit attack.
func TestBlowForABlow_LikelyHitCreditsPing(t *testing.T) {
	var s sim.TurnState
	c := BlowForABlowRed{}
	cs := &sim.CardState{Card: c}
	c.Play(&s, cs)
	testutils.FireOnHitIfLikely(&s, cs)
	if got := s.Value; got != 4+1 {
		t.Errorf("Play() = %d, want 5 (4 likely to hit + 1 ping)", got)
	}
}
