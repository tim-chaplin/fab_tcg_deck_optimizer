package turntests

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card/cards"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/deck"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/hero/heroes"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
)

// Tests that a non-attack pitch funding Aether Slash activates the +1 arcane rider.
func TestPitchAttribution_AetherSlashSingleNonAttackPitchFiresRider(t *testing.T) {
	hand := []deck.Card{cards.AetherSlashRed{}, cards.MaleficIncantationBlue{}}
	d := deck.New(heroes.Viserai{}, nil, fillerDeck())
	summary := sim.EvalOneTurnForTesting(d, sim.Matchup{IncomingDamage: 0}, nil, hand)
	if summary.Value != 5 {
		t.Fatalf("Value = %d, want 5 (Aether Slash 4 + rider 1)", summary.Value)
	}
}

// Tests that an attack-typed pitch funding Aether Slash does not activate the rider.
func TestPitchAttribution_AetherSlashAttackPitchDoesNotFireRider(t *testing.T) {
	hand := []deck.Card{cards.AetherSlashRed{}, testutils.YellowAttack{}}
	d := deck.New(heroes.Viserai{}, nil, fillerDeck())
	summary := sim.EvalOneTurnForTesting(d, sim.Matchup{IncomingDamage: 0}, nil, hand)
	if summary.Value != 4 {
		t.Fatalf("Value = %d, want 4 (Aether Slash base power, no rider)", summary.Value)
	}
}

// Tests that Deathly Duet fires both riders when funded by one attack and one non-attack action.
func TestPitchAttribution_DeathlyDuetBothRidersFireFromMixedFunding(t *testing.T) {
	d := deck.New(heroes.Viserai{}, nil, fillerDeck())
	hand := []deck.Card{
		cards.DeathlyDuetRed{},
		cards.AetherSlashRed{},
		cards.MaleficIncantationBlue{},
	}
	summary := sim.EvalOneTurnForTesting(d, sim.Matchup{IncomingDamage: 0}, nil, hand)
	if got := summary.Value; got != 8 {
		t.Fatalf("Value = %d, want 8 (Deathly Duet 4 + attack rider 2 + 2 runechants)", got)
	}
}

// Tests that a single pitch paying for multiple Aether Slashes activates the bonus on each.
func TestPitchAttribution_OneNonAttackPitchFundsMultipleAetherSlashes(t *testing.T) {
	d := deck.New(heroes.Viserai{}, nil, fillerDeck())

	withNonAttack := []deck.Card{
		cards.MauvrionSkiesRed{},
		cards.AetherSlashRed{}, cards.AetherSlashRed{},
		cards.MaleficIncantationBlue{},
	}
	summary := sim.EvalOneTurnForTesting(d, sim.Matchup{IncomingDamage: 0}, nil, withNonAttack)
	if got := summary.Value; got != 15 {
		t.Errorf("non-attack pitch: Value = %d, want 15", got)
	}

	withAttack := []deck.Card{
		cards.MauvrionSkiesRed{},
		cards.AetherSlashRed{}, cards.AetherSlashRed{},
		testutils.BlueAttack{},
	}
	summary = sim.EvalOneTurnForTesting(d, sim.Matchup{IncomingDamage: 0}, nil, withAttack)
	if got := summary.Value; got != 13 {
		t.Errorf("attack pitch: Value = %d, want 13", got)
	}
}

// fillerDeck is a no-op deck body for EvalOneTurnForTesting calls that supply their own
// initialHand. The cards never enter play this turn (the caller's hand is the only thing
// the chain runner sees) but EvalOneTurnForTesting still wants a non-empty deck body so
// the post-turn deal can pull a turn-2 hand without short-circuiting to a zero state.
func fillerDeck() []deck.Card {
	return []deck.Card{
		testutils.BlueAttack{}, testutils.BlueAttack{},
		testutils.BlueAttack{}, testutils.BlueAttack{},
		testutils.BlueAttack{}, testutils.BlueAttack{},
		testutils.BlueAttack{}, testutils.BlueAttack{},
	}
}
