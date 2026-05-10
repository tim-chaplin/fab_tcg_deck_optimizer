package cards

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
)

// Tests that Block flips +1{d} when this is the only plain blocker.
func TestBattlefrontBastion_BlockAloneFiresPrevention(t *testing.T) {
	cases := []sim.Card{
		BattlefrontBastionRed{},
		BattlefrontBastionYellow{},
		BattlefrontBastionBlue{},
	}
	for _, c := range cases {
		blocker, ok := c.(sim.Blocker)
		if !ok {
			t.Errorf("%s: missing sim.Blocker hook", c.Name())
			continue
		}
		s := sim.TurnState{Defenders: []sim.Card{c}}
		self := &sim.CardState{Card: c}
		blocker.Block(&s, s.Logger(), self)
		if self.BonusDefense != 1 {
			t.Errorf("%s: BonusDefense = %d, want 1 (alone)", c.Name(), self.BonusDefense)
		}
	}
}

// Tests that Block leaves BonusDefense untouched when a second plain blocker shares
// the defenders slot.
func TestBattlefrontBastion_BlockWithOtherPlainBlockerNoBonus(t *testing.T) {
	c := BattlefrontBastionRed{}
	s := sim.TurnState{Defenders: []sim.Card{c, testutils.GenericAttack(0, 1)}}
	self := &sim.CardState{Card: c}
	c.Block(&s, s.Logger(), self)
	if self.BonusDefense != 0 {
		t.Errorf("BonusDefense = %d, want 0 (another plain blocker)", self.BonusDefense)
	}
}
