package sim

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/v2/gameengine"
)

// Tests that the runechant aura handler leaves OpponentMarked alone — arcane damage
// doesn't strip the mark, only physical attacks do.
func TestRunechantAuraHandler_LeavesOpponentMarked(t *testing.T) {
	s := gameengine.NewFromSpec(gameengine.Spec{
		OpponentMarked: true,
	})
	s.CreateAura(NewRunechantAura(1))
	// Fire the runechant aura via the engine's TriggerAttack fire walk (the runechant aura
	// is registered as TriggerAttack); pass a nil triggering card since the runechant
	// handler doesn't read TriggeringCard.
	s.FireAttack(nil)
	if !s.OpponentMarked() {
		t.Error("OpponentMarked = false after runechant pop, want true (arcane doesn't clear mark)")
	}
}
