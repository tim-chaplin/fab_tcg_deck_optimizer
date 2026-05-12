package turntests

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/cards"
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
)

// Tests that Pursue to the Edge of Oblivion's Play registers an OnHit handler that marks
// the opposing hero when LikelyToHit fires.
func TestPursueToTheEdgeOfOblivion_OnHitMarksOpponent(t *testing.T) {
	self := &card.CardState{Card: cards.PursueToTheEdgeOfOblivionRed{}}
	s := sim.TurnState{}
	sim.ResolveChainStep(&s, s.Logger(), self)
	if len(self.OnHit) != 1 {
		t.Fatalf("OnHit handlers = %d, want 1", len(self.OnHit))
	}
	if s.OpponentMarked() {
		t.Errorf("OpponentMarked = true before OnHit fires, want false")
	}
	testutils.FireOnHitIfLikely(&s, s.Logger(), self)
	if !s.OpponentMarked() {
		t.Errorf("OpponentMarked = false after OnHit fires, want true")
	}
}
