package cards

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
)

// Tests that the predicate accepts a sword weapon attack (mode-0 leg).
func TestRazorReflex_PredicateAcceptsSwordWeapon(t *testing.T) {
	swordWeapon := testutils.NewStubCard("sword").
		WithTypes(card.NewTypeSet(card.TypeGeneric, card.TypeWeapon, card.TypeSword))
	if !(RazorReflexRed{}).ARTargetAllowed(swordWeapon) {
		t.Error("sword weapon should pass mode-0 leg of predicate")
	}
}

// Tests that the predicate accepts a cost-1 attack action (mode-1 leg).
func TestRazorReflex_PredicateAcceptsCostOneAttackAction(t *testing.T) {
	if !(RazorReflexRed{}).ARTargetAllowed(testutils.GenericAttack(1, 4)) {
		t.Error("cost-1 attack action should pass mode-1 leg of predicate")
	}
}

// Tests that a non-sword cost-2 attack action fails both legs.
func TestRazorReflex_PredicateRejectsCostTwoNonSword(t *testing.T) {
	if (RazorReflexRed{}).ARTargetAllowed(testutils.GenericAttack(2, 4)) {
		t.Error("cost-2 generic attack action shouldn't match either mode")
	}
}
