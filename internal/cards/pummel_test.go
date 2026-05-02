package cards

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
)

// Tests that the predicate accepts a cost-≥2 attack action.
func TestPummel_PredicateAcceptsCostTwoAttackAction(t *testing.T) {
	if !(PummelRed{}).ARTargetAllowed(testutils.GenericAttack(2, 4)) {
		t.Error("cost-2 attack action should be a legal target")
	}
}

// Tests that the predicate rejects a cost-1 attack action.
func TestPummel_PredicateRejectsCostOneAttack(t *testing.T) {
	if (PummelRed{}).ARTargetAllowed(testutils.GenericAttack(1, 4)) {
		t.Error("cost-1 attack action shouldn't match the cost-≥2 gate")
	}
}

// Tests that the predicate rejects a non-attack action.
func TestPummel_PredicateRejectsNonAttackAction(t *testing.T) {
	if (PummelRed{}).ARTargetAllowed(testutils.GenericAction()) {
		t.Error("non-attack action shouldn't match")
	}
}
