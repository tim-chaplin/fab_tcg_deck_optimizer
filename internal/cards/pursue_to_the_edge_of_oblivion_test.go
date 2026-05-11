package cards

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
)

// Tests that Pursue to the Edge of Oblivion's Play registers an OnHit handler that marks
// the opposing hero when LikelyToHit fires.
func TestPursueToTheEdgeOfOblivion_OnHitMarksOpponent(t *testing.T) {
	self := &sim.CardState{Card: PursueToTheEdgeOfOblivionRed{}}
	s := sim.TurnState{}
	sim.ResolveChainStep(&s, s.Logger(), self)
	if len(self.OnHit) != 1 {
		t.Fatalf("OnHit handlers = %d, want 1", len(self.OnHit))
	}
	if s.OpponentMarked {
		t.Errorf("OpponentMarked = true before OnHit fires, want false")
	}
	testutils.FireOnHitIfLikely(&s, s.Logger(), self)
	if !s.OpponentMarked {
		t.Errorf("OpponentMarked = false after OnHit fires, want true")
	}
}
