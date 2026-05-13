package turntests

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/cards"

	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/gameengine"
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
		{"HealingBalmRed", cards.HealingBalmRed{}, 3},
		{"HealingBalmYellow", cards.HealingBalmYellow{}, 2},
		{"HealingBalmBlue", cards.HealingBalmBlue{}, 1},
		{"SunKissRed", cards.SunKissRed{}, 3},
		{"SunKissYellow", cards.SunKissYellow{}, 2},
		{"SunKissBlue", cards.SunKissBlue{}, 1},
		{"FiddlersGreenRed", cards.FiddlersGreenRed{}, 3},
		{"FiddlersGreenYellow", cards.FiddlersGreenYellow{}, 2},
		{"FiddlersGreenBlue", cards.FiddlersGreenBlue{}, 1},
	}
	for _, tc := range cases {
		s := gameengine.New()
		s.ResolveChainStep(s.Logger(), &card.CardState{Card: tc.card})
		if got := s.Value(); got != tc.want {
			t.Errorf("%s: Play() = %d, want %d", tc.name, got, tc.want)
		}
	}
}
