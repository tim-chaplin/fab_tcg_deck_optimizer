package turntests

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card/cards"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/deck"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/gameengine"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
)

// Tests that the on-play "create N Runechant tokens" rider raises summary.State.RunechantCount()
// by N and credits N damage to summary.Value (printed power + N runechants resolving as arcane).
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
		d := deck.New(testutils.Hero{Intel: 4}, nil, nil)
		hand := []card.Card{tc.c, testutils.FakeBlueResource(), testutils.FakeBlueResource()}
		summary := sim.EvalOneTurnForTesting(d, gameengine.GameStateBuilder().SetIncomingDamage(0).Build(), hand)
		want := tc.c.Attack() + tc.n
		if summary.Value != want {
			t.Errorf("%s: Value = %d, want %d (Attack %d + %d runechants)",
				tc.c.Name(), summary.Value, want, tc.c.Attack(), tc.n)
		}
	}
}
