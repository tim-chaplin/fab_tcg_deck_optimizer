package turntests

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card/cards"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/deck"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/hero/heroes"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
)

// Tests that Strike Gold's on-hit rider lands a Gold token in Items when the attack hits.
func TestStrikeGold_OnHitCreatesGoldToken(t *testing.T) {
	d := deck.New(heroes.Viserai{}, nil, fillerDeck())
	hand := []deck.Card{
		cards.StrikeGoldRed{},
		testutils.BluePitch{},
		testutils.BluePitch{},
		testutils.BluePitch{},
	}
	summary := sim.EvalOneTurnForTesting(d, sim.Matchup{IncomingDamage: 0}, nil, hand)
	if summary.Value != 4 {
		t.Fatalf("Value = %d, want 4 (Strike Gold Red 4 power)\nBestLine: %s",
			summary.Value, formatBestLine(summary.BestLine))
	}
	if got := summary.State.GoldCount(); got != 1 {
		t.Fatalf("Gold count at start of next turn = %d, want 1 (on-hit token)", got)
	}
}

// Tests that Strike Gold's on-hit rider does not fire when the attack misses LikelyToHit.
func TestStrikeGold_BlockableMissDoesNotCreateGold(t *testing.T) {
	d := deck.New(heroes.Viserai{}, nil, fillerDeck())
	hand := []deck.Card{
		cards.StrikeGoldYellow{},
		testutils.BluePitch{},
		testutils.BluePitch{},
		testutils.BluePitch{},
	}
	summary := sim.EvalOneTurnForTesting(d, sim.Matchup{IncomingDamage: 0}, nil, hand)
	if summary.Value != 3 {
		t.Fatalf("Value = %d, want 3 (Strike Gold Yellow 3 power)\nBestLine: %s",
			summary.Value, formatBestLine(summary.BestLine))
	}
	if got := summary.State.GoldCount(); got != 0 {
		t.Fatalf("Gold count = %d, want 0 (power-3 attack misses LikelyToHit window)", got)
	}
}

// Tests that a Gold token created on turn 1 carries to turn 2 in the Items list.
func TestStrikeGold_GoldTokenPlayableNextTurn(t *testing.T) {
	deckCards := []deck.Card{
		// Turn 1 hand.
		cards.StrikeGoldRed{},
		testutils.BluePitch{},
		testutils.BluePitch{},
		testutils.BluePitch{},
		// Turn 2 hand: red pitches to fund the Gold ability ({2}).
		testutils.RedAttack{},
		testutils.RedAttack{},
		testutils.RedAttack{},
		testutils.RedAttack{},
		// Filler so the dealer can pull a hand.
		testutils.BlueAttack{}, testutils.BlueAttack{},
		testutils.BlueAttack{}, testutils.BlueAttack{},
	}
	d := deck.New(heroes.Viserai{}, nil, deckCards)
	summary := sim.EvalOneTurnForTesting(d, sim.Matchup{IncomingDamage: 0}, nil, nil)
	if summary.State.GoldCount() != 1 {
		t.Fatalf("after turn 1: Gold = %d, want 1 (Strike Gold Red on-hit)", summary.State.GoldCount())
	}
}
