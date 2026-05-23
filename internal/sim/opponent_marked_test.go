package sim

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/gameengine"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/token"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/triggertype"
)

// Tests that the runechant aura handler leaves OpponentMarked alone — arcane damage
// doesn't strip the mark, only physical attacks do.
func TestRunechantAuraHandler_LeavesOpponentMarked(t *testing.T) {
	ge := &gameengine.GameEngine{GameState: gameengine.GameStateBuilder().SetOpponentMarked(true).Build()}
	ge.AppendAura(token.NewRunechant(1))
	// The Runechant aura fires on triggertype.CardOrAbility filtered to attacks, so the
	// firing card must be an attack for its IsAttack filter to match.
	ge.FireTriggers(triggertype.CardOrAbility, FakeRedAttack{})
	if !ge.OpponentMarked() {
		t.Error("OpponentMarked = false after runechant pop, want true (arcane doesn't clear mark)")
	}
}
