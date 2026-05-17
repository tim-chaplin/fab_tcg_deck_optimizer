package turntests

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/v2/card/cards"

	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/gameengine"
)

// Tests that Regurgitating Slog with no Sloggism in the graveyard does not gain Dominate.
func TestRegurgitatingSlog_NoSloggismNoDominate(t *testing.T) {
	for _, c := range []card.Card{cards.RegurgitatingSlogRed{}, cards.RegurgitatingSlogYellow{}, cards.RegurgitatingSlogBlue{}} {
		ge := gameengine.New()
		self := &card.CardState{Card: c}
		ge.ResolveChainStep(ge.Logger(), self)
		if self.GrantedDominate {
			t.Errorf("%s [%d{p}]: GrantedDominate = true with no Sloggism, want false", c.Name(), c.Pitch())
		}
	}
}

// Tests that with a Sloggism in the graveyard, Regurgitating Slog banishes it and gains
// Dominate.
func TestRegurgitatingSlog_BanishesSloggismForDominate(t *testing.T) {
	for _, c := range []card.Card{cards.RegurgitatingSlogRed{}, cards.RegurgitatingSlogYellow{}, cards.RegurgitatingSlogBlue{}} {
		ge := &gameengine.GameEngine{GameState: gameengine.GameStateBuilder().SetGraveyard([]card.Card{cards.SloggismRed{}}).Build()}
		self := &card.CardState{Card: c}
		ge.ResolveChainStep(ge.Logger(), self)
		if !self.GrantedDominate {
			t.Errorf("%s [%d{p}]: GrantedDominate = false after banish, want true", c.Name(), c.Pitch())
		}
	}
}
