package turntests

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card/cards"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/gameengine"
)

// Tests that Block keeps printed Defense on mode 0.
func TestBrothersInArms_Mode0NoBonus(t *testing.T) {
	cases := []card.Card{
		cards.BrothersInArmsRed{},
		cards.BrothersInArmsYellow{},
		cards.BrothersInArmsBlue{},
	}
	for _, c := range cases {
		blocker := c.(card.Blocker)
		ge := gameengine.New()
		pc := &card.CardState{Card: c}
		blocker.Block(ge, ge.Logger(), pc)
		if pc.BonusDefense != 0 {
			t.Errorf("%s: mode 0 BonusDefense = %d, want 0", c.Name(), pc.BonusDefense)
		}
	}
}

// Tests that Block flips +2 BonusDefense on mode 1.
func TestBrothersInArms_Mode1FiresBonus(t *testing.T) {
	cases := []card.Card{
		cards.BrothersInArmsRed{},
		cards.BrothersInArmsYellow{},
		cards.BrothersInArmsBlue{},
	}
	for _, c := range cases {
		blocker := c.(card.Blocker)
		ge := gameengine.New()
		pc := &card.CardState{Card: c, Mode: 1}
		blocker.Block(ge, ge.Logger(), pc)
		if pc.BonusDefense != 2 {
			t.Errorf("%s: mode 1 BonusDefense = %d, want 2", c.Name(), pc.BonusDefense)
		}
	}
}

// Tests that BlockCost reports 0 on mode 0 and 1 on mode 1.
func TestBrothersInArms_BlockCostsTrackMode(t *testing.T) {
	c := cards.BrothersInArmsRed{}
	if got := c.BlockCost(0); got != 0 {
		t.Errorf("BlockCost(0) = %d, want 0", got)
	}
	if got := c.BlockCost(1); got != 1 {
		t.Errorf("BlockCost(1) = %d, want 1", got)
	}
}
