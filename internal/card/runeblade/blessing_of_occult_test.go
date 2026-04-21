package runeblade

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

// TestBlessingOfOccult_PlayIsAuraOnly: Play is a no-op beyond flipping AuraCreated. The
// Runechant payoff is deferred to PlayNextTurn so same-turn attacks don't consume tokens that
// only exist in next turn's upkeep.
func TestBlessingOfOccult_PlayIsAuraOnly(t *testing.T) {
	cases := []card.Card{BlessingOfOccultRed{}, BlessingOfOccultYellow{}, BlessingOfOccultBlue{}}
	for _, c := range cases {
		var s card.TurnState
		if got := c.Play(&s); got != 0 {
			t.Errorf("%s: Play() = %d, want 0 (Runechants deferred to PlayNextTurn)", c.Name(), got)
		}
		if s.Runechants != 0 {
			t.Errorf("%s: Runechants = %d, want 0", c.Name(), s.Runechants)
		}
		if !s.AuraCreated {
			t.Errorf("%s: AuraCreated should be set", c.Name())
		}
	}
}

// TestBlessingOfOccult_PlayNextTurnCreatesRunechants: the leave-arena trigger fires at the start
// of the next turn, destroys the aura, and creates N Runechant tokens (Red=3, Yellow=2, Blue=1).
func TestBlessingOfOccult_PlayNextTurnCreatesRunechants(t *testing.T) {
	cases := []struct {
		c card.DelayedPlay
		n int
	}{
		{BlessingOfOccultRed{}, 3},
		{BlessingOfOccultYellow{}, 2},
		{BlessingOfOccultBlue{}, 1},
	}
	for _, tc := range cases {
		var s card.TurnState
		got := tc.c.PlayNextTurn(&s)
		if got.Damage != tc.n {
			t.Errorf("%s: PlayNextTurn Damage = %d, want %d", tc.c.(card.Card).Name(), got.Damage, tc.n)
		}
		if s.Runechants != tc.n {
			t.Errorf("%s: Runechants = %d, want %d", tc.c.(card.Card).Name(), s.Runechants, tc.n)
		}
		if !s.SelfDestroyed {
			t.Errorf("%s: SelfDestroyed should be true (aura destroyed on leave)", tc.c.(card.Card).Name())
		}
	}
}
