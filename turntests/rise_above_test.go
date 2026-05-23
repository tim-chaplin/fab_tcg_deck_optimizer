package turntests

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card/cards"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/gameengine"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
)

// Tests that Cost is 0 when a hand card can pay the alt cost, else printed 2.
func TestRiseAbove_VariableCost(t *testing.T) {
	for _, c := range []card.Card{cards.RiseAboveRed{}, cards.RiseAboveYellow{}, cards.RiseAboveBlue{}} {
		held := gameengine.New()
		held.SetHand([]card.Card{testutils.FakeRedAttack()})
		if got := c.Cost(held); got != 0 {
			t.Errorf("%s: Cost(Hand) = %d, want 0", c.Name(), got)
		}
		empty := gameengine.New()
		if got := c.Cost(empty); got != 2 {
			t.Errorf("%s: Cost(empty) = %d, want 2", c.Name(), got)
		}
		vc, ok := c.(card.VariableCost)
		if !ok {
			t.Errorf("%s: missing card.VariableCost", c.Name())
			continue
		}
		if vc.MinCost() != 0 || vc.MaxCost() != 2 {
			t.Errorf("%s: bounds = [%d, %d], want [0, 2]", c.Name(), vc.MinCost(), vc.MaxCost())
		}
	}
}

// Tests that the alt cost pops a hand card and prepends it to the deck.
func TestRiseAbove_AltCostMovesHandCardToDeckTop(t *testing.T) {
	for _, c := range []card.Card{cards.RiseAboveRed{}, cards.RiseAboveYellow{}, cards.RiseAboveBlue{}} {
		spare := testutils.FakeRedAttack()
		ge := gameengine.New()
		ge.SetHand([]card.Card{spare})
		pc := &card.CardState{Card: c}
		ge.ResolveChainStep(ge.Logger(), pc)
		if len(ge.Hand()) != 0 {
			t.Errorf("%s: Hand = %d entries, want 0", c.Name(), len(ge.Hand()))
		}
		top, ok := ge.PeekDeck()
		if !ok || top.Name() != spare.Name() {
			t.Errorf("%s: deck top = %v, want %v", c.Name(), top, spare)
		}
	}
}

// Tests that an empty hand leaves the deck untouched.
func TestRiseAbove_EmptyHandLeavesDeckUntouched(t *testing.T) {
	for _, c := range []card.Card{cards.RiseAboveRed{}, cards.RiseAboveYellow{}, cards.RiseAboveBlue{}} {
		ge := gameengine.New()
		preSize := ge.Deck().Size()
		pc := &card.CardState{Card: c}
		ge.ResolveChainStep(ge.Logger(), pc)
		if got := ge.Deck().Size(); got != preSize {
			t.Errorf("%s: deck size = %d, want %d", c.Name(), got, preSize)
		}
	}
}
