package cards

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
)

// Tests that Exposed accepts attack action cards as targets.
func TestExposed_AcceptsAttackAction(t *testing.T) {
	if !(ExposedBlue{}).ARTargetAllowed(testutils.GenericAttack(0, 4), 0) {
		t.Error("attack action should be a legal target")
	}
}

// Tests that Exposed accepts weapon swings as targets.
func TestExposed_AcceptsWeaponAttack(t *testing.T) {
	if !(ExposedBlue{}).ARTargetAllowed(testutils.RunebladeWeapon{}, 0) {
		t.Error("weapon swing should be a legal target")
	}
}

// Tests that Exposed rejects non-attack cards.
func TestExposed_RejectsNonAttack(t *testing.T) {
	if (ExposedBlue{}).ARTargetAllowed(testutils.NonAttack{}, 0) {
		t.Error("non-attack should be rejected")
	}
}

// Tests that Exposed's Play marks the opposing hero.
func TestExposed_PlayMarksOpponent(t *testing.T) {
	s := sim.TurnState{}
	(ExposedBlue{}).Play(&s, s.Logger(), &sim.CardState{Card: ExposedBlue{}})
	if !s.OpponentMarked {
		t.Error("OpponentMarked = false after Play, want true")
	}
}
