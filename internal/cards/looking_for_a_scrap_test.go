package cards

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
)

// Tests that with no banishable graveyard target the rider stays off.
func TestLookingForAScrap_NoBanishableRiderOff(t *testing.T) {
	for _, c := range []card.Card{LookingForAScrapRed{}, LookingForAScrapYellow{}, LookingForAScrapBlue{}} {
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

// Tests that a 1-power graveyard card lets Looking for a Scrap banish for the +1{p} /
// go-again rider.
func TestLookingForAScrap_BanishesOnePowerForBonus(t *testing.T) {
	for _, c := range []card.Card{LookingForAScrapRed{}, LookingForAScrapYellow{}, LookingForAScrapBlue{}} {
		s := sim.NewTurnStateFromCards(nil, []card.Card{testutils.GenericAttack(0, 1)})
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
