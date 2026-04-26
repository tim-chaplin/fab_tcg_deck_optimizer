package runeblade

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

func TestShrillOfSkullform_BaseDamage(t *testing.T) {
	// Without any auras played this turn, Shrill returns its printed power.
	cases := []struct {
		c    card.Card
		want int
	}{
		{ShrillOfSkullformRed{}, 4},
		{ShrillOfSkullformYellow{}, 3},
		{ShrillOfSkullformBlue{}, 2},
	}
	for _, tc := range cases {
		var s card.TurnState
		tc.c.Play(&s, &card.CardState{Card: tc.c})
		got := s.Value
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
		{ShrillOfSkullformRed{}, 7},
		{ShrillOfSkullformYellow{}, 6},
		{ShrillOfSkullformBlue{}, 5},
	}
	for _, tc := range cases {
		s := card.TurnState{CardsPlayed: []card.Card{stubAura{}}}
		tc.c.Play(&s, &card.CardState{Card: tc.c})
		got := s.Value
		if got != tc.want {
			t.Errorf("%s with aura: Play() = %d, want %d", tc.c.Name(), got, tc.want)
		}
	}
}
