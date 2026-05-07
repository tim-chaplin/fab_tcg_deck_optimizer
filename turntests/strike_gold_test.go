package turntests

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/cards"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/heroes"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
)

// Tests that Strike Gold's on-hit rider lands a Gold token in Items when the attack hits.
func TestStrikeGold_OnHitCreatesGoldToken(t *testing.T) {
	d := sim.New(heroes.Viserai{}, nil, fillerDeck())
	hand := []sim.Card{
		cards.StrikeGoldRed{},
		cards.TitaniumBaubleBlue{},
		cards.TitaniumBaubleBlue{},
		cards.TitaniumBaubleBlue{},
	}
	state := d.EvalOneTurnForTesting(sim.Matchup{IncomingDamage: 0}, sim.TurnState{}, hand)
	if state.Value != 4 {
		t.Fatalf("Value = %d, want 4 (Strike Gold Red 4 power)\nBestLine: %s",
			state.Value, formatBestLine(state.BestLine))
	}
	if got := state.Gold(); got != 1 {
		t.Fatalf("Gold count at start of next turn = %d, want 1 (on-hit token)", got)
	}
}

// Tests that Strike Gold's on-hit rider does not fire when the attack misses LikelyToHit.
func TestStrikeGold_BlockableMissDoesNotCreateGold(t *testing.T) {
	d := sim.New(heroes.Viserai{}, nil, fillerDeck())
	hand := []sim.Card{
		cards.StrikeGoldYellow{},
		cards.TitaniumBaubleBlue{},
		cards.TitaniumBaubleBlue{},
		cards.TitaniumBaubleBlue{},
	}
	state := d.EvalOneTurnForTesting(sim.Matchup{IncomingDamage: 0}, sim.TurnState{}, hand)
	if state.Value != 3 {
		t.Fatalf("Value = %d, want 3 (Strike Gold Yellow 3 power)\nBestLine: %s",
			state.Value, formatBestLine(state.BestLine))
	}
	if got := state.Gold(); got != 0 {
		t.Fatalf("Gold count = %d, want 0 (power-3 attack misses LikelyToHit window)", got)
	}
}

// Tests that a Gold token created on turn 1 carries to turn 2 in the Items list.
func TestStrikeGold_GoldAbilityPlayableNextTurn(t *testing.T) {
	deck := []sim.Card{
		// Turn 1 hand.
		cards.StrikeGoldRed{},
		cards.TitaniumBaubleBlue{},
		cards.TitaniumBaubleBlue{},
		cards.TitaniumBaubleBlue{},
		// Turn 2 hand: red pitches to fund the Gold ability ({2}).
		testutils.RedAttack{},
		testutils.RedAttack{},
		testutils.RedAttack{},
		testutils.RedAttack{},
		// Filler so the dealer can pull a hand.
		testutils.BlueAttack{}, testutils.BlueAttack{},
		testutils.BlueAttack{}, testutils.BlueAttack{},
	}
	d := sim.New(heroes.Viserai{}, nil, deck)
	state := d.EvalOneTurnForTesting(sim.Matchup{IncomingDamage: 0}, sim.TurnState{}, nil)
	if state.Gold() != 1 {
		t.Fatalf("after turn 1: Gold = %d, want 1 (Strike Gold Red on-hit)", state.Gold())
	}
}
