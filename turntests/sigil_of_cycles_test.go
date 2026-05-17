package turntests

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/v2/card/cards/notimplemented"

	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/gameengine"
)

// TestSigilOfCycles_SetsAuraCreated verifies the Blue-only variant flips AuraCreated and returns 0.
func TestSigilOfCycles_SetsAuraCreated(t *testing.T) {
	ge := gameengine.New()
	ge.ResolveChainStep(ge.Logger(), &card.CardState{Card: notimplemented.SigilOfCyclesBlue{}})
	if got := ge.Value(); got != 0 {
		t.Errorf("Play() = %d, want 0", got)
	}
	if !ge.AuraCreated() {
		t.Error("AuraCreated = false, want true")
	}
}
