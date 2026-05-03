package cards

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

// Tests that every variant returns +1 via sim.DefendsTogetherBonus.
func TestRightBehindYou_DefendsTogetherBonusReturnsOne(t *testing.T) {
	cases := []sim.Card{
		RightBehindYouRed{},
		RightBehindYouYellow{},
		RightBehindYouBlue{},
	}
	for _, c := range cases {
		dab, ok := c.(sim.DefendsTogetherBonus)
		if !ok {
			t.Errorf("%s: missing sim.DefendsTogetherBonus marker", c.Name())
			continue
		}
		if got := dab.DefendsTogetherBonus(); got != 1 {
			t.Errorf("%s: DefendsTogetherBonus() = %d, want 1", c.Name(), got)
		}
	}
}
