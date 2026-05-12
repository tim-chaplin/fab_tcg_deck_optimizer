package turntests

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/cards"
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
)

// Tests that the on-play "create N Runechant tokens" rider raises s.RunechantCount() by N,
// sets AuraCreated, and credits N damage to Value.
func TestRunechantOnPlay_CreatesNTokens(t *testing.T) {
	cases := []struct {
		c card.Card
		n int
	}{
		{cards.HocusPocusRed{}, 1},
		{cards.HocusPocusYellow{}, 1},
		{cards.HocusPocusBlue{}, 1},
		{cards.ReadTheRunesRed{}, 3},
		{cards.ReadTheRunesYellow{}, 2},
		{cards.ReadTheRunesBlue{}, 1},
		{cards.SpellbladeAssaultRed{}, 2},
		{cards.SpellbladeAssaultYellow{}, 2},
		{cards.SpellbladeAssaultBlue{}, 2},
		{cards.SpellbladeStrikeRed{}, 1},
		{cards.SpellbladeStrikeYellow{}, 1},
		{cards.SpellbladeStrikeBlue{}, 1},
	}
	for _, tc := range cases {
		var s sim.TurnState
		sim.ResolveChainStep(&s, s.Logger(), &card.CardState{Card: tc.c})
		if s.RunechantCount() != tc.n {
			t.Errorf("%s: Runechants = %d, want %d", tc.c.Name(), s.RunechantCount(), tc.n)
		}
		if !s.AuraCreated() {
			t.Errorf("%s: AuraCreated = false, want true", tc.c.Name())
		}
		want := tc.c.Attack() + tc.n
		if s.Value() != want {
			t.Errorf("%s: Value = %d, want %d (Attack %d + %d runechants)",
				tc.c.Name(), s.Value(), want, tc.c.Attack(), tc.n)
		}
	}
}
