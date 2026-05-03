package cards

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

// Tests that every variant returns +1 via sim.DefendsAloneBonus.
func TestBattlefrontBastion_DefendsAloneBonusReturnsOne(t *testing.T) {
	cases := []sim.Card{
		BattlefrontBastionRed{},
		BattlefrontBastionYellow{},
		BattlefrontBastionBlue{},
	}
	for _, c := range cases {
		dab, ok := c.(sim.DefendsAloneBonus)
		if !ok {
			t.Errorf("%s: missing sim.DefendsAloneBonus marker", c.Name())
			continue
		}
		if got := dab.DefendsAloneBonus(); got != 1 {
			t.Errorf("%s: DefendsAloneBonus() = %d, want 1", c.Name(), got)
		}
	}
}
