package turntests

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/v2/card/cards"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/gameengine"
)

// Tests that Fate Foreseen blocks for printed defense and emits an Opt 1 log entry.
func TestFateForeseen_BlocksAndCallsOpt1(t *testing.T) {
	cases := []struct {
		c     card.Card
		block int
	}{
		{cards.FateForeseenRed{}, 4},
		{cards.FateForeseenYellow{}, 3},
		{cards.FateForeseenBlue{}, 2},
	}

	for _, tc := range cases {
		top := testutils.NewStubCard("top")
		ge := &gameengine.GameEngine{GameState: gameengine.GameStateBuilder().SetCards([]card.Card{top}).Build()}
		ge.SetIncomingDamage(10)
		ge.ResolveChainStep(ge.Logger(), &card.CardState{Card: tc.c})
		if ge.Value() != tc.block {
			t.Errorf("%s: Play(IncomingDamage=10) Value = %d, want %d (block only)",
				tc.c.Name(), ge.Value(), tc.block)
		}
		if len(ge.LogEntries()) != 2 {
			t.Errorf("%s: Log len = %d, want 2 (Opted... + chain step)", tc.c.Name(), len(ge.LogEntries()))
			continue
		}
		// Play emits the Opted line during the DR resolution; the chain step is
		// auto-appended after Play returns, so the Opted entry lands first.
		want := "Opted [top], put [top] on top, put [] on bottom"
		if got := ge.LogEntries()[0].Text; got != want {
			t.Errorf("%s: Opt log entry = %q, want %q", tc.c.Name(), got, want)
		}
	}
}
