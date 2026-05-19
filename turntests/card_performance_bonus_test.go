package turntests

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card/cards"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/gameengine"
)

// Tests that Performance Bonus registers an OnHit handler.
func TestPerformanceBonus_RegistersOnHit(t *testing.T) {
	ge := gameengine.New()
	pc := &card.CardState{Card: cards.PerformanceBonusBlue{}}
	ge.ResolveChainStep(ge.Logger(), pc)
	if len(pc.OnHit) != 1 {
		t.Fatalf("OnHit handlers = %d, want 1 (gold-create rider registered)", len(pc.OnHit))
	}
}

// Tests that the from-arsenal play grants Go again on top of the on-hit rider.
func TestPerformanceBonus_ArsenalGrantsGoAgain(t *testing.T) {
	ge := gameengine.New()
	pc := &card.CardState{Card: cards.PerformanceBonusRed{}, FromArsenal: true}
	ge.ResolveChainStep(ge.Logger(), pc)
	if !pc.GrantedGoAgain {
		t.Fatalf("GrantedGoAgain = false, want true (played from arsenal)")
	}
}

// Tests that hand-played (non-arsenal) Performance Bonus does NOT grant Go again.
func TestPerformanceBonus_NonArsenalNoGoAgain(t *testing.T) {
	ge := gameengine.New()
	pc := &card.CardState{Card: cards.PerformanceBonusRed{}, FromArsenal: false}
	ge.ResolveChainStep(ge.Logger(), pc)
	if pc.GrantedGoAgain {
		t.Fatalf("GrantedGoAgain = true with FromArsenal=false, want false")
	}
}
