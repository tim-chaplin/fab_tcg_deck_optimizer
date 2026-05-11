package cards

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
)

// Tests that Regurgitating Slog with no Sloggism in the graveyard does not gain Dominate.
func TestRegurgitatingSlog_NoSloggismNoDominate(t *testing.T) {
	for _, c := range []card.Card{RegurgitatingSlogRed{}, RegurgitatingSlogYellow{}, RegurgitatingSlogBlue{}} {
		s := sim.NewTurnStateFromCards(nil, nil)
		self := &card.CardState{Card: c}
		sim.ResolveChainStep(s, s.Logger(), self)
		if self.GrantedDominate {
			t.Errorf("%s [%d{p}]: GrantedDominate = true with no Sloggism, want false", c.Name(), c.Pitch())
		}
	}
}

// Tests that with a Sloggism in the graveyard, Regurgitating Slog banishes it and gains
// Dominate.
func TestRegurgitatingSlog_BanishesSloggismForDominate(t *testing.T) {
	for _, c := range []card.Card{RegurgitatingSlogRed{}, RegurgitatingSlogYellow{}, RegurgitatingSlogBlue{}} {
		s := sim.NewTurnStateFromCards(nil, []card.Card{SloggismRed{}})
		self := &card.CardState{Card: c}
		sim.ResolveChainStep(s, s.Logger(), self)
		if !self.GrantedDominate {
			t.Errorf("%s [%d{p}]: GrantedDominate = false after banish, want true", c.Name(), c.Pitch())
		}
	}
}
