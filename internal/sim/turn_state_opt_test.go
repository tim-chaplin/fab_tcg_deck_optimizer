package sim

import (
	"strings"
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/deck"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/gameengine"
)

// withOptHero is a no-op stub kept for source compatibility — tests now set the hero on
// each engine directly. Retained so the per-test call shape doesn't need updating.
func withOptHero(t *testing.T, h FakeHero, fn func()) {
	t.Helper()
	currentOptTestHero = h
	defer func() { currentOptTestHero = FakeHero{} }()
	fn()
}

// currentOptTestHero is the hero installed on each engine built inside withOptHero.
var currentOptTestHero FakeHero

// newOptTestEngine builds a *gameengine.GameEngine with the deck seeded and the
// withOptHero-installed currentOptTestHero attached.
func newOptTestEngine(deckCards, graveyard []card.Card) *gameengine.GameEngine {
	ge := &gameengine.GameEngine{GameState: gameengine.GameStateBuilder().
		SetCards(deckCards).
		SetGraveyard(graveyard).
		Build()}
	ge.SetHero(currentOptTestHero)
	return ge
}

// Tests that Opt with the default passthrough handler keeps the deck order unchanged —
// every revealed card returns to the top in input order, none move to the bottom.
func TestTurnStateOpt_PassthroughKeepsDeckOrder(t *testing.T) {
	a := NewFakeCard("a")
	b := NewFakeCard("b")
	c := NewFakeCard("c")
	d := NewFakeCard("d")
	withOptHero(t, FakeHero{Intel: 4}, func() {
		s := newOptTestEngine([]card.Card{a, b, c, d}, nil)
		s.Opt(s.Logger(), 2)
		got := s.Deck()
		want := []card.Card{a, b, c, d}
		if !sameDeck(got, want) {
			t.Errorf("deck = %v, want %v (passthrough)", got, want)
		}
	})
}

// Tests that the handler can move cards to the bottom of the deck. Handler bottoms the
// first revealed card and keeps the second on top; the un-opted tail of the deck stays
// in place, and the bottomed card lands at the end.
func TestTurnStateOpt_BottomsHandlerSpecifiedCards(t *testing.T) {
	a := NewFakeCard("a")
	b := NewFakeCard("b")
	c := NewFakeCard("c")
	withOptHero(t, FakeHero{
		OptStrategy: func(cards []card.Card) (top, bottom []card.Card) {
			// Bottom cards[0] (a), keep cards[1] (b) on top.
			return []card.Card{cards[1]}, []card.Card{cards[0]}
		},
	}, func() {
		s := newOptTestEngine([]card.Card{a, b, c}, nil)
		s.Opt(s.Logger(), 2)
		got := s.Deck()
		// Handler saw [a, b]; returned top=[b], bottom=[a]. Deck becomes [b] + [c] + [a].
		want := []card.Card{b, c, a}
		if !sameDeck(got, want) {
			t.Errorf("deck = %v, want %v (handler bottoms a)", got, want)
		}
	})
}

// Tests that the handler can re-order cards on top in addition to bottoming.
func TestTurnStateOpt_HandlerReorderCanReverseTop(t *testing.T) {
	a := NewFakeCard("a")
	b := NewFakeCard("b")
	c := NewFakeCard("c")
	withOptHero(t, FakeHero{
		OptStrategy: func(cards []card.Card) (top, bottom []card.Card) {
			return []card.Card{cards[1], cards[0]}, nil
		},
	}, func() {
		s := newOptTestEngine([]card.Card{a, b, c}, nil)
		s.Opt(s.Logger(), 2)
		got := s.Deck()
		want := []card.Card{b, a, c}
		if !sameDeck(got, want) {
			t.Errorf("deck = %v, want %v (top reversed)", got, want)
		}
	})
}

// Tests that a request larger than the deck length clamps to whatever cards are there.
// The handler sees only the available cards and reshapes them; the empty-tail-deck
// remains empty.
func TestTurnStateOpt_ClampsNToDeckLength(t *testing.T) {
	a := NewFakeCard("a")
	b := NewFakeCard("b")
	withOptHero(t, FakeHero{
		OptStrategy: func(cards []card.Card) (top, bottom []card.Card) {
			if len(cards) != 2 {
				t.Errorf("handler saw %d cards, want 2 (clamp)", len(cards))
			}
			return cards, nil
		},
	}, func() {
		s := newOptTestEngine([]card.Card{a, b}, nil)
		s.Opt(s.Logger(), 5)
		got := s.Deck()
		want := []card.Card{a, b}
		if !sameDeck(got, want) {
			t.Errorf("deck = %v, want %v (clamped)", got, want)
		}
	})
}

// Tests that Opt on an empty deck never invokes the handler and is a safe no-op.
func TestTurnStateOpt_EmptyDeckSkipsHandler(t *testing.T) {
	called := false
	withOptHero(t, FakeHero{
		OptStrategy: func(cards []card.Card) (top, bottom []card.Card) {
			called = true
			return cards, nil
		},
	}, func() {
		s := newOptTestEngine(nil, nil)
		s.Opt(s.Logger(), 3)
		if called {
			t.Error("handler should not be called on empty deck")
		}
	})
}

