package cards

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
)

// Tests that with no Nimblism in the graveyard the rider stays off.
func TestNimbleStrike_NoNimblismRiderOff(t *testing.T) {
	for _, c := range []sim.Card{NimbleStrikeRed{}, NimbleStrikeYellow{}, NimbleStrikeBlue{}} {
		self := &card.CardState{Card: c}
		s := sim.NewTurnStateFromCards(nil, nil)
		sim.ResolveChainStep(s, s.Logger(), self)
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
		s := sim.NewTurnStateFromCards(nil, []sim.Card{NimblismRed{}})
		self := &card.CardState{Card: c}
		sim.ResolveChainStep(s, s.Logger(), self)
		if !self.GrantedGoAgain {
			t.Errorf("%s [%d{p}]: GrantedGoAgain = false after banish, want true", c.Name(), c.Pitch())
		}
		if self.BonusAttack != 1 {
			t.Errorf("%s [%d{p}]: BonusAttack = %d after banish, want 1", c.Name(), c.Pitch(), self.BonusAttack)
		}
	}
}
