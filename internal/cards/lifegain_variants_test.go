package cards

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

// TestLifegainPerVariant guards against a regression where all three colour variants of a
// lifegain card credit the Red value. Printed gain is Red 3, Yellow 2, Blue 1 — a straight 1-to-1
// mapping to the Play return per variant.
func TestLifegainPerVariant(t *testing.T) {
	cases := []struct {
		name string
		card sim.Card
		want int
	}{
		{"HealingBalmRed", HealingBalmRed{}, 3},
		{"HealingBalmYellow", HealingBalmYellow{}, 2},
		{"HealingBalmBlue", HealingBalmBlue{}, 1},
		{"SunKissRed", SunKissRed{}, 3},
		{"SunKissYellow", SunKissYellow{}, 2},
		{"SunKissBlue", SunKissBlue{}, 1},
		{"FiddlersGreenRed", FiddlersGreenRed{}, 3},
		{"FiddlersGreenYellow", FiddlersGreenYellow{}, 2},
		{"FiddlersGreenBlue", FiddlersGreenBlue{}, 1},
	}
	for _, tc := range cases {
		var s sim.TurnState
		tc.card.Play(&s, &sim.CardState{Card: tc.card})
		if got := s.Value; got != tc.want {
			t.Errorf("%s: Play() = %d, want %d", tc.name, got, tc.want)
		}
	}
}
