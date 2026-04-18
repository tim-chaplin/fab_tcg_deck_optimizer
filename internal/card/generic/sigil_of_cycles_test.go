package generic

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

// TestSigilOfCycles_SetsAuraCreated verifies the Blue-only variant flips AuraCreated and returns 0.
func TestSigilOfCycles_SetsAuraCreated(t *testing.T) {
	s := card.TurnState{}
	if got := (SigilOfCyclesBlue{}).Play(&s); got != 0 {
		t.Errorf("Play() = %d, want 0", got)
	}
	if !s.AuraCreated {
		t.Error("AuraCreated = false, want true")
	}
}
