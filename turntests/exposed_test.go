package turntests

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card/cards"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/gameengine"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
)

// Tests that Exposed accepts attack action cards as targets.
func TestExposed_AcceptsAttackAction(t *testing.T) {
	if !(cards.ExposedBlue{}).ARTargetAllowed(nil, testutils.GenericAttack(0, 4), 0) {
		t.Error("attack action should be a legal target")
	}
}

// Tests that Exposed accepts weapon swings as targets.
func TestExposed_AcceptsWeaponAttack(t *testing.T) {
	if !(cards.ExposedBlue{}).ARTargetAllowed(nil, testutils.RunebladeWeapon{}, 0) {
		t.Error("weapon swing should be a legal target")
	}
}

// Tests that Exposed rejects non-attack cards.
func TestExposed_RejectsNonAttack(t *testing.T) {
	if (cards.ExposedBlue{}).ARTargetAllowed(nil, testutils.NonAttack{}, 0) {
		t.Error("non-attack should be rejected")
	}
}

// Tests that Exposed's Play marks the opposing hero.
func TestExposed_PlayMarksOpponent(t *testing.T) {
	ge := gameengine.New()
	ge.ResolveChainStep(ge.Logger(), &card.CardState{Card: cards.ExposedBlue{}})
	if !ge.OpponentMarked() {
		t.Error("OpponentMarked = false after Play, want true")
	}
}
