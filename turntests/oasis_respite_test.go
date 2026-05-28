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

// Tests that each printing prevents its full Defense() amount (4/3/2) and that the
// "may gain 1{h}" rider adds +1 iff the current hero opts into LowerHealthWanter.
func TestOasisRespite_PreventsAndLifeRider(t *testing.T) {
	cases := []struct {
		card            card.Card
		wantOff, wantOn int
	}{
		{cards.OasisRespiteRed{}, 4, 5},
		{cards.OasisRespiteYellow{}, 3, 4},
		{cards.OasisRespiteBlue{}, 2, 3},
	}
	for _, tc := range cases {
		hand := []card.Card{tc.card, testutils.FakeBlueResource()}
		dOff := deck.New(testutils.Hero{Intel: 4}, nil, nil)
		sOff := gameengine.GameStateBuilder().
			SetHero(fakeLowHeroOff{}).
			SetIncomingPhysicalDamage(10).
			Build()
		summary := sim.EvalOneTurnForTesting(dOff, sOff, hand)
		if summary.Value != tc.wantOff {
			t.Errorf("%s: hero off Value = %d, want %d", tc.card.Name(), summary.Value, tc.wantOff)
		}

		dOn := deck.New(testutils.Hero{Intel: 4}, nil, nil)
		sOn := gameengine.GameStateBuilder().
			SetHero(fakeLowHeroOn{}).
			SetIncomingPhysicalDamage(10).
			Build()
		summary = sim.EvalOneTurnForTesting(dOn, sOn, hand)
		if summary.Value != tc.wantOn {
			t.Errorf("%s: hero on Value = %d, want %d", tc.card.Name(), summary.Value, tc.wantOn)
		}
	}
}
