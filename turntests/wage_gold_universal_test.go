package turntests

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/cards"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/deck"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/hero"
)

// Tests that Wage Gold's Universal keyword triggers Viserai's Runeblade hero ability.
func TestWageGold_UniversalTriggersViseraiOnPlay(t *testing.T) {
	d := deck.New(hero.Viserai{}, nil, fillerDeck())
	hand := []deck.Card{
		cards.HighStrikerBlue{},
		cards.WageGoldRed{},
		testutils.BluePitch{},
		testutils.BluePitch{},
	}
	state := sim.EvalOneTurnForTesting(d, sim.Matchup{IncomingDamage: 0}, nil, hand)
	if state.Value != 8 {
		t.Fatalf("Value = %d, want 8 (Wage Gold 7 + Viserai-via-Universal runechant +1)\nBestLine: %s",
			state.Value, formatBestLine(state.BestLine))
	}
}
