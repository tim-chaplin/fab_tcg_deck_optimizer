package turntests

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card/cards"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/gameengine"
)

// Tests that with no Nimblism in the graveyard the rider stays off.
func TestNimbleStrike_NoNimblismRiderOff(t *testing.T) {
	for _, c := range []card.Card{cards.NimbleStrikeRed{}, cards.NimbleStrikeYellow{}, cards.NimbleStrikeBlue{}} {
		pc := &card.CardState{Card: c}
		ge := gameengine.New()
		ge.ResolveAttackStep(ge.Logger(), pc)
		if pc.GrantedGoAgain {
			t.Errorf("%s [%d{p}]: GrantedGoAgain = true with empty graveyard, want false", c.Name(), c.Pitch())
		}
		if pc.BonusAttack != 0 {
			t.Errorf("%s [%d{p}]: BonusAttack = %d with empty graveyard, want 0", c.Name(), c.Pitch(), pc.BonusAttack)
		}
	}
}

// Tests that a Nimblism in the graveyard lets Nimble Strike banish for the +1{p} /
// go-again rider.
func TestNimbleStrike_BanishesNimblismForBonus(t *testing.T) {
	for _, c := range []card.Card{cards.NimbleStrikeRed{}, cards.NimbleStrikeYellow{}, cards.NimbleStrikeBlue{}} {
		ge := &gameengine.GameEngine{GameState: gameengine.GameStateBuilder().SetGraveyard([]card.Card{cards.NimblismRed{}}).Build()}
		pc := &card.CardState{Card: c}
		ge.ResolveAttackStep(ge.Logger(), pc)
		if !pc.GrantedGoAgain {
			t.Errorf("%s [%d{p}]: GrantedGoAgain = false after banish, want true", c.Name(), c.Pitch())
		}
		if pc.BonusAttack != 1 {
			t.Errorf("%s [%d{p}]: BonusAttack = %d after banish, want 1", c.Name(), c.Pitch(), pc.BonusAttack)
		}
	}
}
