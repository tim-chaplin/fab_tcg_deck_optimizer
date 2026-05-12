package turntests

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/cards"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/heroes"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/deck"
)

// Tests that Wage Gold's on-hit rider creates a Gold token when the attack hits.
func TestWageGold_OnHitCreatesGoldToken(t *testing.T) {
	d := deck.New(heroes.Viserai{}, nil, fillerDeck())
	hand := []deck.Card{
		cards.WageGoldRed{},
		testutils.BluePitch{},
		testutils.BluePitch{},
		testutils.BluePitch{},
	}
	state := sim.EvalOneTurnForTesting(d, sim.Matchup{IncomingDamage: 0}, sim.TurnState{}, hand)
	if state.Value != 7 {
		t.Fatalf("Value = %d, want 7 (Wage Gold Red 7 power)\nBestLine: %s",
			state.Value, formatBestLine(state.BestLine))
	}
	if got := state.GoldCount(); got != 1 {
		t.Fatalf("Gold count at start of next turn = %d, want 1 (on-hit token)", got)
	}
}

// Tests that Wage Gold's on-hit rider skips Gold when the attack misses LikelyToHit.
func TestWageGold_BlockableMissDoesNotCreateGold(t *testing.T) {
	d := deck.New(heroes.Viserai{}, nil, fillerDeck())
	hand := []deck.Card{
		cards.WageGoldBlue{},
		testutils.BluePitch{},
		testutils.BluePitch{},
		testutils.BluePitch{},
	}
	state := sim.EvalOneTurnForTesting(d, sim.Matchup{IncomingDamage: 0}, sim.TurnState{}, hand)
	if state.Value != 5 {
		t.Fatalf("Value = %d, want 5 (Wage Gold Blue 5 power)\nBestLine: %s",
			state.Value, formatBestLine(state.BestLine))
	}
	if got := state.GoldCount(); got != 0 {
		t.Fatalf("Gold count = %d, want 0 (power-5 attack misses LikelyToHit window)", got)
	}
}
