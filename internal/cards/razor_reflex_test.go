package cards

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
)

// Tests that mode 0 accepts a sword weapon attack.
func TestRazorReflex_Mode0AcceptsSwordWeapon(t *testing.T) {
	swordWeapon := testutils.NewStubCard("sword").
		WithTypes(card.NewTypeSet(card.TypeGeneric, card.TypeWeapon, card.TypeSword))
	if !(RazorReflexRed{}).ARTargetAllowed(swordWeapon, 0) {
		t.Error("mode 0 should accept a sword weapon")
	}
}

// Tests that mode 0 rejects a non-sword target.
func TestRazorReflex_Mode0RejectsNonSword(t *testing.T) {
	if (RazorReflexRed{}).ARTargetAllowed(testutils.GenericAttack(1, 4), 0) {
		t.Error("mode 0 should reject a non-sword attack")
	}
}

// Tests that mode 1 accepts a cost-≤1 attack action.
func TestRazorReflex_Mode1AcceptsCostOneAttackAction(t *testing.T) {
	if !(RazorReflexRed{}).ARTargetAllowed(testutils.GenericAttack(1, 4), 1) {
		t.Error("mode 1 should accept a cost-1 attack action")
	}
}

// Tests that mode 1 rejects a cost-≥2 attack action.
func TestRazorReflex_Mode1RejectsCostTwoAttack(t *testing.T) {
	if (RazorReflexRed{}).ARTargetAllowed(testutils.GenericAttack(2, 4), 1) {
		t.Error("mode 1 should reject cost-2 attack actions")
	}
}
