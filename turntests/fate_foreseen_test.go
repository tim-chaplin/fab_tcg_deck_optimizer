package turntests

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card/cards"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/gameengine"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
)

// Tests that Fate Foreseen blocks for printed defense.
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
		top := testutils.FakeRedAction().WithName("top")
		ge := &gameengine.GameEngine{GameState: gameengine.GameStateBuilder().SetCards([]card.Card{top}).Build()}
		ge.SetIncomingPhysicalDamage(10)
		ge.ResolveAttackStep(ge.Logger(), &card.CardState{Card: tc.c})
		if ge.Value() != tc.block {
			t.Errorf("%s: Play(IncomingPhysicalDamage=10) Value = %d, want %d (block only)",
				tc.c.Name(), ge.Value(), tc.block)
		}
	}
}
