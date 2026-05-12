package turntests

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/cards"
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
)

// Tests that Performance Bonus registers an OnHit handler.
func TestPerformanceBonus_RegistersOnHit(t *testing.T) {
	s := sim.NewTurnStateFromCards(nil, nil)
	self := &card.CardState{Card: cards.PerformanceBonusBlue{}}
	sim.ResolveChainStep(s, s.Logger(), self)
	if len(self.OnHit) != 1 {
		t.Fatalf("OnHit handlers = %d, want 1 (gold-create rider registered)", len(self.OnHit))
	}
}

// Tests that the from-arsenal play grants Go again on top of the on-hit rider.
func TestPerformanceBonus_ArsenalGrantsGoAgain(t *testing.T) {
	s := sim.NewTurnStateFromCards(nil, nil)
	self := &card.CardState{Card: cards.PerformanceBonusRed{}, FromArsenal: true}
	sim.ResolveChainStep(s, s.Logger(), self)
	if !self.GrantedGoAgain {
		t.Fatalf("GrantedGoAgain = false, want true (played from arsenal)")
	}
}

// Tests that hand-played (non-arsenal) Performance Bonus does NOT grant Go again.
func TestPerformanceBonus_NonArsenalNoGoAgain(t *testing.T) {
	s := sim.NewTurnStateFromCards(nil, nil)
	self := &card.CardState{Card: cards.PerformanceBonusRed{}, FromArsenal: false}
	sim.ResolveChainStep(s, s.Logger(), self)
	if self.GrantedGoAgain {
		t.Fatalf("GrantedGoAgain = true with FromArsenal=false, want false")
	}
}
