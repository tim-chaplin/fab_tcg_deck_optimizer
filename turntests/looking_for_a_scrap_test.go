package turntests

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card/cards"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/gameengine"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
)

// Tests that with no banishable graveyard target the rider stays off.
func TestLookingForAScrap_NoBanishableRiderOff(t *testing.T) {
	for _, c := range []card.Card{cards.LookingForAScrapRed{}, cards.LookingForAScrapYellow{}, cards.LookingForAScrapBlue{}} {
		pc := &card.CardState{Card: c}
		ge := gameengine.New()
		ge.ResolveChainStep(ge.Logger(), pc)
		if pc.GrantedGoAgain {
			t.Errorf("%s [%d{p}]: GrantedGoAgain = true with empty graveyard, want false", c.Name(), c.Pitch())
		}
		if pc.BonusAttack != 0 {
			t.Errorf("%s [%d{p}]: BonusAttack = %d with empty graveyard, want 0", c.Name(), c.Pitch(), pc.BonusAttack)
		}
	}
}

// Tests that a 1-power graveyard card lets Looking for a Scrap banish for the +1{p} /
// go-again rider.
func TestLookingForAScrap_BanishesOnePowerForBonus(t *testing.T) {
	for _, c := range []card.Card{cards.LookingForAScrapRed{}, cards.LookingForAScrapYellow{}, cards.LookingForAScrapBlue{}} {
		ge := &gameengine.GameEngine{GameState: gameengine.GameStateBuilder().SetGraveyard([]card.Card{testutils.GenericAttack(0, 1)}).Build()}
		pc := &card.CardState{Card: c}
		ge.ResolveChainStep(ge.Logger(), pc)
		if !pc.GrantedGoAgain {
			t.Errorf("%s [%d{p}]: GrantedGoAgain = false after banish, want true", c.Name(), c.Pitch())
		}
		if pc.BonusAttack != 1 {
			t.Errorf("%s [%d{p}]: BonusAttack = %d after banish, want 1", c.Name(), c.Pitch(), pc.BonusAttack)
		}
	}
}
