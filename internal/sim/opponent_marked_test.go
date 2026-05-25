package sim

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/gameengine"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/token"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/triggertype"
)

// Tests that a Runechant pop whose count clears the damage-likely-to-hit window strips
// OpponentMarked — arcane damage consumes the mark like any other damage.
func TestRunechantAuraHandler_ClearsOpponentMarkedOnLandingArcane(t *testing.T) {
	ge := &gameengine.GameEngine{GameState: gameengine.GameStateBuilder().
		SetOpponentMarked(true).
		AddAura(token.NewRunechant(1)).
		Build()}
	// The Runechant aura fires on triggertype.CardOrAbility filtered to attacks, so the
	// firing card must be an attack for its IsAttack filter to match. Count 1 clears
	// LikelyDamageHits(1, false) so the pop lands as arcane damage.
	ge.FireTriggers(triggertype.CardOrAbility, testutils.FakeRedAttack())
	if ge.OpponentMarked() {
		t.Error("OpponentMarked = true after landing arcane pop, want false (arcane damage strips the mark)")
	}
}

// Tests that a Runechant pop whose count doesn't clear the damage-likely-to-hit window
// leaves OpponentMarked alone — the damage is fully prevented in our model so the mark
// isn't consumed.
func TestRunechantAuraHandler_LeavesOpponentMarkedWhenArcaneFizzles(t *testing.T) {
	ge := &gameengine.GameEngine{GameState: gameengine.GameStateBuilder().
		SetOpponentMarked(true).
		AddAura(token.NewRunechant(2)).
		Build()}
	// Count 2 fails LikelyDamageHits(2, false) so no damage lands.
	ge.FireTriggers(triggertype.CardOrAbility, testutils.FakeRedAttack())
	if !ge.OpponentMarked() {
		t.Error("OpponentMarked = false after fizzled arcane pop, want true (no damage landed, mark survives)")
	}
}
