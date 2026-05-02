package e2etest

// End-to-end tests for Attack Reaction partition validation and buff flow-through.

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
	// Arcanic Crackle (0 cost, 3 power, Go again) + Lunging Press (AR, 0 cost). Optimal
	// order [Lunging Press, Arcanic Crackle] buffs the attack to 4{p}.
	hand := []sim.Card{
		cards.ArcanicCrackleRed{},
		cards.LungingPressBlue{},
	}
	got := d.EvalOneTurnForTesting(0, nil, hand).PrevTurnValue
	if got != 5 {
		t.Fatalf("PrevTurnValue = %d, want 5 (Arcanic Crackle 3 + Lunging Press +1 buff + 1 arcane)", got)
	}
}

// Tests that an AR with no legal target in hand can't enter the chain (validator rejects
// every attack-role assignment, AR stays Held, PrevTurnValue is zero).
func TestAttackReaction_NoTargetAtAllNothingHappens(t *testing.T) {
	d := sim.New(heroes.Viserai{}, nil, fillerDeck())
	hand := []sim.Card{cards.LungingPressBlue{}}
	got := d.EvalOneTurnForTesting(0, nil, hand).PrevTurnValue
	if got != 0 {
		t.Fatalf("PrevTurnValue = %d, want 0 (no target, AR can't play)", got)
	}
}

// Tests that an AR can only target an attack that is actually played, not one that's been
// pitched. With Lunging Press (AR, pitch 3) + a cost-1 attack action, the only feasible
// chain pitches LP and plays the attack alone — LP isn't in the chain to buff anything,
// and the partition that puts LP in Attack has no target so the validator rejects it.
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

// Tests that Thrust's +3{p} buff lands on a swinging Sword weapon (exercises the wmask
// path through partitionHasValidARTargets).
func TestAttackReaction_ThrustBuffsSwingingSwordWeapon(t *testing.T) {
	d := sim.New(heroes.Viserai{}, []sim.Weapon{weapons.NebulaBlade{}}, fillerDeck())
	// Thrust [R] (cost 1) + ToughenUp [B] (pitches 3) funds Nebula Blade's 2-cost swing +
	// Thrust's 1. Optimal chain: pitch Toughen Up, play Thrust, swing Nebula Blade.
	hand := []sim.Card{
		cards.ThrustRed{},
		cards.ToughenUpBlue{},
	}
	got := d.EvalOneTurnForTesting(0, nil, hand).PrevTurnValue
	if got != 5 {
		t.Fatalf("PrevTurnValue = %d, want 5 (Nebula Blade 1 + Thrust +3 buff + runechant 1)", got)
	}
}
