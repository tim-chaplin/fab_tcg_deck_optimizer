package turntests

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/cards"

	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/gameengine"
)

// Tests that with no Nimblism in the graveyard the rider stays off.
func TestJackBeQuick_NoNimblismRiderOff(t *testing.T) {
	self := &card.CardState{Card: cards.JackBeQuickRed{}}
	s := &gameengine.GameEngine{GameState: gameengine.GameStateBuilder().Build()}
	s.ResolveChainStep(s.Logger(), self)
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
	s := &gameengine.GameEngine{GameState: gameengine.GameStateBuilder().SetGraveyard([]card.Card{cards.NimblismRed{}}).Build()}
	self := &card.CardState{Card: cards.JackBeQuickRed{}}
	s.ResolveChainStep(s.Logger(), self)
	if !self.GrantedGoAgain {
		t.Errorf("GrantedGoAgain = false after banish, want true")
	}
	if self.BonusAttack != 1 {
		t.Errorf("BonusAttack = %d after banish, want 1", self.BonusAttack)
	}
}
