package turntests

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/v2/card/cards"

	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/gameengine"
)

// Tests that each variant credits printed_power - 2 to gs.Value().
func TestFreewheelingRenegades_AlwaysDebuffedByTwo(t *testing.T) {
	cases := []struct {
		c    card.Card
		want int
	}{
		{cards.FreewheelingRenegadesRed{}, 4},
		{cards.FreewheelingRenegadesYellow{}, 3},
		{cards.FreewheelingRenegadesBlue{}, 2},
	}
	for _, tc := range cases {
		ge := gameengine.New()
		ge.ResolveChainStep(ge.Logger(), &card.CardState{Card: tc.c})
		if got := ge.Value(); got != tc.want {
			t.Errorf("%s: Play() = %d, want %d (printed - 2)", tc.c.Name(), got, tc.want)
		}
	}
}
