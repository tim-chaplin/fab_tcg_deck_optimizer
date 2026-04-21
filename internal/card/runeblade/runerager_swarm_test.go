package runeblade

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

func TestRuneragerSwarm_NoAuraNoGoAgain(t *testing.T) {
	// No aura played/created this turn → returns base power and does NOT grant self go-again.
	cases := []struct {
		c    card.Card
		want int
	}{
		{RuneragerSwarmRed{}, 3},
		{RuneragerSwarmYellow{}, 2},
		{RuneragerSwarmBlue{}, 1},
	}
	for _, tc := range cases {
		s := card.TurnState{}
		if got := tc.c.Play(&s); got != tc.want {
			t.Errorf("%s: Play() = %d, want %d", tc.c.Name(), got, tc.want)
		}
		if s.SelfGoAgain {
			t.Errorf("%s: SelfGoAgain should stay false with no aura", tc.c.Name())
		}
	}
}

func TestRuneragerSwarm_AuraPlayedGrantsGoAgain(t *testing.T) {
	// An aura in CardsPlayed satisfies the "played an aura this turn" condition.
	for _, c := range []card.Card{RuneragerSwarmRed{}, RuneragerSwarmYellow{}, RuneragerSwarmBlue{}} {
		s := card.TurnState{CardsPlayed: []card.Card{stubAura{}}}
		c.Play(&s)
		if !s.SelfGoAgain {
			t.Errorf("%s: SelfGoAgain should be set when an aura has been played", c.Name())
		}
	}
}

func TestRuneragerSwarm_AuraCreatedGrantsGoAgain(t *testing.T) {
	// TurnState.AuraCreated (e.g. from a runechant-creating effect earlier in the chain) also
	// satisfies the condition.
	for _, c := range []card.Card{RuneragerSwarmRed{}, RuneragerSwarmYellow{}, RuneragerSwarmBlue{}} {
		s := card.TurnState{AuraCreated: true}
		c.Play(&s)
		if !s.SelfGoAgain {
			t.Errorf("%s: SelfGoAgain should be set when AuraCreated is true", c.Name())
		}
	}
}
