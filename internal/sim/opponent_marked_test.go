package sim

import "testing"

// Tests that the runechant aura handler leaves OpponentMarked alone — arcane damage doesn't
// strip the mark, only physical attacks do.
func TestRunechantAuraHandler_LeavesOpponentMarked(t *testing.T) {
	s := NewTurnStatePtr(TurnStateSpec{OpponentMarked: true, Auras: []Aura{NewRunechantAura(1)}})
	s.FireAuraForTesting(0)
	if !s.OpponentMarked() {
		t.Error("OpponentMarked = false after runechant pop, want true (arcane doesn't clear mark)")
	}
}
