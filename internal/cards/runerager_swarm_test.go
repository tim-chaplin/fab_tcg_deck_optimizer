package cards

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
)

func TestRuneragerSwarm_NoAuraNoGoAgain(t *testing.T) {
	// No aura played/created this turn → returns base power and does NOT grant self go-again.
	cases := []struct {
		c    sim.Card
		want int
	}{
		{RuneragerSwarmRed{}, 3},
		{RuneragerSwarmYellow{}, 2},
		{RuneragerSwarmBlue{}, 1},
	}
	for _, tc := range cases {
		s := sim.TurnState{}
		self := &sim.CardState{Card: tc.c}
		tc.c.Play(&s, self)
		if got := s.Value; got != tc.want {
			t.Errorf("%s: Play() = %d, want %d", tc.c.Name(), got, tc.want)
		}
		if self.GrantedGoAgain {
			t.Errorf("%s: GrantedGoAgain should stay false with no aura", tc.c.Name())
		}
	}
}

func TestRuneragerSwarm_AuraPlayedGrantsGoAgain(t *testing.T) {
	// An aura in CardsPlayed satisfies the "played an aura this turn" condition.
	for _, c := range []sim.Card{RuneragerSwarmRed{}, RuneragerSwarmYellow{}, RuneragerSwarmBlue{}} {
		s := sim.TurnState{CardsPlayed: []sim.Card{testutils.Aura{}}}
		self := &sim.CardState{Card: c}
		c.Play(&s, self)
		if !self.GrantedGoAgain {
			t.Errorf("%s: GrantedGoAgain should be set when an aura has been played", c.Name())
		}
	}
}

func TestRuneragerSwarm_AuraCreatedGrantsGoAgain(t *testing.T) {
	// TurnState.AuraCreated (e.g. from a runechant-creating effect earlier in the chain) also
	// satisfies the condition.
	for _, c := range []sim.Card{RuneragerSwarmRed{}, RuneragerSwarmYellow{}, RuneragerSwarmBlue{}} {
		s := sim.TurnState{AuraCreated: true}
		self := &sim.CardState{Card: c}
		c.Play(&s, self)
		if !self.GrantedGoAgain {
			t.Errorf("%s: GrantedGoAgain should be set when AuraCreated is true", c.Name())
		}
	}
}
