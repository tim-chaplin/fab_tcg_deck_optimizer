package generic

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

// TestBlowForABlow_LikelyHitCreditsPing: Red (4) lands in the likely-to-hit set; the on-hit
// 1-damage rider credits +1.
func TestBlowForABlow_LikelyHitCreditsPing(t *testing.T) {
	var s card.TurnState
	c := BlowForABlowRed{}
	if got := c.Play(&s, &card.CardState{Card: c}); got != 4+1 {
		t.Errorf("Play() = %d, want 5 (4 likely to hit + 1 ping)", got)
	}
}
