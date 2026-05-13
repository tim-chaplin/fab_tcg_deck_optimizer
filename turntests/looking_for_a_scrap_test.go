package turntests

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/cards"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/gameengine"
)

// Tests that with no banishable graveyard target the rider stays off.
func TestLookingForAScrap_NoBanishableRiderOff(t *testing.T) {
	for _, c := range []card.Card{cards.LookingForAScrapRed{}, cards.LookingForAScrapYellow{}, cards.LookingForAScrapBlue{}} {
		self := &card.CardState{Card: c}
		s := &gameengine.GameEngine{GameState: gameengine.GameStateBuilder().Build()}
		s.ResolveChainStep(s.Logger(), self)
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
	for _, c := range []card.Card{cards.LookingForAScrapRed{}, cards.LookingForAScrapYellow{}, cards.LookingForAScrapBlue{}} {
		s := &gameengine.GameEngine{GameState: gameengine.GameStateBuilder().SetGraveyard([]card.Card{testutils.GenericAttack(0, 1)}).Build()}
		self := &card.CardState{Card: c}
		s.ResolveChainStep(s.Logger(), self)
		if !self.GrantedGoAgain {
			t.Errorf("%s [%d{p}]: GrantedGoAgain = false after banish, want true", c.Name(), c.Pitch())
		}
		if self.BonusAttack != 1 {
			t.Errorf("%s [%d{p}]: BonusAttack = %d after banish, want 1", c.Name(), c.Pitch(), self.BonusAttack)
		}
	}
}
