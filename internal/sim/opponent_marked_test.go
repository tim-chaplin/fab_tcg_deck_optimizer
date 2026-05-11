package sim

import "testing"

// Tests that the runechant aura handler leaves OpponentMarked alone — arcane damage doesn't
// strip the mark, only physical attacks do.
func TestRunechantAuraHandler_LeavesOpponentMarked(t *testing.T) {
	s := &TurnState{OpponentMarked: true}
	aura := NewRunechantAura(1)
	s.Auras = append(s.Auras, aura)
	s.SetCurrentAuraIdxForTesting(0)
	runechantAuraHandler(s, s.logger, &s.Auras[0].Trigger, &s.Auras[0])
	if !s.OpponentMarked {
		t.Error("OpponentMarked = false after runechant pop, want true (arcane doesn't clear mark)")
	}
}
