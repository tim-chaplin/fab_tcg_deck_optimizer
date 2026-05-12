package turntests

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/cards"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/gameengine"
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
		sOff := gameengine.NewFromSpec(gameengine.Spec{Hero: stubLowHeroOff{}, IncomingDamage: 10})
		sOff.ResolveChainStep(sOff.Logger(), &card.CardState{Card: tc.card})
		if sOff.Value() != tc.wantOff {
			t.Errorf("%s: hero off Value = %d, want %d", tc.card.Name(), sOff.Value(), tc.wantOff)
		}
		sOn := gameengine.NewFromSpec(gameengine.Spec{Hero: stubLowHeroOn{}, IncomingDamage: 10})
		sOn.ResolveChainStep(sOn.Logger(), &card.CardState{Card: tc.card})
		if sOn.Value() != tc.wantOn {
			t.Errorf("%s: hero on Value = %d, want %d", tc.card.Name(), sOn.Value(), tc.wantOn)
		}
	}
}
