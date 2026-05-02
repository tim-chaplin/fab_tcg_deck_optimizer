package cards

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
)

// Tests that Nip at the Heels buffs a base-3 attack action by +1{p}.
func TestNipAtTheHeels_BuffsLowPowerAttack(t *testing.T) {
	target := &sim.CardState{Card: testutils.GenericAttack(0, 3)}
	s := sim.TurnState{CardsRemaining: []*sim.CardState{target}}
	(NipAtTheHeelsBlue{}).Play(&s, &sim.CardState{Card: NipAtTheHeelsBlue{}})
	if target.BonusAttack != 1 {
		t.Errorf("3-power target BonusAttack = %d, want 1", target.BonusAttack)
	}
}

// Tests that a 4-power attack is rejected (base-power gate reads printed Attack()).
func TestNipAtTheHeels_RejectsHighPowerAttack(t *testing.T) {
	target := &sim.CardState{Card: testutils.GenericAttack(0, 4)}
	s := sim.TurnState{CardsRemaining: []*sim.CardState{target}}
	(NipAtTheHeelsBlue{}).Play(&s, &sim.CardState{Card: NipAtTheHeelsBlue{}})
	if target.BonusAttack != 0 {
		t.Errorf("4-power target BonusAttack = %d, want 0 (base power > 3)", target.BonusAttack)
	}
}

// Tests that the predicate accepts a low-power weapon.
func TestNipAtTheHeels_AcceptsLowPowerWeapon(t *testing.T) {
	weapon := testutils.RunebladeWeapon{}
	if !(NipAtTheHeelsBlue{}).ARTargetAllowed(weapon) {
		t.Error("0-power weapon should be a legal target (≤ 3)")
	}
}
