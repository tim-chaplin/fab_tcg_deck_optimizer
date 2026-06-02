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

// Tests that playing a real instant card turns Volatile Fluxor's +3 rider on: a 0-power instant
// followed by Red Fluxor swings for 3.
func TestVolatileFluxor_InstantPlayedTurnsOnRider(t *testing.T) {
	hero := testutils.Hero{Intel: 4}
	instant := testutils.FakeRedInstant().WithName("Instant")
	state := gameengine.GameStateBuilder().SetHero(hero).Build()
	d := deck.New(hero, nil, nil)
	summary := sim.EvalOneTurnForTesting(d, state, []card.Card{instant, cards.VolatileFluxorRed{}})

	if got, want := summary.Value, 3; got != want {
		t.Errorf("Value = %d, want %d (Fluxor +3 from the played instant)\nBestLine: %s",
			got, want, sim.FormatBestLine(summary.BestLine))
	}
}

// Tests that activating Amulet of Lightning's instant ability does NOT count as playing an
// instant: with the Amulet crackable (Lightning fused) and Fluxor in hand, cracking the Amulet
// must not turn on Fluxor's +3 rider, so the turn is worth 0.
func TestVolatileFluxor_CrackingAmuletIsNotPlayingAnInstant(t *testing.T) {
	hero := testutils.Hero{Intel: 4}
	fuse := testutils.FakeRedAction().
		WithName("Fuse").
		WithGoAgain().
		WithPlay(func(ge card.GameEngine, _ card.Logger, _ *card.CardState) { ge.SetLightningFused(true) })
	state := gameengine.GameStateBuilder().
		SetHero(hero).
		CreateItemFromCard(cards.AmuletOfLightningBlue{}).
		Build()
	d := deck.New(hero, nil, nil)
	summary := sim.EvalOneTurnForTesting(d, state, []card.Card{fuse, cards.VolatileFluxorRed{}})

	if got := summary.Value; got != 0 {
		t.Errorf("Value = %d, want 0 — cracking the Amulet (an instant ability) must not satisfy "+
			"Fluxor's 'played an instant' rider\nBestLine: %s", got, sim.FormatBestLine(summary.BestLine))
	}
}
