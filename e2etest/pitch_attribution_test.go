package e2etest

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/cards"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/heroes"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
)

// Tests that a non-attack pitch funding Aether Slash activates the +1 arcane rider.
func TestPitchAttribution_AetherSlashSingleNonAttackPitchFiresRider(t *testing.T) {
	hand := []sim.Card{cards.AetherSlashRed{}, cards.MaleficIncantationBlue{}}
	d := sim.New(heroes.Viserai{}, nil, fillerDeck())
	state := d.EvalOneTurnForTesting(0, nil, hand)
	if state.PrevTurnValue != 5 {
		t.Fatalf("PrevTurnValue = %d, want 5 (Aether Slash 4 + rider 1)", state.PrevTurnValue)
	}
}

// Tests that an attack-typed pitch funding Aether Slash does not activate the rider.
func TestPitchAttribution_AetherSlashAttackPitchDoesNotFireRider(t *testing.T) {
	hand := []sim.Card{cards.AetherSlashRed{}, testutils.YellowAttack{}}
	d := sim.New(heroes.Viserai{}, nil, fillerDeck())
	state := d.EvalOneTurnForTesting(0, nil, hand)
	if state.PrevTurnValue != 4 {
		t.Fatalf("PrevTurnValue = %d, want 4 (Aether Slash base power, no rider)", state.PrevTurnValue)
	}
}

// Tests that Deathly Duet fires both riders when funded by one attack and one non-attack action.
func TestPitchAttribution_DeathlyDuetBothRidersFireFromMixedFunding(t *testing.T) {
	d := sim.New(heroes.Viserai{}, nil, fillerDeck())
	hand := []sim.Card{
		cards.DeathlyDuetRed{},
		cards.AetherSlashRed{},
		cards.MaleficIncantationBlue{},
	}
	if got := d.EvalOneTurnForTesting(0, nil, hand).PrevTurnValue; got != 8 {
		t.Fatalf("PrevTurnValue = %d, want 8 (Deathly Duet 4 + attack rider 2 + 2 runechants)", got)
	}
}

// Tests that a single pitch paying for multiple Aether Slashes activates the bonus on each.
func TestPitchAttribution_OneNonAttackPitchFundsMultipleAetherSlashes(t *testing.T) {
	d := sim.New(heroes.Viserai{}, nil, fillerDeck())

	withNonAttack := []sim.Card{
		cards.MauvrionSkiesRed{},
		cards.AetherSlashRed{}, cards.AetherSlashRed{},
		cards.MaleficIncantationBlue{},
	}
	if got := d.EvalOneTurnForTesting(0, nil, withNonAttack).PrevTurnValue; got != 15 {
		t.Errorf("non-attack pitch: PrevTurnValue = %d, want 15", got)
	}

	withAttack := []sim.Card{
		cards.MauvrionSkiesRed{},
		cards.AetherSlashRed{}, cards.AetherSlashRed{},
		testutils.BlueAttack{},
	}
	if got := d.EvalOneTurnForTesting(0, nil, withAttack).PrevTurnValue; got != 13 {
		t.Errorf("attack pitch: PrevTurnValue = %d, want 13", got)
	}
}

// fillerDeck is a no-op deck body for EvalOneTurnForTesting calls that supply their own
// initialHand. The cards never enter play this turn (the caller's hand is the only thing
// the chain runner sees) but EvalOneTurnForTesting still wants a non-empty Deck.Cards so
// the post-turn deal can pull a turn-2 hand without short-circuiting to a zero state.
func fillerDeck() []sim.Card {
	return []sim.Card{
		testutils.BlueAttack{}, testutils.BlueAttack{},
		testutils.BlueAttack{}, testutils.BlueAttack{},
		testutils.BlueAttack{}, testutils.BlueAttack{},
		testutils.BlueAttack{}, testutils.BlueAttack{},
	}
}
