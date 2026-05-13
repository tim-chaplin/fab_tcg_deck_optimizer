package turntests

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/cards"

	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/gameengine"
)

// Tests that with no Nimblism in the graveyard the rider stays off.
func TestNimbleStrike_NoNimblismRiderOff(t *testing.T) {
	for _, c := range []card.Card{cards.NimbleStrikeRed{}, cards.NimbleStrikeYellow{}, cards.NimbleStrikeBlue{}} {
		self := &card.CardState{Card: c}
		ge := gameengine.New()
		ge.ResolveChainStep(ge.Logger(), self)
		if self.GrantedGoAgain {
			t.Errorf("%s [%d{p}]: GrantedGoAgain = true with empty graveyard, want false", c.Name(), c.Pitch())
		}
		if self.BonusAttack != 0 {
			t.Errorf("%s [%d{p}]: BonusAttack = %d with empty graveyard, want 0", c.Name(), c.Pitch(), self.BonusAttack)
		}
	}
}

// Tests that a Nimblism in the graveyard lets Nimble Strike banish for the +1{p} /
// go-again rider.
func TestNimbleStrike_BanishesNimblismForBonus(t *testing.T) {
	for _, c := range []card.Card{cards.NimbleStrikeRed{}, cards.NimbleStrikeYellow{}, cards.NimbleStrikeBlue{}} {
		ge := &gameengine.GameEngine{GameState: gameengine.GameStateBuilder().SetGraveyard([]card.Card{cards.NimblismRed{}}).Build()}
		self := &card.CardState{Card: c}
		ge.ResolveChainStep(ge.Logger(), self)
		if !self.GrantedGoAgain {
			t.Errorf("%s [%d{p}]: GrantedGoAgain = false after banish, want true", c.Name(), c.Pitch())
		}
		if self.BonusAttack != 1 {
			t.Errorf("%s [%d{p}]: BonusAttack = %d after banish, want 1", c.Name(), c.Pitch(), self.BonusAttack)
		}
	}
}
