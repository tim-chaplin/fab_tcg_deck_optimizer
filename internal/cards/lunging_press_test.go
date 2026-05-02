package cards

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
)

// Tests that Lunging Press fizzles silently when no target follows it in the chain.
func TestLungingPress_NoTargetFizzles(t *testing.T) {
	s := sim.TurnState{}
	(LungingPressBlue{}).Play(&s, &sim.CardState{Card: LungingPressBlue{}})
	if s.Value != 0 {
		t.Errorf("Play() Value = %d, want 0", s.Value)
	}
}

// Tests that Lunging Press lands +1{p} on the first attack action card in CardsRemaining.
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

// Tests that a plain non-attack Action is skipped (predicate requires Action+Attack).
func TestLungingPress_SkipsNonAttackAction(t *testing.T) {
	target := &sim.CardState{Card: testutils.GenericAction()}
	s := sim.TurnState{CardsRemaining: []*sim.CardState{target}}
	(LungingPressBlue{}).Play(&s, &sim.CardState{Card: LungingPressBlue{}})
	if target.BonusAttack != 0 {
		t.Errorf("non-attack target BonusAttack = %d, want 0", target.BonusAttack)
	}
}

// Tests that a weapon is skipped (predicate requires "attack action card", not "attack").
func TestLungingPress_SkipsWeapons(t *testing.T) {
	target := &sim.CardState{Card: testutils.RunebladeWeapon{}}
	s := sim.TurnState{CardsRemaining: []*sim.CardState{target}}
	(LungingPressBlue{}).Play(&s, &sim.CardState{Card: LungingPressBlue{}})
	if target.BonusAttack != 0 {
		t.Errorf("weapon target BonusAttack = %d, want 0", target.BonusAttack)
	}
}
