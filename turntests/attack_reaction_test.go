package turntests

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/cards"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/weapons"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/deck"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/hero"
)

// Tests that Lunging Press's +1{p} buff lands on a paired attack action card.
func TestAttackReaction_BuffLandsOnTarget(t *testing.T) {
	d := deck.New(hero.Viserai{}, nil, fillerDeck())
	hand := []deck.Card{
		cards.ArcanicCrackleRed{},
		cards.LungingPressBlue{},
	}
	got := sim.EvalOneTurnForTesting(d, sim.Matchup{IncomingDamage: 0}, nil, hand).Value
	if got != 5 {
		t.Fatalf("Value = %d, want 5 (Arcanic Crackle 3 + Lunging Press +1 buff + 1 arcane)", got)
	}
}

// Tests that an AR with no legal target in hand can't enter the chain.
func TestAttackReaction_NoTargetAtAllNothingHappens(t *testing.T) {
	d := deck.New(hero.Viserai{}, nil, fillerDeck())
	hand := []deck.Card{cards.LungingPressBlue{}}
	got := sim.EvalOneTurnForTesting(d, sim.Matchup{IncomingDamage: 0}, nil, hand).Value
	if got != 0 {
		t.Fatalf("Value = %d, want 0 (no target, AR can't play)", got)
	}
}

// Tests that an Attack Reaction can't target another Attack Reaction.
func TestAttackReaction_CantTargetAnotherAR(t *testing.T) {
	d := deck.New(hero.Viserai{}, nil, fillerDeck())
	hand := []deck.Card{cards.LungingPressBlue{}, cards.LungingPressBlue{}}
	got := sim.EvalOneTurnForTesting(d, sim.Matchup{IncomingDamage: 0}, nil, hand).Value
	if got != 0 {
		t.Fatalf("Value = %d, want 0 (ARs can't target each other)", got)
	}
}

// Tests that Attack Reactions can't play before the attack they target (and thereby trigger
// Viserai's hero ability).
func TestAttackReaction_ReactionsComeAfterAttacks(t *testing.T) {
	d := deck.New(hero.Viserai{}, nil, fillerDeck())
	hand := []deck.Card{cards.LungingPressBlue{}, cards.HocusPocusRed{}}
	got := sim.EvalOneTurnForTesting(d, sim.Matchup{IncomingDamage: 0}, nil, hand).Value
	if got != 5 {
		t.Fatalf("Value = %d, want 5 (HP 3 + LP +1 buff + HP runechant 1; Viserai must not fire)", got)
	}
}

// Tests that an AR can only target an attack that is actually played, not one that's been
// pitched.
func TestAttackReaction_PitchedAttackIsNotATarget(t *testing.T) {
	d := deck.New(hero.Viserai{}, nil, fillerDeck())
	hand := []deck.Card{
		cards.LungingPressBlue{},
		testutils.GenericAttack(1, 0),
	}
	got := sim.EvalOneTurnForTesting(d, sim.Matchup{IncomingDamage: 0}, nil, hand).Value
	if got != 0 {
		t.Fatalf("Value = %d, want 0 (attack pitched ⇒ AR has nothing to buff)", got)
	}
}

// Tests that Thrust's +3{p} buff lands on a swinging Sword weapon.
func TestAttackReaction_ThrustBuffsSwingingSwordWeapon(t *testing.T) {
	d := deck.New(hero.Viserai{}, []deck.Weapon{weapons.NebulaBlade{}}, fillerDeck())
	hand := []deck.Card{
		cards.ThrustRed{},
		cards.ToughenUpBlue{},
	}
	got := sim.EvalOneTurnForTesting(d, sim.Matchup{IncomingDamage: 0}, nil, hand).Value
	if got != 5 {
		t.Fatalf("Value = %d, want 5 (Nebula Blade 1 + Thrust +3 buff + runechant 1)", got)
	}
}

// Tests that Blade Flash on Nebula Blade's swing chains a follow-up no-go-again attack action.
func TestAttackReaction_BladeFlashGoAgainExtendsChain(t *testing.T) {
	d := deck.New(hero.Viserai{}, []deck.Weapon{weapons.NebulaBlade{}}, fillerDeck())
	hand := []deck.Card{
		cards.BladeFlashBlue{},
		cards.ToughenUpBlue{},
		testutils.NoGoAgainAttackStub{},
	}
	got := sim.EvalOneTurnForTesting(d, sim.Matchup{IncomingDamage: 0}, nil, hand).Value
	if got != 3 {
		t.Fatalf("Value = %d, want 3 (Nebula 1 + runechant 1 + chained 1-power attack)", got)
	}
}

// Tests that without Blade Flash the same hand can't chain past Nebula Blade's swing.
func TestAttackReaction_NebulaSwingWithoutBladeFlashCantChain(t *testing.T) {
	d := deck.New(hero.Viserai{}, []deck.Weapon{weapons.NebulaBlade{}}, fillerDeck())
	hand := []deck.Card{
		cards.ToughenUpBlue{},
		testutils.NoGoAgainAttackStub{},
	}
	got := sim.EvalOneTurnForTesting(d, sim.Matchup{IncomingDamage: 0}, nil, hand).Value
	if got != 2 {
		t.Fatalf("Value = %d, want 2 (Nebula 1 + runechant 1; no AP for the follow-up)", got)
	}
}
