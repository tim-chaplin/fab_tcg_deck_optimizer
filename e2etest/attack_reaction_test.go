package e2etest

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/cards"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/heroes"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/weapons"
)

// Tests that Lunging Press's +1{p} buff lands on a paired attack action card.
func TestAttackReaction_BuffLandsOnTarget(t *testing.T) {
	d := sim.New(heroes.Viserai{}, nil, fillerDeck())
	hand := []sim.Card{
		cards.ArcanicCrackleRed{},
		cards.LungingPressBlue{},
	}
	got := d.EvalOneTurnForTesting(0, nil, hand).PrevTurnValue
	if got != 5 {
		t.Fatalf("PrevTurnValue = %d, want 5 (Arcanic Crackle 3 + Lunging Press +1 buff + 1 arcane)", got)
	}
}

// Tests that an AR with no legal target in hand can't enter the chain.
func TestAttackReaction_NoTargetAtAllNothingHappens(t *testing.T) {
	d := sim.New(heroes.Viserai{}, nil, fillerDeck())
	hand := []sim.Card{cards.LungingPressBlue{}}
	got := d.EvalOneTurnForTesting(0, nil, hand).PrevTurnValue
	if got != 0 {
		t.Fatalf("PrevTurnValue = %d, want 0 (no target, AR can't play)", got)
	}
}

// Tests that an Attack Reaction can't target another Attack Reaction.
func TestAttackReaction_CantTargetAnotherAR(t *testing.T) {
	d := sim.New(heroes.Viserai{}, nil, fillerDeck())
	hand := []sim.Card{cards.LungingPressBlue{}, cards.LungingPressBlue{}}
	got := d.EvalOneTurnForTesting(0, nil, hand).PrevTurnValue
	if got != 0 {
		t.Fatalf("PrevTurnValue = %d, want 0 (ARs can't target each other)", got)
	}
}

// Tests that Attack Reactions can't play before the attack they target (and thereby trigger
// Viserai's hero ability).
func TestAttackReaction_DoesNotTriggerViseraiAsNonAttackAction(t *testing.T) {
	d := sim.New(heroes.Viserai{}, nil, fillerDeck())
	hand := []sim.Card{cards.LungingPressBlue{}, cards.HocusPocusRed{}}
	got := d.EvalOneTurnForTesting(0, nil, hand).PrevTurnValue
	if got != 5 {
		t.Fatalf("PrevTurnValue = %d, want 5 (HP 3 + LP +1 buff + HP runechant 1; Viserai must not fire)", got)
	}
}

// Tests that an AR can only target an attack that is actually played, not one that's been
// pitched.
func TestAttackReaction_PitchedAttackIsNotATarget(t *testing.T) {
	d := sim.New(heroes.Viserai{}, nil, fillerDeck())
	hand := []sim.Card{
		cards.LungingPressBlue{},
		testutils.GenericAttack(1, 0),
	}
	got := d.EvalOneTurnForTesting(0, nil, hand).PrevTurnValue
	if got != 0 {
		t.Fatalf("PrevTurnValue = %d, want 0 (attack pitched ⇒ AR has nothing to buff)", got)
	}
}

// Tests that Thrust's +3{p} buff lands on a swinging Sword weapon.
func TestAttackReaction_ThrustBuffsSwingingSwordWeapon(t *testing.T) {
	d := sim.New(heroes.Viserai{}, []sim.Weapon{weapons.NebulaBlade{}}, fillerDeck())
	hand := []sim.Card{
		cards.ThrustRed{},
		cards.ToughenUpBlue{},
	}
	got := d.EvalOneTurnForTesting(0, nil, hand).PrevTurnValue
	if got != 5 {
		t.Fatalf("PrevTurnValue = %d, want 5 (Nebula Blade 1 + Thrust +3 buff + runechant 1)", got)
	}
}
