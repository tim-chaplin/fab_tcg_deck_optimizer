package sim

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/v2/gameengine"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/token"
)

// Tests that the runechant aura handler leaves OpponentMarked alone — arcane damage
// doesn't strip the mark, only physical attacks do.
func TestRunechantAuraHandler_LeavesOpponentMarked(t *testing.T) {
	ge := &gameengine.GameEngine{GameState: gameengine.GameStateBuilder().SetOpponentMarked(true).Build()}
	ge.CreateAura(token.NewRunechant(1))
	// Fire the runechant aura via the engine's TriggerAttack fire walk (the runechant aura
	// is registered as TriggerAttack); pass a nil triggering card since the runechant
	// handler doesn't read TriggeringCard.
	ge.FireAttack(nil)
	if !ge.OpponentMarked() {
		t.Error("OpponentMarked = false after runechant pop, want true (arcane doesn't clear mark)")
	}
}
