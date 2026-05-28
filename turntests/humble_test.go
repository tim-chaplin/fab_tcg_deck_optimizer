package turntests

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card/cards"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/deck"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/gameengine"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/hero/heroes"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
)

// Tests that each Humble variant scores its printed power when a pitch funds its cost.
func TestHumble_ScoresPrintedPower(t *testing.T) {
	for _, tc := range []struct {
		c    card.Card
		want int
	}{
		{cards.HumbleRed{}, 6},
		{cards.HumbleYellow{}, 5},
		{cards.HumbleBlue{}, 4},
	} {
		d := deck.New(heroes.Viserai, nil, nil)
		hand := []card.Card{tc.c, testutils.FakeBlueResource()}
		state := gameengine.GameStateBuilder().SetIncomingPhysicalDamage(0).Build()
		summary := sim.EvalOneTurnForTesting(d, state, hand)
		if got := summary.Value; got != tc.want {
			t.Errorf("%s: Value = %d, want %d (printed power, BluePitch funds cost 2)",
				tc.c.Name(), got, tc.want)
		}
	}
}
