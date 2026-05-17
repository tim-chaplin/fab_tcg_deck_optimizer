package cards

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
)

// Tests that Thrust's predicate accepts a sword attack action card.
func TestThrust_AcceptsSwordAttackAction(t *testing.T) {
	swordAction := testutils.NewStubCard("SwordAction").
		WithTypes(card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeAttack, card.TypeSword))
	if !(ThrustRed{}).ARTargetAllowed(nil, swordAction, 0) {
		t.Error("sword action card should be a legal target")
	}
}

// Tests that Thrust's predicate accepts a Sword weapon swing.
func TestThrust_AcceptsSwordWeapon(t *testing.T) {
	swordSwing := testutils.NewStubCard("SwordWeaponAbility").
		WithTypes(card.NewTypeSet(card.TypeGeneric, card.TypeWeapon, card.TypeSword, card.TypeAttack))
	if !(ThrustRed{}).ARTargetAllowed(nil, swordSwing, 0) {
		t.Error("sword weapon swing should be a legal target")
	}
}

// Tests that a non-Sword attack is rejected.
func TestThrust_RejectsNonSwordAttack(t *testing.T) {
	if (ThrustRed{}).ARTargetAllowed(nil, testutils.GenericAttack(0, 0), 0) {
		t.Error("non-sword attack action should be rejected")
	}
}
