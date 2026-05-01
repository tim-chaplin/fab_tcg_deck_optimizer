package cards

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

// Tests that the on-play "create N Runechant tokens" rider increments TurnState.Runechants
// by N, sets AuraCreated, and credits N damage to Value.
func TestRunechantOnPlay_CreatesNTokens(t *testing.T) {
	cases := []struct {
		c sim.Card
		n int
	}{
		{HocusPocusRed{}, 1},
		{HocusPocusYellow{}, 1},
		{HocusPocusBlue{}, 1},
		{ReadTheRunesRed{}, 3},
		{ReadTheRunesYellow{}, 2},
		{ReadTheRunesBlue{}, 1},
		{SpellbladeAssaultRed{}, 2},
		{SpellbladeAssaultYellow{}, 2},
		{SpellbladeAssaultBlue{}, 2},
		{SpellbladeStrikeRed{}, 1},
		{SpellbladeStrikeYellow{}, 1},
		{SpellbladeStrikeBlue{}, 1},
	}
	for _, tc := range cases {
		var s sim.TurnState
		tc.c.Play(&s, &sim.CardState{Card: tc.c})
		if s.Runechants != tc.n {
			t.Errorf("%s: Runechants = %d, want %d", tc.c.Name(), s.Runechants, tc.n)
		}
		if !s.AuraCreated {
			t.Errorf("%s: AuraCreated = false, want true", tc.c.Name())
		}
		want := tc.c.Attack() + tc.n
		if s.Value != want {
			t.Errorf("%s: Value = %d, want %d (Attack %d + %d runechants)",
				tc.c.Name(), s.Value, want, tc.c.Attack(), tc.n)
		}
	}
}
