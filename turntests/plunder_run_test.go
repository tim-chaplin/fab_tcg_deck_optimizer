package turntests

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/cards"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/heroes"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
)

// Tests that Plunder Run's "next time an attack action card hits" trigger waits across a
// missed attack and fires on the next attack action card that lands.
func TestPlunderRun_TriggerWaitsAcrossMissAndFiresOnHit(t *testing.T) {
	d := sim.New(heroes.Viserai{}, nil, fillerDeck())
	hand := []sim.Card{
		testutils.BluePitch{},
		cards.PlunderRunRed{},
		cards.RuneragerSwarmRed{},
		cards.CriticalStrikeYellow{},
	}
	state := d.EvalOneTurnForTesting(sim.Matchup{IncomingDamage: 0}, sim.TurnState{}, hand)
	if state.Value != 8 {
		t.Fatalf("Value = %d, want 8 (Runerager 3 + Viserai runechant 1 + CS 4)\nBestLine: %s",
			state.Value, formatBestLine(state.BestLine))
	}
}
