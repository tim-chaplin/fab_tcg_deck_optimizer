package generic

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

// TestSigilOfFyendal_SetsAuraCreatedAndGainsHealth verifies the Blue-only variant flips
// AuraCreated (same-turn aura check) and returns 1 for the guaranteed 1{h} gain on leave.
func TestSigilOfFyendal_SetsAuraCreatedAndGainsHealth(t *testing.T) {
	s := card.TurnState{}
	if got := (SigilOfFyendalBlue{}).Play(&s); got != 1 {
		t.Errorf("Play() = %d, want 1 (1{h} on leave)", got)
	}
	if !s.AuraCreated {
		t.Error("AuraCreated = false, want true")
	}
}
