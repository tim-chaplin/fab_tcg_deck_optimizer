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

// Tests that Wage Gold's on-hit rider creates a Gold token when the attack hits.
func TestWageGold_OnHitCreatesGoldToken(t *testing.T) {
	d := deck.New(heroes.Viserai{}, nil, fillerDeck())
	hand := []card.Card{
		cards.WageGoldRed{},
		testutils.BluePitch{},
		testutils.BluePitch{},
		testutils.BluePitch{},
	}
	summary := sim.EvalOneTurnForTesting(d, gameengine.GameStateBuilder().SetIncomingDamage(0).Build(), hand)
	if summary.Value != 7 {
		t.Fatalf("Value = %d, want 7 (Wage Gold Red 7 power)\nBestLine: %s",
			summary.Value, formatBestLine(summary.BestLine))
	}
	if got := summary.State.GoldCount(); got != 1 {
		t.Fatalf("Gold count at start of next turn = %d, want 1 (on-hit token)", got)
	}
}

// Tests that Wage Gold's on-hit rider skips Gold when the attack misses LikelyToHit.
func TestWageGold_BlockableMissDoesNotCreateGold(t *testing.T) {
	d := deck.New(heroes.Viserai{}, nil, fillerDeck())
	hand := []card.Card{
		cards.WageGoldBlue{},
		testutils.BluePitch{},
		testutils.BluePitch{},
		testutils.BluePitch{},
	}
	summary := sim.EvalOneTurnForTesting(d, gameengine.GameStateBuilder().SetIncomingDamage(0).Build(), hand)
	if summary.Value != 5 {
		t.Fatalf("Value = %d, want 5 (Wage Gold Blue 5 power)\nBestLine: %s",
			summary.Value, formatBestLine(summary.BestLine))
	}
	if got := summary.State.GoldCount(); got != 0 {
		t.Fatalf("Gold count = %d, want 0 (power-5 attack misses LikelyToHit window)", got)
	}
}
