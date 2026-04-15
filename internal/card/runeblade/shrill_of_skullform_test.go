package runeblade

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

// stubAura is a minimal Card with the "Aura" type for tests.
type stubAura struct{}

func (stubAura) Name() string                { return "StubAura" }
func (stubAura) Cost() int                   { return 0 }
func (stubAura) Pitch() int                  { return 0 }
func (stubAura) Attack() int                 { return 0 }
func (stubAura) Defense() int                { return 0 }
func (stubAura) Types() map[string]bool      { return map[string]bool{"Aura": true} }
func (stubAura) Play(*card.TurnState) int    { return 0 }

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
		got := tc.c.Play(&s)
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
		got := tc.c.Play(&s)
		if got != tc.want {
			t.Errorf("%s with aura: Play() = %d, want %d", tc.c.Name(), got, tc.want)
		}
	}
}
