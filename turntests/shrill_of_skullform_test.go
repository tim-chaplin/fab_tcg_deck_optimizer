package turntests

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/v2/card/cards"

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
		ge := gameengine.New()
		ge.ResolveChainStep(ge.Logger(), &card.CardState{Card: tc.c})
		got := ge.Value()
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
		ge := &gameengine.GameEngine{GameState: gameengine.GameStateBuilder().
			SetCardsPlayed([]card.Card{testutils.Aura{}}).
			SetAuraCreated(true).
			Build()}
		ge.ResolveChainStep(ge.Logger(), &card.CardState{Card: tc.c})
		got := ge.Value()
		if got != tc.want {
			t.Errorf("%s with aura: Play() = %d, want %d", tc.c.Name(), got, tc.want)
		}
	}
}
