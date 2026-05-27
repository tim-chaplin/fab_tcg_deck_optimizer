package cards

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
)

// Tests that Thrust's predicate accepts a sword attack action card.
func TestThrust_AcceptsSwordAttackAction(t *testing.T) {
	swordAction := testutils.FakeRedAttack().WithTypes(card.TypeSword)
	if !(ThrustRed{}).ARTargetAllowed(nil, &card.CardState{Card: swordAction}, 0) {
		t.Error("sword action card should be a legal target")
	}
}

// Tests that Thrust's predicate accepts a Sword weapon swing.
func TestThrust_AcceptsSwordWeapon(t *testing.T) {
	swordSwing := testutils.FakeWeaponSwing().WithTypes(card.TypeSword)
	if !(ThrustRed{}).ARTargetAllowed(nil, &card.CardState{Card: swordSwing}, 0) {
		t.Error("sword weapon swing should be a legal target")
	}
}

// Tests that a non-Sword attack is rejected.
func TestThrust_RejectsNonSwordAttack(t *testing.T) {
	if (ThrustRed{}).ARTargetAllowed(nil, &card.CardState{Card: testutils.FakeRedAttack()}, 0) {
		t.Error("non-sword attack action should be rejected")
	}
}
