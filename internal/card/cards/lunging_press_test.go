package cards

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
)

// Tests that Lunging Press's predicate accepts an attack action card.
func TestLungingPress_AcceptsAttackAction(t *testing.T) {
	if !(LungingPressBlue{}).ARTargetAllowed(nil, &card.CardState{Card: testutils.FakeRedAttack()}, 0) {
		t.Error("attack action card should be a legal target")
	}
}

// Tests that the predicate rejects a plain non-attack Action.
func TestLungingPress_RejectsNonAttackAction(t *testing.T) {
	if (LungingPressBlue{}).ARTargetAllowed(nil, &card.CardState{Card: testutils.FakeRedAction()}, 0) {
		t.Error("non-attack action should be rejected")
	}
}

// Tests that the predicate rejects weapons (printed text is "attack action card", not
// "attack"). The fake is a Resource+Weapon — only the absence of TypeAttack matters to
// the predicate.
func TestLungingPress_RejectsWeapon(t *testing.T) {
	weapon := testutils.FakeRedResource().WithTypes(card.TypeWeapon)
	if (LungingPressBlue{}).ARTargetAllowed(nil, &card.CardState{Card: weapon}, 0) {
		t.Error("weapon should be rejected")
	}
}
