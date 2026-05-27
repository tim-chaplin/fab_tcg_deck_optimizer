package cards

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
)

// Tests that Blade Flash's predicate accepts a sword attack action card.
func TestBladeFlash_AcceptsSwordAttackAction(t *testing.T) {
	swordAction := testutils.FakeRedAttack().WithTypes(card.TypeSword)
	if !(BladeFlashBlue{}).ARTargetAllowed(nil, &card.CardState{Card: swordAction}, 0) {
		t.Error("sword action card should be a legal target")
	}
}

// Tests that Blade Flash's predicate accepts a Sword weapon swing.
func TestBladeFlash_AcceptsSwordWeapon(t *testing.T) {
	swordSwing := testutils.FakeWeaponSwing().WithTypes(card.TypeSword)
	if !(BladeFlashBlue{}).ARTargetAllowed(nil, &card.CardState{Card: swordSwing}, 0) {
		t.Error("sword weapon swing should be a legal target")
	}
}

// Tests that a non-Sword attack is rejected.
func TestBladeFlash_RejectsNonSwordAttack(t *testing.T) {
	if (BladeFlashBlue{}).ARTargetAllowed(nil, &card.CardState{Card: testutils.FakeRedAttack()}, 0) {
		t.Error("non-sword attack action should be rejected")
	}
}

// Tests that a non-attack sword card is rejected. FakeRedResource has the Generic +
// Resource type line — Resource isn't the predicate's concern; what matters is that
// TypeAttack is absent so the rejection branch fires.
func TestBladeFlash_RejectsNonAttackSword(t *testing.T) {
	swordEquipment := testutils.FakeRedResource().WithTypes(card.TypeSword)
	if (BladeFlashBlue{}).ARTargetAllowed(nil, &card.CardState{Card: swordEquipment}, 0) {
		t.Error("non-attack sword should be rejected")
	}
}
