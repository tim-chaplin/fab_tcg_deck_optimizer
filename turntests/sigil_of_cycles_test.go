package turntests

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/cards/notimplemented"

	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/gameengine"
)

// TestSigilOfCycles_SetsAuraCreated verifies the Blue-only variant flips AuraCreated and returns 0.
func TestSigilOfCycles_SetsAuraCreated(t *testing.T) {
	s := gameengine.New()
	s.ResolveChainStep(s.Logger(), &card.CardState{Card: notimplemented.SigilOfCyclesBlue{}})
	if got := s.Value(); got != 0 {
		t.Errorf("Play() = %d, want 0", got)
	}
	if !s.AuraCreated() {
		t.Error("AuraCreated = false, want true")
	}
}
