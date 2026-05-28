package turntests

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card/cards"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/gameengine"
)

// Tests that with no arcane incoming the default branch credits 1{h}.
func TestArcanePolarity_NoArcaneIncomingCreditsOne(t *testing.T) {
	cases := []card.Card{
		cards.ArcanePolarityRed{},
		cards.ArcanePolarityYellow{},
		cards.ArcanePolarityBlue{},
	}
	for _, c := range cases {
		ge := gameengine.New()
		ge.ResolveAttackStep(ge.Logger(), &card.CardState{Card: c})
		if ge.Value() != 1 {
			t.Errorf("%s: Value = %d, want 1", c.Name(), ge.Value())
		}
	}
}

// Tests that IncomingArcaneDamage > 0 swaps to the per-pitch alternate gain.
func TestArcanePolarity_ArcaneIncomingCreditsLargeGain(t *testing.T) {
	cases := []struct {
		c    card.Card
		gain int
	}{
		{cards.ArcanePolarityRed{}, 4},
		{cards.ArcanePolarityYellow{}, 3},
		{cards.ArcanePolarityBlue{}, 2},
	}
	for _, tc := range cases {
		ge := &gameengine.GameEngine{GameState: gameengine.GameStateBuilder().SetIncomingArcaneDamage(1).Build()}
		ge.ResolveAttackStep(ge.Logger(), &card.CardState{Card: tc.c})
		if ge.Value() != tc.gain {
			t.Errorf("%s: Value = %d, want %d", tc.c.Name(), ge.Value(), tc.gain)
		}
	}
}
