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

// Tests that the on-hit 1{h} gain credits +1 on a likely-hit attack — Red lands its 4{p}
// inside the LikelyDamageHits window so the on-hit heal fires.
func TestLifeForALife_LikelyHitCreditsHeal(t *testing.T) {
	d := deck.New(testutils.Hero{Intel: 4}, nil, fillerDeck())
	hand := []card.Card{cards.LifeForALifeRed{}, testutils.BluePitch{}}
	summary := sim.EvalOneTurnForTesting(d, gameengine.GameStateBuilder().SetIncomingDamage(0).Build(), hand)
	if summary.Value != 5 {
		t.Errorf("Red: Value = %d, want 5 (4 likely-hit + 1 heal)", summary.Value)
	}
}

// Tests that the heal rider doesn't fire on blockable variants — Yellow 3{p} and Blue 2{p}
// fall outside LikelyDamageHits' {1,4,7} window so the on-hit rider stays silent.
func TestLifeForALife_BlockableSuppressesHeal(t *testing.T) {
	cases := []struct {
		c    card.Card
		want int
	}{
		{cards.LifeForALifeYellow{}, 3},
		{cards.LifeForALifeBlue{}, 2},
	}
	for _, tc := range cases {
		d := deck.New(testutils.Hero{Intel: 4}, nil, fillerDeck())
		hand := []card.Card{tc.c, testutils.BluePitch{}}
		summary := sim.EvalOneTurnForTesting(d, gameengine.GameStateBuilder().SetIncomingDamage(0).Build(), hand)
		if summary.Value != tc.want {
			t.Errorf("%s: Value = %d, want %d (blockable, no heal)", tc.c.Name(), summary.Value, tc.want)
		}
	}
}
