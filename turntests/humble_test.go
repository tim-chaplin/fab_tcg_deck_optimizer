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

// Tests that each Humble variant contributes its printed power as damage when a pitch funds
// its cost 2 — the hero-ability suppression rider is unmodelled, so it scores as a vanilla
// attack.
func TestHumble_ScoresPrintedPower(t *testing.T) {
	for _, tc := range []struct {
		c    card.Card
		want int
	}{
		{cards.HumbleRed{}, 6},
		{cards.HumbleYellow{}, 5},
		{cards.HumbleBlue{}, 4},
	} {
		d := deck.New(heroes.Viserai{}, nil, fillerDeck())
		hand := []card.Card{tc.c, testutils.BluePitch{}}
		summary := sim.EvalOneTurnForTesting(d, gameengine.GameStateBuilder().SetIncomingDamage(0).Build(), hand)
		if got := summary.Value; got != tc.want {
			t.Errorf("%s: Value = %d, want %d (printed power, BluePitch funds cost 2)", tc.c.Name(), got, tc.want)
		}
	}
}
