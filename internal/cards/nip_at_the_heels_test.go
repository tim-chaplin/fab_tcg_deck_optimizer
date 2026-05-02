package cards

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
)

// Tests that Nip at the Heels accepts a base-≤3 attack action.
func TestNipAtTheHeels_AcceptsLowPowerAttack(t *testing.T) {
	if !(NipAtTheHeelsBlue{}).ARTargetAllowed(testutils.GenericAttack(0, 3)) {
		t.Error("3-power attack should be a legal target")
	}
}

// Tests that the base-power gate reads printed Attack(): a 4-power attack is rejected.
func TestNipAtTheHeels_RejectsHighPowerAttack(t *testing.T) {
	if (NipAtTheHeelsBlue{}).ARTargetAllowed(testutils.GenericAttack(0, 4)) {
		t.Error("4-power attack should be rejected (base power > 3)")
	}
}

// Tests that the predicate accepts low-power weapons too.
func TestNipAtTheHeels_AcceptsLowPowerWeapon(t *testing.T) {
	weapon := testutils.RunebladeWeapon{}
	if !(NipAtTheHeelsBlue{}).ARTargetAllowed(weapon) {
		t.Error("0-power weapon should be a legal target (≤ 3)")
	}
}
