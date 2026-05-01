package sim_test

import (
	"strings"
	"testing"

	. "github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
)

// withOptHero swaps CurrentHero to the supplied testutils.Hero for the test's lifetime
// and restores nil afterwards. Mirrors the pattern in lower_health_wanter_test.go.
func withOptHero(t *testing.T, h testutils.Hero, fn func()) {
	t.Helper()
	prev := CurrentHero
	CurrentHero = h
	defer func() { CurrentHero = prev }()
	fn()
}

// Tests that Opt with the default passthrough handler keeps the deck order unchanged —
// every revealed card returns to the top in input order, none move to the bottom.
func TestTurnStateOpt_PassthroughKeepsDeckOrder(t *testing.T) {
	a := testutils.NewStubCard("a")
	b := testutils.NewStubCard("b")
	c := testutils.NewStubCard("c")
	d := testutils.NewStubCard("d")
	withOptHero(t, testutils.Hero{Intel: 4}, func() {
		s := NewTurnState([]Card{a, b, c, d}, nil)
		s.Opt(2)
		got := s.Deck()
		want := []Card{a, b, c, d}
		if !sameDeck(got, want) {
			t.Errorf("deck = %v, want %v (passthrough)", got, want)
		}
	})
}

// Tests that the handler can move cards to the bottom of the deck. Handler bottoms the
// first revealed card and keeps the second on top; the un-opted tail of the deck stays
// in place, and the bottomed card lands at the end.
func TestTurnStateOpt_BottomsHandlerSpecifiedCards(t *testing.T) {
	a := testutils.NewStubCard("a")
	b := testutils.NewStubCard("b")
	c := testutils.NewStubCard("c")
	withOptHero(t, testutils.Hero{
		OptStrategy: func(cards []Card) (top, bottom []Card) {
			// Bottom cards[0] (a), keep cards[1] (b) on top.
			return []Card{cards[1]}, []Card{cards[0]}
		},
	}, func() {
		s := NewTurnState([]Card{a, b, c}, nil)
		s.Opt(2)
		got := s.Deck()
		// Handler saw [a, b]; returned top=[b], bottom=[a]. Deck becomes [b] + [c] + [a].
		want := []Card{b, c, a}
		if !sameDeck(got, want) {
			t.Errorf("deck = %v, want %v (handler bottoms a)", got, want)
		}
	})
}

// Tests that the handler can re-order cards on top in addition to bottoming.
func TestTurnStateOpt_HandlerReorderCanReverseTop(t *testing.T) {
	a := testutils.NewStubCard("a")
	b := testutils.NewStubCard("b")
	c := testutils.NewStubCard("c")
	withOptHero(t, testutils.Hero{
		OptStrategy: func(cards []Card) (top, bottom []Card) {
			return []Card{cards[1], cards[0]}, nil
		},
	}, func() {
		s := NewTurnState([]Card{a, b, c}, nil)
		s.Opt(2)
		got := s.Deck()
		want := []Card{b, a, c}
		if !sameDeck(got, want) {
			t.Errorf("deck = %v, want %v (top reversed)", got, want)
		}
	})
}

// Tests that a request larger than the deck length clamps to whatever cards are there.
// The handler sees only the available cards and reshapes them; the empty-tail-deck
// remains empty.
func TestTurnStateOpt_ClampsNToDeckLength(t *testing.T) {
	a := testutils.NewStubCard("a")
	b := testutils.NewStubCard("b")
	withOptHero(t, testutils.Hero{
		OptStrategy: func(cards []Card) (top, bottom []Card) {
			if len(cards) != 2 {
				t.Errorf("handler saw %d cards, want 2 (clamp)", len(cards))
			}
			return cards, nil
		},
	}, func() {
		s := NewTurnState([]Card{a, b}, nil)
		s.Opt(5)
		got := s.Deck()
		want := []Card{a, b}
		if !sameDeck(got, want) {
			t.Errorf("deck = %v, want %v (clamped)", got, want)
		}
	})
}