// Tests that Opt(0) and Opt(-1) are no-ops aside from the cacheable flip — n <= 0 has no
// cards to reshape so the handler isn't invoked.
func TestTurnStateOpt_NonPositiveNSkipsHandler(t *testing.T) {
	for _, n := range []int{0, -1, -42} {
		called := false
		withOptHero(t, FakeHero{
			OptStrategy: func(cards []card.Card) (top, bottom []card.Card) {
				called = true
				return cards, nil
			},
		}, func() {
			s := newOptTestEngine([]card.Card{NewFakeCard("x")}, nil)
			s.Opt(s.Logger(), n)
			if called {
				t.Errorf("Opt(%d) called the handler, want skip", n)
			}
		})
	}
}

// Tests that Opt always flips IsCacheable to false, even on the no-op paths (n <= 0,
// empty deck) — the chain reading the deck implies an order dependency regardless of
// whether the handler ran.
func TestTurnStateOpt_AlwaysFlipsCacheable(t *testing.T) {
	cases := []struct {
		name string
		deck []card.Card
		n    int
	}{
		{"populated deck", []card.Card{NewFakeCard("x")}, 1},
		{"empty deck", nil, 3},
		{"zero n", []card.Card{NewFakeCard("x")}, 0},
	}
	withOptHero(t, FakeHero{}, func() {
		for _, tc := range cases {
			s := newOptTestEngine(tc.deck, nil)
			if !s.IsCacheable() {
				t.Fatalf("%s: pre IsCacheable should be true", tc.name)
			}
			s.Opt(s.Logger(), tc.n)
			if s.IsCacheable() {
				t.Errorf("%s: Opt(%d) should flip IsCacheable to false", tc.name, tc.n)
			}
		}
	})
}

// Tests that a handler returning fewer cards than received panics.
func TestTurnStateOpt_PanicsOnDroppedCard(t *testing.T) {
	withOptHero(t, FakeHero{
		OptStrategy: func(cards []card.Card) (top, bottom []card.Card) {
			return []card.Card{cards[0]}, nil // drops cards[1]
		},
	}, func() {
		s := newOptTestEngine([]card.Card{NewFakeCard("a"), NewFakeCard("b")}, nil)
		assertPanics(t, "dropped card", "Opt:", func() { s.Opt(s.Logger(), 2) })
	})
}

// Tests that a handler returning more cards than received panics.
func TestTurnStateOpt_PanicsOnExtraCard(t *testing.T) {
	a := NewFakeCard("a")
	b := NewFakeCard("b")
	foreign := NewFakeCard("foreign")
	withOptHero(t, FakeHero{
		OptStrategy: func(cards []card.Card) (top, bottom []card.Card) {
			return []card.Card{a, b, foreign}, nil
		},
	}, func() {
		s := newOptTestEngine([]card.Card{a, b}, nil)
		assertPanics(t, "extra card", "Opt:", func() { s.Opt(s.Logger(), 2) })
	})
}

// Tests that a handler substituting one input card for a non-input card panics — the
// length check passes but the multiset check trips on the foreign card.
func TestTurnStateOpt_PanicsOnSubstitutedCard(t *testing.T) {
	a := NewFakeCard("a")
	b := NewFakeCard("b")
	foreign := NewFakeCard("foreign")
	withOptHero(t, FakeHero{
		OptStrategy: func(cards []card.Card) (top, bottom []card.Card) {
			return []card.Card{a, foreign}, nil
		},
	}, func() {
		s := newOptTestEngine([]card.Card{a, b}, nil)
		assertPanics(t, "substituted card", "Opt:", func() { s.Opt(s.Logger(), 2) })
	})
}

// Tests that a handler duplicating one input card (and silently dropping another) panics.
// Multiset check catches the over-count of the duplicate before the leftover dropped
// card surfaces in the post-loop count.
func TestTurnStateOpt_PanicsOnDuplicatedCard(t *testing.T) {
	a := NewFakeCard("a")
	b := NewFakeCard("b")
	withOptHero(t, FakeHero{
		OptStrategy: func(cards []card.Card) (top, bottom []card.Card) {
			return []card.Card{a, a}, nil // duplicates a, drops b
		},
	}, func() {
		s := newOptTestEngine([]card.Card{a, b}, nil)
		assertPanics(t, "duplicated card", "Opt:", func() { s.Opt(s.Logger(), 2) })
	})
}

// sameDeck reports whether got contains the same cards in the same order as want by
// destructively draining the deck — Opt's effect is observable through subsequent draws,
// so the natural verification is to walk the deck via Draw. Mutates got; callers must
// discard it after the check.
func sameDeck(got *deck.Deck, want []card.Card) bool {
	if got.Size() != len(want) {
		return false
	}
	drawn := got.Draw(got.Size())
	for i, c := range drawn {
		if c.(card.Card) != want[i] {
			return false
		}
	}
	return true
}

// assertPanics runs fn and fails the test if it doesn't panic with a message containing
// substr. label appears in the failure message so call-site context survives diagnosis.
func assertPanics(t *testing.T, label, substr string, fn func()) {
	t.Helper()
	defer func() {
		r := recover()
		if r == nil {
			t.Errorf("%s: expected panic, got none", label)
			return
		}
		msg, ok := r.(string)
		if !ok {
			msg = "non-string panic"
		}
		if !strings.Contains(msg, substr) {
			t.Errorf("%s: panic message %q does not contain %q", label, msg, substr)
		}
	}()
	fn()
}
