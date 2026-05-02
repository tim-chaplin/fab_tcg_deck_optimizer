package cards

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
)

// Tests that Lunging Press fizzles silently when no attack action card follows it in the
// chain — orderings the partition validator already let through but where the AR happens
// to play after every legal target.
func TestLungingPress_NoTargetFizzles(t *testing.T) {
	s := sim.TurnState{}
	(LungingPressBlue{}).Play(&s, &sim.CardState{Card: LungingPressBlue{}})
	if s.Value != 0 {
		t.Errorf("Play() Value = %d, want 0", s.Value)
	}
}

// Tests that Lunging Press lands its +1{p} on the first attack action card in
// CardsRemaining — the canonical "AR buffs the next target" path.
func TestLungingPress_BuffsNextAttackAction(t *testing.T) {
	target := &sim.CardState{Card: testutils.GenericAttack(0, 0)}
	s := sim.TurnState{CardsRemaining: []*sim.CardState{target}}
	(LungingPressBlue{}).Play(&s, &sim.CardState{Card: LungingPressBlue{}})
	if s.Value != 0 {
		t.Errorf("Play() Value = %d, want 0 (AR contributes 0; buff rides on target)", s.Value)
	}
	if target.BonusAttack != 1 {
		t.Errorf("target BonusAttack = %d, want 1", target.BonusAttack)
	}
}

// Tests that a non-attack action card in CardsRemaining gets skipped — Lunging Press's
// printed predicate says "attack action card", so a plain Action doesn't qualify.
func TestLungingPress_SkipsNonAttackAction(t *testing.T) {
	target := &sim.CardState{Card: testutils.GenericAction()}
	s := sim.TurnState{CardsRemaining: []*sim.CardState{target}}
	(LungingPressBlue{}).Play(&s, &sim.CardState{Card: LungingPressBlue{}})
	if target.BonusAttack != 0 {
		t.Errorf("non-attack target BonusAttack = %d, want 0", target.BonusAttack)
	}
}

// Tests that a weapon swing doesn't qualify as a target — Lunging Press says "attack action
// card", which excludes weapons even though they're attacks.
func TestLungingPress_SkipsWeapons(t *testing.T) {
	target := &sim.CardState{Card: testutils.RunebladeWeapon{}}
	s := sim.TurnState{CardsRemaining: []*sim.CardState{target}}
	(LungingPressBlue{}).Play(&s, &sim.CardState{Card: LungingPressBlue{}})
	if target.BonusAttack != 0 {
		t.Errorf("weapon target BonusAttack = %d, want 0", target.BonusAttack)
	}
}
