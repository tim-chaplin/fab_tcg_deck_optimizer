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

// Tests that Plunder Run's "next time an attack action card hits" trigger waits across a
// missed attack and fires on the next attack action card that lands.
func TestPlunderRun_TriggerWaitsAcrossMissAndFiresOnHit(t *testing.T) {
	d := deck.New(heroes.Viserai, nil, nil)
	hand := []card.Card{
		testutils.BluePitch{},
		cards.PlunderRunRed{},
		cards.RuneragerSwarmRed{},
		cards.CriticalStrikeYellow{},
	}
	summary := sim.EvalOneTurnForTesting(d, gameengine.GameStateBuilder().SetHero(heroes.Viserai).Build(), hand)
	if summary.Value != 8 {
		t.Fatalf("Value = %d, want 8 (Runerager 3 + Viserai runechant 1 + CS 4)\nBestLine: %s",
			summary.Value, formatBestLine(summary.BestLine))
	}
}
