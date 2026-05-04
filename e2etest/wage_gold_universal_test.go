package e2etest

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/cards"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/heroes"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
)

// Pins Wage Gold's Universal keyword: when Viserai (Runeblade) plays a non-attack action
// followed by Wage Gold, the Universal fold makes Wage Gold count as Runeblade and
// Viserai's hero ability fires. The runechant credits +1 Value at creation and is then
// consumed by Wage Gold's own attack — so we assert the +1 in Value rather than a surviving
// token (it's gone by next-turn start).
func TestWageGold_UniversalTriggersViseraiOnPlay(t *testing.T) {
	d := sim.New(heroes.Viserai{}, nil, fillerDeck())
	hand := []sim.Card{
		cards.HighStrikerBlue{},
		cards.WageGoldRed{},
		testutils.BluePitch{},
		testutils.BluePitch{},
	}
	state := d.EvalOneTurnForTesting(sim.Matchup{IncomingDamage: 0}, nil, hand)
	if state.Value != 8 {
		t.Fatalf("Value = %d, want 8 (Wage Gold 7 + Viserai-via-Universal runechant +1)\nBestLine: %s",
			state.Value, formatBestLine(state.BestLine))
	}
}
