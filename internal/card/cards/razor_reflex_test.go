package cards

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
)

// Tests that mode 0 accepts a sword weapon attack.
func TestRazorReflex_Mode0AcceptsSwordWeapon(t *testing.T) {
	swordSwing := testutils.FakeWeaponSwing().WithTypes(card.TypeSword)
	if !(RazorReflexRed{}).ARTargetAllowed(nil, &card.CardState{Card: swordSwing}, 0) {
		t.Error("mode 0 should accept a sword weapon swing")
	}
}

// Tests that mode 0 rejects a non-sword target.
func TestRazorReflex_Mode0RejectsNonSword(t *testing.T) {
	if (RazorReflexRed{}).ARTargetAllowed(nil, &card.CardState{Card: testutils.FakeRedAttack().WithCost(1)}, 0) {
		t.Error("mode 0 should reject a non-sword attack")
	}
}

// Tests that mode 0 rejects a sword-typed attack action card — the printed text says
// "weapon attack", so an action card sharing the Sword subtype shouldn't qualify.
func TestRazorReflex_Mode0RejectsSwordAttackActionCard(t *testing.T) {
	swordAction := testutils.FakeRedAttack().WithTypes(card.TypeSword)
	if (RazorReflexRed{}).ARTargetAllowed(nil, &card.CardState{Card: swordAction}, 0) {
		t.Error("mode 0 should reject a sword attack action card (only weapon attacks qualify)")
	}
}

// Tests that mode 1 accepts a cost-≤1 attack action.
func TestRazorReflex_Mode1AcceptsCostOneAttackAction(t *testing.T) {
	if !(RazorReflexRed{}).ARTargetAllowed(nil, &card.CardState{Card: testutils.FakeRedAttack().WithCost(1)}, 1) {
		t.Error("mode 1 should accept a cost-1 attack action")
	}
}

// Tests that mode 1 rejects a cost-≥2 attack action.
func TestRazorReflex_Mode1RejectsCostTwoAttack(t *testing.T) {
	if (RazorReflexRed{}).ARTargetAllowed(nil, &card.CardState{Card: testutils.FakeRedAttack().WithCost(2)}, 1) {
		t.Error("mode 1 should reject cost-2 attack actions")
	}
}
