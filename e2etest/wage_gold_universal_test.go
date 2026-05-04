package e2etest

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/cards"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/heroes"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
)

// Tests that Wage Gold's Universal keyword triggers Viserai's Runeblade hero ability.
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
