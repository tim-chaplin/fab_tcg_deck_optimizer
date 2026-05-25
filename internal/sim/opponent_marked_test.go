package sim

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/gameengine"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/token"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/triggertype"
)

// Tests that a Runechant pop strips OpponentMarked — arcane damage consumes the mark like
// any other damage, gated the same way as physical attacks (any positive damage, not
// LikelyDamageHits).
func TestRunechantAuraHandler_ClearsOpponentMarkedOnPop(t *testing.T) {
	ge := &gameengine.GameEngine{GameState: gameengine.GameStateBuilder().
		SetOpponentMarked(true).
		AddAura(token.NewRunechant(1)).
		Build()}
	// The Runechant aura fires on triggertype.CardOrAbility filtered to attacks, so the
	// firing card must be an attack for its IsAttack filter to match.
	ge.FireTriggers(triggertype.CardOrAbility, testutils.FakeRedAttack())
	if ge.OpponentMarked() {
		t.Error("OpponentMarked = true after Runechant pop, want false (arcane damage strips the mark)")
	}
}
