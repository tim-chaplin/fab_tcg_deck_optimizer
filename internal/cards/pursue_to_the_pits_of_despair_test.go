package cards

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
)

// Tests that Pursue to the Pits of Despair's Play registers an OnHit handler that marks
// the opposing hero when LikelyToHit fires.
func TestPursueToThePitsOfDespair_OnHitMarksOpponent(t *testing.T) {
	self := &sim.CardState{Card: PursueToThePitsOfDespairRed{}}
	s := sim.TurnState{}
	(PursueToThePitsOfDespairRed{}).Play(&s, self)
	if len(self.OnHit) != 1 {
		t.Fatalf("OnHit handlers = %d, want 1", len(self.OnHit))
	}
	// Printed 5{p} doesn't fit the 1/4/7 LikelyDamageHits window; bump to 7 to drain.
	self.BonusAttack = 2
	testutils.FireOnHitIfLikely(&s, self)
	if !s.OpponentMarked {
		t.Errorf("OpponentMarked = false after OnHit fires, want true")
	}
}
