package turntests

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/cards"

	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/gameengine"
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
		ge := gameengine.New()
		ge.ResolveChainStep(ge.Logger(), &card.CardState{Card: tc.c})
		if ge.RunechantCount() != tc.n {
			t.Errorf("%s: Runechants = %d, want %d", tc.c.Name(), ge.RunechantCount(), tc.n)
		}
		if !ge.AuraCreated() {
			t.Errorf("%s: AuraCreated = false, want true", tc.c.Name())
		}
		want := tc.c.Attack() + tc.n
		if ge.Value() != want {
			t.Errorf("%s: Value = %d, want %d (Attack %d + %d runechants)",
				tc.c.Name(), ge.Value(), want, tc.c.Attack(), tc.n)
		}
	}
}
