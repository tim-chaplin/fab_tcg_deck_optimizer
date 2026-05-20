package turntests

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card/cards"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/deck"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
)

// siftFiller returns a deck of n RedAttack cards — generic filler that keeps Sift's
// mid-turn draws and the end-of-turn refill supplied without running dry.
func siftFiller(n int) []deck.Card {
	cs := make([]deck.Card, n)
	for i := range cs {
		cs[i] = testutils.RedAttack{}
	}
	return cs
}

// Tests that the hand cycle bottoms four spare hand cards: all four end up back in the
// deck. The twelve-card filler keeps them past the mid-turn and end-of-turn draws.
func TestSift_CyclesFourHandCardsToDeckBottom(t *testing.T) {
	for _, c := range []card.Card{cards.SiftRed{}, cards.SiftYellow{}, cards.SiftBlue{}} {
		spares := []card.Card{
			testutils.BluePitch{}, testutils.BluePitch{}, testutils.BluePitch{}, testutils.BluePitch{},
		}
		d := deck.New(testutils.Hero{Intel: 4}, nil, siftFiller(12))
		summary := sim.EvalOneTurnForTesting(d, nil, append([]card.Card{c}, spares...))
		if got := summary.State.Deck().NameCounts()[testutils.BluePitch{}.DisplayName()]; got != 4 {
			t.Errorf("%s: deck holds %d copies of the cycled card, want 4", c.Name(), got)
		}
	}
}

// Tests that the cycle is capped at hand size: with a single spare only one card is
// bottomed, not four.
func TestSift_CapsCycleAtHandSize(t *testing.T) {
	for _, c := range []card.Card{cards.SiftRed{}, cards.SiftYellow{}, cards.SiftBlue{}} {
		spare := testutils.BluePitch{}
		d := deck.New(testutils.Hero{Intel: 4}, nil, siftFiller(12))
		summary := sim.EvalOneTurnForTesting(d, nil, []card.Card{c, spare})
		if got := summary.State.Deck().NameCounts()[spare.DisplayName()]; got != 1 {
			t.Errorf("%s: deck holds %d copies of the cycled card, want 1", c.Name(), got)
		}
	}
}

// Tests that an empty hand drops the rider cleanly: Sift resolves with no spare to cycle.
func TestSift_EmptyHandResolves(t *testing.T) {
	for _, c := range []card.Card{cards.SiftRed{}, cards.SiftYellow{}, cards.SiftBlue{}} {
		d := deck.New(testutils.Hero{Intel: 4}, nil, siftFiller(12))
		summary := sim.EvalOneTurnForTesting(d, nil, []card.Card{c})
		if got := summary.State.Deck().NameCounts()[testutils.BluePitch{}.DisplayName()]; got != 0 {
			t.Errorf("%s: deck holds %d copies of BluePitch, want 0 (none were in hand)", c.Name(), got)
		}
	}
}
