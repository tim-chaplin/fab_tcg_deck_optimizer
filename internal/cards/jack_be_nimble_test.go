package cards

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
)

// Tests that with no Nimblism in the graveyard the rider stays off.
func TestJackBeNimble_NoNimblismRiderOff(t *testing.T) {
	self := &card.CardState{Card: JackBeNimbleRed{}}
	s := sim.NewTurnStateFromCards(nil, nil)
	sim.ResolveChainStep(s, s.Logger(), self)
	if self.GrantedGoAgain {
		t.Errorf("GrantedGoAgain = true with empty graveyard, want false")
	}
	if self.BonusAttack != 0 {
		t.Errorf("BonusAttack = %d with empty graveyard, want 0", self.BonusAttack)
	}
}

// Tests that a Nimblism in the graveyard lets Jack Be Nimble banish for the +1{p} /
// go-again rider.
func TestJackBeNimble_BanishesNimblismForBonus(t *testing.T) {
	s := sim.NewTurnStateFromCards(nil, []card.Card{NimblismRed{}})
	self := &card.CardState{Card: JackBeNimbleRed{}}
	sim.ResolveChainStep(s, s.Logger(), self)
	if !self.GrantedGoAgain {
		t.Errorf("GrantedGoAgain = false after banish, want true")
	}
	if self.BonusAttack != 1 {
		t.Errorf("BonusAttack = %d after banish, want 1", self.BonusAttack)
	}
}
