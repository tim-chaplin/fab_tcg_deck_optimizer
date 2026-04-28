package cards

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

// TestBlowForABlow_LikelyHitCreditsPing: Red (4) lands in the likely-to-hit set; the on-hit
// 1-damage rider credits +1.
func TestBlowForABlow_LikelyHitCreditsPing(t *testing.T) {
	var s sim.TurnState
	c := BlowForABlowRed{}
	c.Play(&s, &sim.CardState{Card: c})
	if got := s.Value; got != 4+1 {
		t.Errorf("Play() = %d, want 5 (4 likely to hit + 1 ping)", got)
	}
}
