package cards

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

// Tests that Performance Bonus registers an OnHit handler. End-to-end Gold creation is
// covered by the e2e suite.
func TestPerformanceBonus_RegistersOnHit(t *testing.T) {
	s := sim.NewTurnState(nil, nil)
	self := &sim.CardState{Card: PerformanceBonusBlue{}}
	(PerformanceBonusBlue{}).Play(s, self)
	if len(self.OnHit) != 1 {
		t.Fatalf("OnHit handlers = %d, want 1 (gold-create rider registered)", len(self.OnHit))
	}
}

// Tests that the from-arsenal play grants Go again on top of the on-hit rider.
func TestPerformanceBonus_ArsenalGrantsGoAgain(t *testing.T) {
	s := sim.NewTurnState(nil, nil)
	self := &sim.CardState{Card: PerformanceBonusRed{}, FromArsenal: true}
	(PerformanceBonusRed{}).Play(s, self)
	if !self.GrantedGoAgain {
		t.Fatalf("GrantedGoAgain = false, want true (played from arsenal)")
	}
}

// Tests that hand-played (non-arsenal) Performance Bonus does NOT grant Go again.
func TestPerformanceBonus_NonArsenalNoGoAgain(t *testing.T) {
	s := sim.NewTurnState(nil, nil)
	self := &sim.CardState{Card: PerformanceBonusRed{}, FromArsenal: false}
	(PerformanceBonusRed{}).Play(s, self)
	if self.GrantedGoAgain {
		t.Fatalf("GrantedGoAgain = true with FromArsenal=false, want false")
	}
}
