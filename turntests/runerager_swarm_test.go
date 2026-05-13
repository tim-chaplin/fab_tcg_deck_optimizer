package turntests

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/cards"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/gameengine"
)

func TestRuneragerSwarm_NoAuraNoGoAgain(t *testing.T) {
	// No aura played/created this turn → returns base power and does NOT grant self go-again.
	cases := []struct {
		c    card.Card
		want int
	}{
		{cards.RuneragerSwarmRed{}, 3},
		{cards.RuneragerSwarmYellow{}, 2},
		{cards.RuneragerSwarmBlue{}, 1},
	}
	for _, tc := range cases {
		s := &gameengine.GameEngine{GameState: gameengine.NewState()}
		self := &card.CardState{Card: tc.c}
		s.ResolveChainStep(s.Logger(), self)
		if got := s.Value(); got != tc.want {
			t.Errorf("%s: Play() = %d, want %d", tc.c.Name(), got, tc.want)
		}
		if self.GrantedGoAgain {
			t.Errorf("%s: GrantedGoAgain should stay false with no aura", tc.c.Name())
		}
	}
}

func TestRuneragerSwarm_AuraPlayedGrantsGoAgain(t *testing.T) {
	// An aura in CardsPlayed satisfies the "played an aura this turn" condition.
	for _, c := range []card.Card{cards.RuneragerSwarmRed{}, cards.RuneragerSwarmYellow{}, cards.RuneragerSwarmBlue{}} {
		s := gameengine.NewFromSpec(gameengine.Spec{CardsPlayed: []card.Card{testutils.Aura{}}, AuraCreated: true})
		self := &card.CardState{Card: c}
		s.ResolveChainStep(s.Logger(), self)
		if !self.GrantedGoAgain {
			t.Errorf("%s: GrantedGoAgain should be set when an aura has been played", c.Name())
		}
	}
}

func TestRuneragerSwarm_AuraCreatedGrantsGoAgain(t *testing.T) {
	// TurnState.AuraCreated (e.g. from a runechant-creating effect earlier in the chain) also
	// satisfies the condition.
	for _, c := range []card.Card{cards.RuneragerSwarmRed{}, cards.RuneragerSwarmYellow{}, cards.RuneragerSwarmBlue{}} {
		s := gameengine.NewFromSpec(gameengine.Spec{AuraCreated: true})
		self := &card.CardState{Card: c}
		s.ResolveChainStep(s.Logger(), self)
		if !self.GrantedGoAgain {
			t.Errorf("%s: GrantedGoAgain should be set when AuraCreated is true", c.Name())
		}
	}
}
