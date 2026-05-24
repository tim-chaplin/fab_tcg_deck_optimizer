package turntests

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card/cards"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/deck"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/gameengine"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
)

// Tests that Rise Above defending (DR phase) does not ghost-duplicate a card.
func TestBest_RiseAboveAsDR_NoGhostCard(t *testing.T) {
	h := []card.Card{
		cards.RiseAboveRed{},
		testutils.FakeBlueResource(),
		testutils.FakeRedAttack(),
		testutils.FakeYellowResource(),
	}
	d := deck.New(testutils.Hero{Intel: 4}, nil, nil)
	summary := sim.EvalOneTurnForTesting(
		d,
		gameengine.GameStateBuilder().SetIncomingDamage(3).Build(),
		h,
	)

	post := summary.State.Deck().Size() + len(summary.State.Graveyard()) + len(summary.State.Hand()) + len(summary.State.Banished())
	if summary.State.Arsenal() != nil {
		post++
	}
	if post != 4 {
		t.Errorf("total card count = %d, want 4. deck=%d grav=%v hand=%v banished=%v arsenal=%v\n  BestLine: %s",
			post, summary.State.Deck().Size(), summary.State.Graveyard(), summary.State.Hand(), summary.State.Banished(), summary.State.Arsenal(),
			sim.FormatBestLine(summary.BestLine))
	}
}