// Tests that Opt on an empty deck never invokes the handler and is a safe no-op.
func TestTurnStateOpt_EmptyDeckSkipsHandler(t *testing.T) {
	called := false
	withOptHero(t, testutils.Hero{
		OptStrategy: func(cards []Card) (top, bottom []Card) {
			called = true
			return cards, nil
		},
	}, func() {
		s := NewTurnState(nil, nil)
		s.Opt(3)
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
		withOptHero(t, testutils.Hero{
			OptStrategy: func(cards []Card) (top, bottom []Card) {
				called = true
				return cards, nil
			},
		}, func() {
			s := NewTurnState([]Card{testutils.NewStubCard("x")}, nil)
			s.Opt(n)
			if called {
				t.Errorf("Opt(%d) called the handler, want skip", n)
			}
		})
	}
}

// Tests that Opt emits a "Opted X, put Y on top, put Z on bottom" log entry naming the
// revealed cards and the handler's split when it ran.
func TestTurnStateOpt_LogsOutcome(t *testing.T) {
	a := testutils.NewStubCard("a")
	b := testutils.NewStubCard("b")
	c := testutils.NewStubCard("c")
	withOptHero(t, testutils.Hero{
		OptStrategy: func(cards []Card) (top, bottom []Card) {
			return []Card{cards[1]}, []Card{cards[0]}
		},
	}, func() {
		s := NewTurnState([]Card{a, b, c}, nil)
		s.Opt(2)
		if len(s.Log) != 1 {
			t.Fatalf("Log len = %d, want 1", len(s.Log))
		}
		want := "Opted [a, b], put [b] on top, put [a] on bottom"
		if got := s.Log[0].Text; got != want {
			t.Errorf("log entry = %q, want %q", got, want)
		}
		if s.Log[0].N != 0 {
			t.Errorf("log N = %d, want 0 (Opt is value-neutral; reshape effect surfaces in later turns)", s.Log[0].N)
		}
	})
}

// Tests that the log entry renders an empty top or bottom list as "[]".
func TestTurnStateOpt_LogShowsEmptyListsAsBrackets(t *testing.T) {
	a := testutils.NewStubCard("a")
	b := testutils.NewStubCard("b")
	withOptHero(t, testutils.Hero{
		OptStrategy: func(cards []Card) (top, bottom []Card) {
			return cards, nil // bottom empty
		},
	}, func() {
		s := NewTurnState([]Card{a, b}, nil)
		s.Opt(2)
		want := "Opted [a, b], put [a, b] on top, put [] on bottom"
		if got := s.Log[0].Text; got != want {
			t.Errorf("log entry = %q, want %q", got, want)
		}
	})
}

// Tests that the no-op paths skip the log entry entirely.
func TestTurnStateOpt_NoOpPathsSkipLog(t *testing.T) {
	cases := []struct {
		name string
		deck []Card
		n    int
	}{
		{"empty deck", nil, 3},
		{"zero n", []Card{testutils.NewStubCard("x")}, 0},
		{"negative n", []Card{testutils.NewStubCard("x")}, -1},
	}
	withOptHero(t, testutils.Hero{}, func() {
		for _, tc := range cases {
			s := NewTurnState(tc.deck, nil)
			s.Opt(tc.n)
			if len(s.Log) != 0 {
				t.Errorf("%s: Log = %v, want empty", tc.name, s.Log)
			}
		}
	})
}

// Tests that Opt always flips IsCacheable to false, even on the no-op paths (n <= 0,
// empty deck) — the chain reading the deck implies an order dependency regardless of
// whether the handler ran.
func TestTurnStateOpt_AlwaysFlipsCacheable(t *testing.T) {
	cases := []struct {
		name string
		deck []Card
		n    int
	}{
		{"populated deck", []Card{testutils.NewStubCard("x")}, 1},
		{"empty deck", nil, 3},
		{"zero n", []Card{testutils.NewStubCard("x")}, 0},
	}
	withOptHero(t, testutils.Hero{}, func() {
		for _, tc := range cases {
			s := NewTurnState(tc.deck, nil)
			if !s.IsCacheable() {
				t.Fatalf("%s: pre IsCacheable should be true", tc.name)
			}
			s.Opt(tc.n)
			if s.IsCacheable() {
				t.Errorf("%s: Opt(%d) should flip IsCacheable to false", tc.name, tc.n)
			}
		}
	})
}

// Tests that a handler returning fewer cards than received panics.
func TestTurnStateOpt_PanicsOnDroppedCard(t *testing.T) {
	withOptHero(t, testutils.Hero{
		OptStrategy: func(cards []Card) (top, bottom []Card) {
			return []Card{cards[0]}, nil // drops cards[1]
		},
	}, func() {
		s := NewTurnState([]Card{testutils.NewStubCard("a"), testutils.NewStubCard("b")}, nil)
		assertPanics(t, "dropped card", "Opt:", func() { s.Opt(2) })
	})
}

// Tests that a handler returning more cards than received panics.
func TestTurnStateOpt_PanicsOnExtraCard(t *testing.T) {
	a := testutils.NewStubCard("a")
	b := testutils.NewStubCard("b")
	foreign := testutils.NewStubCard("foreign")
	withOptHero(t, testutils.Hero{
		OptStrategy: func(cards []Card) (top, bottom []Card) {
			return []Card{a, b, foreign}, nil
		},
	}, func() {
		s := NewTurnState([]Card{a, b}, nil)
		assertPanics(t, "extra card", "Opt:", func() { s.Opt(2) })
	})
}

// Tests that a handler substituting one input card for a non-input card panics — the
// length check passes but the multiset check trips on the foreign card.
func TestTurnStateOpt_PanicsOnSubstitutedCard(t *testing.T) {
	a := testutils.NewStubCard("a")
	b := testutils.NewStubCard("b")
	foreign := testutils.NewStubCard("foreign")
	withOptHero(t, testutils.Hero{
		OptStrategy: func(cards []Card) (top, bottom []Card) {
			return []Card{a, foreign}, nil
		},
	}, func() {
		s := NewTurnState([]Card{a, b}, nil)
		assertPanics(t, "substituted card", "Opt:", func() { s.Opt(2) })
	})
}

// Tests that a handler duplicating one input card (and silently dropping another) panics.
// Multiset check catches the over-count of the duplicate before the leftover dropped
// card surfaces in the post-loop count.
func TestTurnStateOpt_PanicsOnDuplicatedCard(t *testing.T) {
	a := testutils.NewStubCard("a")
	b := testutils.NewStubCard("b")
	withOptHero(t, testutils.Hero{
		OptStrategy: func(cards []Card) (top, bottom []Card) {
			return []Card{a, a}, nil // duplicates a, drops b
		},
	}, func() {
		s := NewTurnState([]Card{a, b}, nil)
		assertPanics(t, "duplicated card", "Opt:", func() { s.Opt(2) })
	})
}

// sameDeck reports whether two card slices contain the same cards in the same order.
func sameDeck(got, want []Card) bool {
	if len(got) != len(want) {
		return false
	}
	for i := range got {
		if got[i] != want[i] {
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
