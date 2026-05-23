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

// Tests that each variant credits printed_power - 2 to summary.Value.
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
		d := deck.New(testutils.Hero{Intel: 4}, nil, nil)
		hand := []card.Card{tc.c, testutils.BluePitch{}}
		summary := sim.EvalOneTurnForTesting(d, gameengine.GameStateBuilder().SetIncomingDamage(0).Build(), hand)
		if summary.Value != tc.want {
			t.Errorf("%s: Value = %d, want %d (printed - 2)", tc.c.Name(), summary.Value, tc.want)
		}
	}
}
