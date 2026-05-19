package turntests

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card/cards"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/gameengine"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
)

// Tests that Pursue to the Pits of Despair's Play registers an OnHit handler that marks
// the opposing hero when LikelyToHit fires.
func TestPursueToThePitsOfDespair_OnHitMarksOpponent(t *testing.T) {
	pc := &card.CardState{Card: cards.PursueToThePitsOfDespairRed{}}
	ge := gameengine.New()
	ge.ResolveChainStep(ge.Logger(), pc)
	if len(pc.OnHit) != 1 {
		t.Fatalf("OnHit handlers = %d, want 1", len(pc.OnHit))
	}
	// Printed 5{p} doesn't fit the 1/4/7 LikelyDamageHits window; bump to 7 to drain.
	pc.BonusAttack = 2
	testutils.FireOnHitIfLikely(ge, ge.Logger(), pc)
	if !ge.OpponentMarked() {
		t.Errorf("OpponentMarked = false after OnHit fires, want true")
	}
}
