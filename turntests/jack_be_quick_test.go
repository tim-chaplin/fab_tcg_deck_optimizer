package turntests

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card/cards"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/gameengine"
)

// Tests that with no Nimblism in the graveyard the rider stays off.
func TestJackBeQuick_NoNimblismRiderOff(t *testing.T) {
	pc := &card.CardState{Card: cards.JackBeQuickRed{}}
	ge := gameengine.New()
	ge.ResolveAttackStep(ge.Logger(), pc)
	if pc.GrantedGoAgain {
		t.Errorf("GrantedGoAgain = true with empty graveyard, want false")
	}
	if pc.BonusAttack != 0 {
		t.Errorf("BonusAttack = %d with empty graveyard, want 0", pc.BonusAttack)
	}
}

// Tests that a Nimblism in the graveyard lets Jack Be Quick banish for the +1{p} /
// go-again rider.
func TestJackBeQuick_BanishesNimblismForBonus(t *testing.T) {
	ge := &gameengine.GameEngine{GameState: gameengine.GameStateBuilder().SetGraveyard([]card.Card{cards.NimblismRed{}}).Build()}
	pc := &card.CardState{Card: cards.JackBeQuickRed{}}
	ge.ResolveAttackStep(ge.Logger(), pc)
	if !pc.GrantedGoAgain {
		t.Errorf("GrantedGoAgain = false after banish, want true")
	}
	if pc.BonusAttack != 1 {
		t.Errorf("BonusAttack = %d after banish, want 1", pc.BonusAttack)
	}
}
