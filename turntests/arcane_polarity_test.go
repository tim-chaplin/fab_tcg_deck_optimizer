package turntests

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/cards"

	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/gameengine"
)

// Tests that with no arcane incoming the default branch credits 1{h}.
func TestArcanePolarity_NoArcaneIncomingCreditsOne(t *testing.T) {
	cases := []card.Card{
		cards.ArcanePolarityRed{},
		cards.ArcanePolarityYellow{},
		cards.ArcanePolarityBlue{},
	}
	for _, c := range cases {
		s := gameengine.New()
		s.ResolveChainStep(s.Logger(), &card.CardState{Card: c})
		if s.Value() != 1 {
			t.Errorf("%s: Value = %d, want 1", c.Name(), s.Value())
		}
	}
}

// Tests that ArcaneIncomingDamage > 0 swaps to the per-pitch alternate gain.
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
		s := gameengine.NewFromSpec(gameengine.Spec{ArcaneIncomingDamage: 1})
		s.ResolveChainStep(s.Logger(), &card.CardState{Card: tc.c})
		if s.Value() != tc.gain {
			t.Errorf("%s: Value = %d, want %d", tc.c.Name(), s.Value(), tc.gain)
		}
	}
}
