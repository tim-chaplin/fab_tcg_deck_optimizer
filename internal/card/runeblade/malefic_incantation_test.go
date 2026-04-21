package runeblade

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

// TestMaleficIncantation_PlayCreditsNMinusOneFlat: Play flips AuraCreated and returns n-1 flat
// damage for the future-turn verse-counter ticks that aren't separately modelled; the first
// tick's rune is deferred to PlayNextTurn.
func TestMaleficIncantation_PlayCreditsNMinusOneFlat(t *testing.T) {
	cases := []struct {
		c    card.Card
		want int
	}{
		{MaleficIncantationRed{}, 2},
		{MaleficIncantationYellow{}, 1},
		{MaleficIncantationBlue{}, 0},
	}
	for _, tc := range cases {
		var s card.TurnState
		if got := tc.c.Play(&s, &card.CardState{}); got != tc.want {
			t.Errorf("%s: Play() = %d, want %d", tc.c.Name(), got, tc.want)
		}
		if !s.AuraCreated {
			t.Errorf("%s: AuraCreated should be set", tc.c.Name())
		}
		if s.Runechants != 0 {
			t.Errorf("%s: Runechants = %d, want 0 (first tick deferred)", tc.c.Name(), s.Runechants)
		}
	}
}

// TestMaleficIncantation_PlayNextTurnCreatesOneRuneAndGraveyardsSelf: at next turn's upkeep,
// the first verse counter ticks — 1 live Runechant appears on the new state and Malefic heads
// to the graveyard. Returned Damage=1 credits the token toward next turn's Value.
func TestMaleficIncantation_PlayNextTurnCreatesOneRuneAndGraveyardsSelf(t *testing.T) {
	cases := []card.Card{
		MaleficIncantationRed{},
		MaleficIncantationYellow{},
		MaleficIncantationBlue{},
	}
	for _, c := range cases {
		var s card.TurnState
		dp := c.(card.DelayedPlay)
		r := dp.PlayNextTurn(&s)
		if r.Damage != 1 {
			t.Errorf("%s: Damage = %d, want 1", c.Name(), r.Damage)
		}
		if s.Runechants != 1 {
			t.Errorf("%s: Runechants = %d, want 1", c.Name(), s.Runechants)
		}
		if len(s.Graveyard) != 1 || s.Graveyard[0].ID() != c.ID() {
			t.Errorf("%s: Graveyard = %v, want [self]", c.Name(), s.Graveyard)
		}
	}
}
