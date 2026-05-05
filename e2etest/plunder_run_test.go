package e2etest

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/cards"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/heroes"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
)

// Pins Plunder Run's "next time an attack action card you control hits" rider waiting
// across a missed attack and firing on the next one that lands. Hand: pitch-only Blue,
// Plunder Run R, Runerager Swarm R, Critical Strike Y. Pitched Blue funds CS's 1 cost;
// Plunder Run sets NonAttackActionPlayed so Runerager triggers Viserai's runechant (+1);
// Runerager misses the LikelyToHit window at power 3 so the trigger stays queued; CS at
// 4 lands the queue → DrawOne fires.
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
