package turntests

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card/cards"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/gameengine"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
)

// Tests printed Cost() == 2 always and AlternativeCost reports (0, true) when a Held card
// is available, (0, false) when empty.
func TestRiseAbove_PrintedAndAlternativeCost(t *testing.T) {
	for _, c := range []card.Card{cards.RiseAboveRed{}, cards.RiseAboveYellow{}, cards.RiseAboveBlue{}} {
		if got := c.Cost(); got != 2 {
			t.Errorf("%s: Cost() = %d, want 2 (printed)", c.Name(), got)
		}
		ac, ok := c.(card.AlternativeCost)
		if !ok {
			t.Fatalf("%s: missing card.AlternativeCost", c.Name())
		}
		held := gameengine.New()
		held.SetHand([]card.Card{testutils.FakeRedAttack()})
		if cost, available := ac.AlternativeCost(held); cost != 0 || !available {
			t.Errorf("%s: AlternativeCost(hand=1) = (%d, %v), want (0, true)", c.Name(), cost, available)
		}
		empty := gameengine.New()
		if _, available := ac.AlternativeCost(empty); available {
			t.Errorf("%s: AlternativeCost(empty hand) reports available, want unavailable", c.Name())
		}
	}
}

// Tests that Play pops a hand card and prepends it to the deck only when PaidAlternativeCost
// is set on the CardState.
func TestRiseAbove_AltCostMovesHandCardToDeckTop(t *testing.T) {
	for _, c := range []card.Card{cards.RiseAboveRed{}, cards.RiseAboveYellow{}, cards.RiseAboveBlue{}} {
		spare := testutils.FakeRedAttack()
		ge := gameengine.New()
		ge.SetHand([]card.Card{spare})
		pc := &card.CardState{
			Card:      c,
			Ephemeral: card.Ephemeral{PaidAlternativeCost: true},
		}
		ge.ResolveAttackStep(ge.Logger(), pc)
		if len(ge.Hand()) != 0 {
			t.Errorf("%s: Hand = %d entries, want 0", c.Name(), len(ge.Hand()))
		}
		top, ok := ge.PeekDeck()
		if !ok || top.Name() != spare.Name() {
			t.Errorf("%s: deck top = %v, want %v", c.Name(), top, spare)
		}
	}
}

// Tests that the printed-cost branch (PaidAlternativeCost=false) doesn't touch the hand.
func TestRiseAbove_PrintedCostBranchSkipsAltCostSideEffect(t *testing.T) {
	for _, c := range []card.Card{cards.RiseAboveRed{}, cards.RiseAboveYellow{}, cards.RiseAboveBlue{}} {
		spare := testutils.FakeRedAttack()
		ge := gameengine.New()
		ge.SetHand([]card.Card{spare})
		pc := &card.CardState{Card: c}
		ge.ResolveAttackStep(ge.Logger(), pc)
		if len(ge.Hand()) != 1 {
			t.Errorf("%s: Hand size = %d, want 1 (printed-cost branch should not pop the hand card)", c.Name(), len(ge.Hand()))
		}
	}
}
