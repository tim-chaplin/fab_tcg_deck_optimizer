package cards

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

// Tests that with no Nimblism in the graveyard the rider stays off.
func TestNimbleStrike_NoNimblismRiderOff(t *testing.T) {
	for _, c := range []sim.Card{NimbleStrikeRed{}, NimbleStrikeYellow{}, NimbleStrikeBlue{}} {
		self := &sim.CardState{Card: c}
		c.Play(sim.NewTurnState(nil, nil), self)
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
	for _, c := range []sim.Card{NimbleStrikeRed{}, NimbleStrikeYellow{}, NimbleStrikeBlue{}} {
		s := sim.NewTurnState(nil, []sim.Card{NimblismRed{}})
		self := &sim.CardState{Card: c}
		c.Play(s, self)
		if !self.GrantedGoAgain {
			t.Errorf("%s [%d{p}]: GrantedGoAgain = false after banish, want true", c.Name(), c.Pitch())
		}
		if self.BonusAttack != 1 {
			t.Errorf("%s [%d{p}]: BonusAttack = %d after banish, want 1", c.Name(), c.Pitch(), self.BonusAttack)
		}
	}
}
