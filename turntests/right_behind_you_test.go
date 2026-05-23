package turntests

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card/cards"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/gameengine"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
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
		ge := &gameengine.GameEngine{GameState: gameengine.GameStateBuilder().SetDefenders([]card.Card{c, testutils.FakeRedAttack()}).Build()}
		pc := &card.CardState{Card: c}
		blocker.Block(ge, ge.Logger(), pc)
		if pc.BonusDefense != 1 {
			t.Errorf("%s: BonusDefense = %d, want 1 (defending together)", c.Name(), pc.BonusDefense)
		}
	}
}

// Tests that Block leaves BonusDefense untouched when this is the only plain blocker.
func TestRightBehindYou_BlockAloneNoBonus(t *testing.T) {
	c := cards.RightBehindYouRed{}
	ge := &gameengine.GameEngine{GameState: gameengine.GameStateBuilder().SetDefenders([]card.Card{c}).Build()}
	pc := &card.CardState{Card: c}
	c.Block(ge, ge.Logger(), pc)
	if pc.BonusDefense != 0 {
		t.Errorf("BonusDefense = %d, want 0 (alone)", pc.BonusDefense)
	}
}
