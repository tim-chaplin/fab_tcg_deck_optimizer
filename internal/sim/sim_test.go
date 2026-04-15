package sim

import (
	"math/rand"
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card/fake"
)

func TestRun_AllRedDeckRecycles(t *testing.T) {
	// 40 reds. Optimal play for any 4-red hand is pitch 2, attack 2:
	// dealt 6, prevented 0 = 6. So 2 cards are pitched back per hand,
	// and the deck never empties to fewer than a full hand until many hands.
	deck := make([]card.Card, 40)
	for i := range deck {
		deck[i] = fake.Red{}
	}
	rng := rand.New(rand.NewSource(1))
	stats := Run(deck, 1, 4, rng)

	if stats.Hands < 10 {
		t.Fatalf("expected at least one full cycle (10 hands), got %d", stats.Hands)
	}
	if stats.Avg() != 6 {
		t.Fatalf("every all-red hand should value 6, got avg %v", stats.Avg())
	}
	if stats.FirstCycle.Hands != 10 {
		t.Fatalf("expected 10 first-cycle hands, got %d", stats.FirstCycle.Hands)
	}
	if stats.SecondCycle.Hands == 0 {
		t.Fatalf("expected recycling to produce a second cycle")
	}
}

func TestRun_DeckEventuallyEmpties(t *testing.T) {
	// All-red deck pitches 2 per hand and burns 2 per hand. Deck shrinks
	// by 2 each turn, so it should terminate in finite time.
	deck := make([]card.Card, 40)
	for i := range deck {
		deck[i] = fake.Red{}
	}
	rng := rand.New(rand.NewSource(2))
	stats := Run(deck, 1, 4, rng)

	// Bound: cards consumed = 2 per hand, starting from 40, stops when <4 left.
	// So roughly (40-3)/2 ~= 19 hands max.
	if stats.Hands > 25 {
		t.Fatalf("simulation did not terminate cleanly: %d hands", stats.Hands)
	}
}

func TestRun_DeterministicWithSeed(t *testing.T) {
	deck := mixedDeck()
	a := Run(deck, 100, 4, rand.New(rand.NewSource(42)))
	b := Run(deck, 100, 4, rand.New(rand.NewSource(42)))
	if a.TotalValue != b.TotalValue || a.Hands != b.Hands {
		t.Fatalf("same seed should produce identical results: %+v vs %+v", a, b)
	}
}

func TestRun_MultipleRunsAggregate(t *testing.T) {
	deck := mixedDeck()
	rng := rand.New(rand.NewSource(7))
	stats := Run(deck, 50, 4, rng)
	if stats.Runs != 50 {
		t.Fatalf("Runs = %d, want 50", stats.Runs)
	}
	if stats.Hands < 50*10 {
		t.Fatalf("expected at least 500 hands across 50 runs, got %d", stats.Hands)
	}
	if stats.Avg() <= 0 {
		t.Fatalf("expected positive average value, got %v", stats.Avg())
	}
}

func mixedDeck() []card.Card {
	d := make([]card.Card, 0, 40)
	for i := 0; i < 20; i++ {
		d = append(d, fake.Blue{})
	}
	for i := 0; i < 20; i++ {
		d = append(d, fake.Red{})
	}
	return d
}

