package turntests

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/cards"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/gameengine"
)

// Tests that Block flips +1{d} when at least one other plain blocker shares the slot.
func TestRightBehindYou_BlockTogetherFiresBonus(t *testing.T) {
	cases := []card.Card{
		cards.RightBehindYouRed{},
		cards.RightBehindYouYellow{},
		cards.RightBehindYouBlue{},
	}
	for _, c := range cases {
		blocker, ok := c.(card.Blocker)
		if !ok {
			t.Errorf("%s: missing card.Blocker hook", c.Name())
			continue
		}
		s := &gameengine.GameEngine{GameState: gameengine.GameStateBuilder().SetDefenders([]card.Card{c, testutils.GenericAttack(0, 1)}).Build()}
		self := &card.CardState{Card: c}
		blocker.Block(s, s.Logger(), self)
		if self.BonusDefense != 1 {
			t.Errorf("%s: BonusDefense = %d, want 1 (defending together)", c.Name(), self.BonusDefense)
		}
	}
}

// Tests that Block leaves BonusDefense untouched when this is the only plain blocker.
func TestRightBehindYou_BlockAloneNoBonus(t *testing.T) {
	c := cards.RightBehindYouRed{}
	s := &gameengine.GameEngine{GameState: gameengine.GameStateBuilder().SetDefenders([]card.Card{c}).Build()}
	self := &card.CardState{Card: c}
	c.Block(s, s.Logger(), self)
	if self.BonusDefense != 0 {
		t.Errorf("BonusDefense = %d, want 0 (alone)", self.BonusDefense)
	}
}
