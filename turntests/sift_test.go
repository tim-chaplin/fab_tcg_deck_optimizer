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
		cs[i] = testutils.FakeRedAttack()
	}
	return cs
}

// Tests that Sift bottoms four spare hand cards into the deck.
func TestSift_CyclesFourHandCardsToDeckBottom(t *testing.T) {
	for _, c := range []card.Card{cards.SiftRed{}, cards.SiftYellow{}, cards.SiftBlue{}} {
		spares := []card.Card{
			testutils.FakeBlueResource(), testutils.FakeBlueResource(), testutils.FakeBlueResource(), testutils.FakeBlueResource(),
		}
		d := deck.New(testutils.Hero{Intel: 4}, nil, siftFiller(12))
		summary := sim.EvalOneTurnForTesting(d, nil, append([]card.Card{c}, spares...))
		if got := summary.State.Deck().NameCounts()[testutils.FakeBlueResource().DisplayName()]; got != 4 {
			t.Errorf("%s: deck holds %d copies of the cycled card, want 4", c.Name(), got)
		}
	}
}

// Tests that Sift's cycle is capped at hand size when fewer than four spares are held.
func TestSift_CapsCycleAtHandSize(t *testing.T) {
	for _, c := range []card.Card{cards.SiftRed{}, cards.SiftYellow{}, cards.SiftBlue{}} {
		spare := testutils.FakeBlueResource()
		d := deck.New(testutils.Hero{Intel: 4}, nil, siftFiller(12))
		summary := sim.EvalOneTurnForTesting(d, nil, []card.Card{c, spare})
		if got := summary.State.Deck().NameCounts()[spare.DisplayName()]; got != 1 {
			t.Errorf("%s: deck holds %d copies of the cycled card, want 1", c.Name(), got)
		}
	}
}

// Tests that Sift resolves cleanly with an empty hand and nothing to cycle.
func TestSift_EmptyHandResolves(t *testing.T) {
	for _, c := range []card.Card{cards.SiftRed{}, cards.SiftYellow{}, cards.SiftBlue{}} {
		d := deck.New(testutils.Hero{Intel: 4}, nil, siftFiller(12))
		summary := sim.EvalOneTurnForTesting(d, nil, []card.Card{c})
		if got := summary.State.Deck().NameCounts()[testutils.FakeBlueResource().DisplayName()]; got != 0 {
			t.Errorf("%s: deck holds %d copies of BluePitch, want 0 (none were in hand)", c.Name(), got)
		}
	}
}
