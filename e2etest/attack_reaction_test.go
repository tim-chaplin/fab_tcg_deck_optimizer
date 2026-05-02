package e2etest

// Tests that Attack Reaction partition validation flows end-to-end: a hand with an AR but
// no legal target rejects every attack-role assignment that includes the AR, while the
// same hand with a target accepts it and the +N{p} buff lands on damage output.

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/cards"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/heroes"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/weapons"
)

// Tests that the Lunging Press +1{p} buff lands on a paired attack action card. Arcanic
// Crackle costs 0 so the AR-played chain spends no resources, forcing the optimizer to
// pick the play-both line (chain runs free) over pitching the AR for resource it doesn't
// need. Buff applies in the [Lunging Press, Arcanic Crackle] ordering.
func TestAttackReaction_BuffLandsOnTarget(t *testing.T) {
	d := sim.New(heroes.Viserai{}, nil, fillerDeck())
	// Hand: Arcanic Crackle [R] (0 cost, 3 power, Go again) + Lunging Press [B] (AR, 0
	// cost, +1{p} target attack action). No pitch needed; both play in chain. Optimal
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

// Tests that with no attack target available at all, Lunging Press can't enter the chain:
// the only attack-role assignment leaves it Held, the optimizer can't milk a buff, and
// PrevTurnValue is zero.
func TestAttackReaction_NoTargetAtAllNothingHappens(t *testing.T) {
	d := sim.New(heroes.Viserai{}, nil, fillerDeck())
	// Hand: just Lunging Press [B]. No attack action card in hand to target. The partition
	// validator rejects the {Lunging Press → Attack} assignment, leaving only Held.
	hand := []sim.Card{cards.LungingPressBlue{}}
	got := d.EvalOneTurnForTesting(0, nil, hand).PrevTurnValue
	if got != 0 {
		t.Fatalf("PrevTurnValue = %d, want 0 (no target, AR can't play)", got)
	}
}

// Tests that Thrust's +3{p} buff lands on a swinging sword weapon — exercises the wmask
// path through partitionHasValidARTargets, which sees the weapon only when wmask
// includes it. Equip Nebula Blade (Sword), draft a hand whose only attack candidate is
// the weapon swing, and confirm the AR-buffed swing damages for printed-power + 3 plus
// Nebula Blade's on-hit Runechant credit.
func TestAttackReaction_ThrustBuffsSwingingSwordWeapon(t *testing.T) {
	d := sim.New(heroes.Viserai{}, []sim.Weapon{weapons.NebulaBlade{}}, fillerDeck())
	// Hand: Thrust [R] (cost 1, pitch 1) + ToughenUp [B] (pitches 3 — funds Nebula Blade's
	// 2-cost swing + Thrust's 1-cost). Nebula Blade is the only attack target; as a Sword
	// Weapon it satisfies Thrust's predicate. Optimal chain pitches Toughen Up, plays
	// Thrust then swings Nebula Blade. Damage: Nebula Blade 1{p} + Thrust +3 = 4 hit, plus
	// +1 from the on-hit Runechant = 5.
	hand := []sim.Card{
		cards.ThrustRed{},
		cards.ToughenUpBlue{},
	}
	got := d.EvalOneTurnForTesting(0, nil, hand).PrevTurnValue
	if got != 5 {
		t.Fatalf("PrevTurnValue = %d, want 5 (Nebula Blade 1 + Thrust +3 buff + runechant 1)", got)
	}
}
