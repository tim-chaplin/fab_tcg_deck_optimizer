package generic

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

// TestLifegainPerVariant guards against a regression where all three colour variants of a
// lifegain card credit the Red value. Printed gain is Red 3, Yellow 2, Blue 1 — a straight 1-to-1
// mapping to the Play return per variant.
func TestLifegainPerVariant(t *testing.T) {
	cases := []struct {
		name string
		card card.Card
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
		var s card.TurnState
		if got := tc.card.Play(&s); got != tc.want {
			t.Errorf("%s: Play() = %d, want %d", tc.name, got, tc.want)
		}
	}
}
