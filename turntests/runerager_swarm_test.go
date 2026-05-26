package turntests

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card/cards"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/gameengine"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
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
		ge := gameengine.New()
		pc := &card.CardState{Card: tc.c}
		ge.ResolveAttackStep(ge.Logger(), pc)
		if got := ge.Value(); got != tc.want {
			t.Errorf("%s: Play() = %d, want %d", tc.c.Name(), got, tc.want)
		}
		if pc.GrantedGoAgain {
			t.Errorf("%s: GrantedGoAgain should stay false with no aura", tc.c.Name())
		}
	}
}

func TestRuneragerSwarm_AuraPlayedGrantsGoAgain(t *testing.T) {
	// An aura in CardsPlayed satisfies the "played an aura this turn" condition.
	for _, c := range []card.Card{cards.RuneragerSwarmRed{}, cards.RuneragerSwarmYellow{}, cards.RuneragerSwarmBlue{}} {
		ge := &gameengine.GameEngine{GameState: gameengine.GameStateBuilder().
			SetCardsPlayed([]card.Card{testutils.FakeRedAura()}).
			SetAuraCreated(true).
			Build()}
		pc := &card.CardState{Card: c}
		ge.ResolveAttackStep(ge.Logger(), pc)
		if !pc.GrantedGoAgain {
			t.Errorf("%s: GrantedGoAgain should be set when an aura has been played", c.Name())
		}
	}
}

func TestRuneragerSwarm_AuraCreatedGrantsGoAgain(t *testing.T) {
	// TurnState.AuraCreated (e.g. from a runechant-creating effect earlier in the attack turn) also
	// satisfies the condition.
	for _, c := range []card.Card{cards.RuneragerSwarmRed{}, cards.RuneragerSwarmYellow{}, cards.RuneragerSwarmBlue{}} {
		ge := &gameengine.GameEngine{GameState: gameengine.GameStateBuilder().SetAuraCreated(true).Build()}
		pc := &card.CardState{Card: c}
		ge.ResolveAttackStep(ge.Logger(), pc)
		if !pc.GrantedGoAgain {
			t.Errorf("%s: GrantedGoAgain should be set when AuraCreated is true", c.Name())
		}
	}
}
