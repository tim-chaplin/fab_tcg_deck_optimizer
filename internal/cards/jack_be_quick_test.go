package cards

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

// Tests that with no Nimblism in the graveyard the rider stays off.
func TestJackBeQuick_NoNimblismRiderOff(t *testing.T) {
	self := &sim.CardState{Card: JackBeQuickRed{}}
	s := sim.NewTurnStateFromCards(nil, nil)
	sim.ResolveChainStep(s, s.Logger(), self)
	if self.GrantedGoAgain {
		t.Errorf("GrantedGoAgain = true with empty graveyard, want false")
	}
	if self.BonusAttack != 0 {
		t.Errorf("BonusAttack = %d with empty graveyard, want 0", self.BonusAttack)
	}
}

// Tests that a Nimblism in the graveyard lets Jack Be Quick banish for the +1{p} /
// go-again rider.
func TestJackBeQuick_BanishesNimblismForBonus(t *testing.T) {
	s := sim.NewTurnStateFromCards(nil, []sim.Card{NimblismRed{}})
	self := &sim.CardState{Card: JackBeQuickRed{}}
	sim.ResolveChainStep(s, s.Logger(), self)
	if !self.GrantedGoAgain {
		t.Errorf("GrantedGoAgain = false after banish, want true")
	}
	if self.BonusAttack != 1 {
		t.Errorf("BonusAttack = %d after banish, want 1", self.BonusAttack)
	}
}
