package cards

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
)

// Tests that Blade Flash's predicate accepts a sword attack action card.
func TestBladeFlash_AcceptsSwordAttackAction(t *testing.T) {
	swordAction := testutils.NewFakeCard("SwordAction").
		WithTypes(card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeAttack, card.TypeSword))
	if !(BladeFlashBlue{}).ARTargetAllowed(nil, swordAction, 0) {
		t.Error("sword action card should be a legal target")
	}
}

// Tests that Blade Flash's predicate accepts a Sword weapon swing.
func TestBladeFlash_AcceptsSwordWeapon(t *testing.T) {
	swordSwing := testutils.NewFakeCard("SwordWeaponAbility").
		WithTypes(card.NewTypeSet(card.TypeGeneric, card.TypeWeapon, card.TypeSword, card.TypeAttack))
	if !(BladeFlashBlue{}).ARTargetAllowed(nil, swordSwing, 0) {
		t.Error("sword weapon swing should be a legal target")
	}
}

// Tests that a non-Sword attack is rejected.
func TestBladeFlash_RejectsNonSwordAttack(t *testing.T) {
	if (BladeFlashBlue{}).ARTargetAllowed(nil, testutils.GenericAttack(0, 0), 0) {
		t.Error("non-sword attack action should be rejected")
	}
}

// Tests that a non-attack sword card is rejected.
func TestBladeFlash_RejectsNonAttackSword(t *testing.T) {
	swordEquipment := testutils.NewFakeCard("SwordEquipment").
		WithTypes(card.NewTypeSet(card.TypeGeneric, card.TypeSword))
	if (BladeFlashBlue{}).ARTargetAllowed(nil, swordEquipment, 0) {
		t.Error("non-attack sword should be rejected")
	}
}
