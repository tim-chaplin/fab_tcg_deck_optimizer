package turntests

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/v2/card/cards"

	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/gameengine"
)

// Tests that each printing prevents its full Defense() amount (4/3/2).
func TestPeaceOfMind_PreventsByPrinting(t *testing.T) {
	cases := []struct {
		card card.Card
		want int
	}{
		{cards.PeaceOfMindRed{}, 4},
		{cards.PeaceOfMindYellow{}, 3},
		{cards.PeaceOfMindBlue{}, 2},
	}
	for _, tc := range cases {
		ge := &gameengine.GameEngine{GameState: gameengine.GameStateBuilder().SetIncomingDamage(10).Build()}
		self := &card.CardState{Card: tc.card}
		ge.ResolveChainStep(ge.Logger(), self)
		if ge.Value() != tc.want {
			t.Errorf("%s: Value = %d, want %d", tc.card.Name(), ge.Value(), tc.want)
		}
	}
}

// Tests that each printing creates one Ponder token on play.
func TestPeaceOfMind_CreatesPonder(t *testing.T) {
	for _, c := range []card.Card{cards.PeaceOfMindRed{}, cards.PeaceOfMindYellow{}, cards.PeaceOfMindBlue{}} {
		ge := &gameengine.GameEngine{GameState: gameengine.GameStateBuilder().SetIncomingDamage(10).Build()}
		self := &card.CardState{Card: c}
		ge.ResolveChainStep(ge.Logger(), self)
		if got := ge.PonderCount(); got != 1 {
			t.Errorf("%s: Ponders = %d, want 1", c.Name(), got)
		}
		if !ge.AuraCreated() {
			t.Errorf("%s: AuraCreated = false, want true after creating a Ponder", c.Name())
		}
	}
}
