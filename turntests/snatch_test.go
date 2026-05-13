package turntests

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/cards"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/gameengine"
)

// Tests that the on-hit DrawOne fires on a likely-hit attack, popping the deck top into
// hand.
func TestSnatch_LikelyHitFiresDrawOne(t *testing.T) {
	top := testutils.GenericAttack(0, 3)
	ge := &gameengine.GameEngine{GameState: gameengine.GameStateBuilder().SetCards([]card.Card{top}).Build()}
	c := cards.SnatchRed{}
	cs := &card.CardState{Card: c}
	ge.ResolveChainStep(ge.Logger(), cs)
	testutils.FireOnHitIfLikely(ge, ge.Logger(), cs)
	if got := ge.Value(); got != 4 {
		t.Errorf("Red: Play() = %d, want 4", got)
	}
	if h := ge.Hand(); len(h) != 1 || h[0] != top {
		t.Errorf("Hand = %v, want [top-of-deck]", h)
	}
	if d := ge.Deck(); d.Size() != 0 {
		t.Errorf("Deck size = %d, want 0 (top consumed)", d.Size())
	}
}

// Tests that the on-hit DrawOne doesn't fire on blockable variants.
func TestSnatch_BlockableSuppressesDraw(t *testing.T) {
	cases := []struct {
		c    card.Card
		want int
	}{
		{cards.SnatchYellow{}, 3},
		{cards.SnatchBlue{}, 2},
	}
	for _, tc := range cases {
		top := testutils.GenericAttack(0, 3)
		ge := &gameengine.GameEngine{GameState: gameengine.GameStateBuilder().SetCards([]card.Card{top}).Build()}
		ge.ResolveChainStep(ge.Logger(), &card.CardState{Card: tc.c})
		if got := ge.Value(); got != tc.want {
			t.Errorf("%s: Play() = %d, want %d (blockable, no draw)", tc.c.Name(), got, tc.want)
		}
		if h := ge.Hand(); len(h) != 0 {
			t.Errorf("%s: Hand = %v, want empty (no draw fired)", tc.c.Name(), h)
		}
		if d := ge.Deck(); d.Size() != 1 {
			t.Errorf("%s: Deck size = %d, want 1 (top preserved)", tc.c.Name(), d.Size())
		}
	}
}
