package turntests

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/cards"
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
)

// Tests that Outed doesn't apply the marked-defender bonus when OpponentMarked is false.
func TestOuted_NoMarkUnbuffed(t *testing.T) {
	self := &card.CardState{Card: cards.OutedRed{}}
	s := sim.TurnState{}
	sim.ResolveChainStep(&s, s.Logger(), self)
	if self.BonusAttack != 0 {
		t.Errorf("BonusAttack = %d, want 0 (mark off)", self.BonusAttack)
	}
}

// Tests that Outed self-buffs +1{p} when the opposing hero is marked at Play time.
func TestOuted_MarkedDefenderAddsOne(t *testing.T) {
	self := &card.CardState{Card: cards.OutedRed{}}
	s := sim.NewTurnStateFromSpec(sim.TurnStateSpec{OpponentMarked: true})
	sim.ResolveChainStep(&s, s.Logger(), self)
	if s.Value() != 4 {
		t.Errorf("Play() Value = %d, want 4 (3 printed + 1 marked-defender)", s.Value())
	}
	if self.BonusAttack != 1 {
		t.Errorf("BonusAttack = %d, want 1", self.BonusAttack)
	}
}
