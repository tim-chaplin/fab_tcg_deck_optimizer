package turntests

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/v2/card/cards"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/gameengine"
)

// Tests that Pursue to the Edge of Oblivion's Play registers an OnHit handler that marks
// the opposing hero when LikelyToHit fires.
func TestPursueToTheEdgeOfOblivion_OnHitMarksOpponent(t *testing.T) {
	self := &card.CardState{Card: cards.PursueToTheEdgeOfOblivionRed{}}
	ge := gameengine.New()
	ge.ResolveChainStep(ge.Logger(), self)
	if len(self.OnHit) != 1 {
		t.Fatalf("OnHit handlers = %d, want 1", len(self.OnHit))
	}
	if ge.OpponentMarked() {
		t.Errorf("OpponentMarked = true before OnHit fires, want false")
	}
	testutils.FireOnHitIfLikely(ge, ge.Logger(), self)
	if !ge.OpponentMarked() {
		t.Errorf("OpponentMarked = false after OnHit fires, want true")
	}
}
