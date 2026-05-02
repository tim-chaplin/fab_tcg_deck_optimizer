package cards

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
)

// Tests that Exposed's +1{p} buff lands on the next attack action card.
func TestExposed_BuffsAttackAction(t *testing.T) {
	target := &sim.CardState{Card: testutils.GenericAttack(0, 0)}
	s := sim.TurnState{CardsRemaining: []*sim.CardState{target}}
	(ExposedBlue{}).Play(&s, &sim.CardState{Card: ExposedBlue{}})
	if target.BonusAttack != 1 {
		t.Errorf("attack-action target BonusAttack = %d, want 1", target.BonusAttack)
	}
}

// Tests that Exposed accepts weapon attacks too — the printed "target attack" includes
// weapons.
func TestExposed_AcceptsWeapon(t *testing.T) {
	if !(ExposedBlue{}).ARTargetAllowed(testutils.RunebladeWeapon{}) {
		t.Error("weapon should be a legal target for Exposed")
	}
}

// Tests that a non-attack card is rejected.
func TestExposed_RejectsNonAttack(t *testing.T) {
	if (ExposedBlue{}).ARTargetAllowed(testutils.GenericAction()) {
		t.Error("non-attack action should be rejected")
	}
}
