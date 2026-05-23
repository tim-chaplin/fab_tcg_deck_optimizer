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

// Tests that Wage Gold's Universal keyword triggers Viserai's Runeblade hero ability.
func TestWageGold_UniversalTriggersViseraiOnPlay(t *testing.T) {
	d := deck.New(heroes.Viserai, nil, nil)
	hand := []card.Card{
		cards.HighStrikerBlue{},
		cards.WageGoldRed{},
		testutils.FakeBlueResource(),
		testutils.FakeBlueResource(),
	}
	summary := sim.EvalOneTurnForTesting(d, gameengine.GameStateBuilder().SetHero(heroes.Viserai).Build(), hand)
	if summary.Value != 8 {
		t.Fatalf("Value = %d, want 8 (Wage Gold 7 + Viserai-via-Universal runechant +1)\nBestLine: %s",
			summary.Value, formatBestLine(summary.BestLine))
	}
}
