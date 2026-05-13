package turntests

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/cards"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/gameengine"
)

func TestShrillOfSkullform_BaseDamage(t *testing.T) {
	// Without any auras played this turn, Shrill returns its printed power.
	cases := []struct {
		c    card.Card
		want int
	}{
		{cards.ShrillOfSkullformRed{}, 4},
		{cards.ShrillOfSkullformYellow{}, 3},
		{cards.ShrillOfSkullformBlue{}, 2},
	}
	for _, tc := range cases {
		s := gameengine.NewFromState(nil)
		s.ResolveChainStep(s.Logger(), &card.CardState{Card: tc.c})
		got := s.Value()
		if got != tc.want {
			t.Errorf("%s: Play() = %d, want %d", tc.c.Name(), got, tc.want)
		}
	}
}

func TestShrillOfSkullform_AuraBonus(t *testing.T) {
	// With an aura in CardsPlayed, Shrill gets +3 power.
	cases := []struct {
		c    card.Card
		want int
	}{
		{cards.ShrillOfSkullformRed{}, 7},
		{cards.ShrillOfSkullformYellow{}, 6},
		{cards.ShrillOfSkullformBlue{}, 5},
	}
	for _, tc := range cases {
		s := gameengine.NewFromSpec(gameengine.Spec{CardsPlayed: []card.Card{testutils.Aura{}}, AuraCreated: true})
		s.ResolveChainStep(s.Logger(), &card.CardState{Card: tc.c})
		got := s.Value()
		if got != tc.want {
			t.Errorf("%s with aura: Play() = %d, want %d", tc.c.Name(), got, tc.want)
		}
	}
}
