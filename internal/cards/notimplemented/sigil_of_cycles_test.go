package notimplemented

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

// TestSigilOfCycles_SetsAuraCreated verifies the Blue-only variant flips AuraCreated and returns 0.
func TestSigilOfCycles_SetsAuraCreated(t *testing.T) {
	s := sim.TurnState{}
	(SigilOfCyclesBlue{}).Play(&s, &sim.CardState{Card: SigilOfCyclesBlue{}})
	if got := s.Value; got != 0 {
		t.Errorf("Play() = %d, want 0", got)
	}
	if !s.AuraCreated {
		t.Error("AuraCreated = false, want true")
	}
}
