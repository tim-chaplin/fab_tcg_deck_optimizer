package cards

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

// Tests that Outed doesn't apply the marked-defender bonus when OpponentMarked is false.
func TestOuted_NoMarkUnbuffed(t *testing.T) {
	self := &sim.CardState{Card: OutedRed{}}
	s := sim.TurnState{}
	(OutedRed{}).Play(&s, s.Logger(), self)
	if self.BonusAttack != 0 {
		t.Errorf("BonusAttack = %d, want 0 (mark off)", self.BonusAttack)
	}
}

// Tests that Outed self-buffs +1{p} when the opposing hero is marked at Play time.
func TestOuted_MarkedDefenderAddsOne(t *testing.T) {
	self := &sim.CardState{Card: OutedRed{}}
	s := sim.TurnState{OpponentMarked: true}
	(OutedRed{}).Play(&s, s.Logger(), self)
	if s.Value != 4 {
		t.Errorf("Play() Value = %d, want 4 (3 printed + 1 marked-defender)", s.Value)
	}
	if self.BonusAttack != 1 {
		t.Errorf("BonusAttack = %d, want 1", self.BonusAttack)
	}
}
