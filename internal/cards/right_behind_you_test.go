package cards

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
)

// Tests that Block flips +1{d} when at least one other plain blocker shares the slot.
func TestRightBehindYou_BlockTogetherFiresBonus(t *testing.T) {
	cases := []sim.Card{
		RightBehindYouRed{},
		RightBehindYouYellow{},
		RightBehindYouBlue{},
	}
	for _, c := range cases {
		blocker, ok := c.(sim.Blocker)
		if !ok {
			t.Errorf("%s: missing sim.Blocker hook", c.Name())
			continue
		}
		s := sim.NewTurnStateFromSpec(sim.TurnStateSpec{Defenders: []sim.Card{c, testutils.GenericAttack(0, 1)}})
		self := &sim.CardState{Card: c}
		blocker.Block(&s, s.Logger(), self)
		if self.BonusDefense != 1 {
			t.Errorf("%s: BonusDefense = %d, want 1 (defending together)", c.Name(), self.BonusDefense)
		}
	}
}

// Tests that Block leaves BonusDefense untouched when this is the only plain blocker.
func TestRightBehindYou_BlockAloneNoBonus(t *testing.T) {
	c := RightBehindYouRed{}
	s := sim.NewTurnStateFromSpec(sim.TurnStateSpec{Defenders: []sim.Card{c}})
	self := &sim.CardState{Card: c}
	c.Block(&s, s.Logger(), self)
	if self.BonusDefense != 0 {
		t.Errorf("BonusDefense = %d, want 0 (alone)", self.BonusDefense)
	}
}
