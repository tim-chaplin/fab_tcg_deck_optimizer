package sim_test

// Adaptive-round count experiment. Runs adaptive eval on a range of decks (random Viserai
// at multiple seeds + the saved annealed list) and reports how many parallel chunks each
// eval needed to hit the SE target. The chunk size is numWorkers × adaptiveCheckInterval =
// 8 × 1000 = 8000 shuffles by default; if every deck converges inside one chunk, we never
// pay the barrier-merge cost and the parallel-shuffle path's overhead amortises away.

import (
	"math/rand"
	"testing"

	. "github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/deck"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/registry"

	"github.com/tim-chaplin/fab-deck-optimizer/v2/hero"
)

// Logs Stats.Runs for various deck shapes — `Runs <= numWorkers × adaptiveCheckInterval`
// means adaptive converged inside the first chunk (no barrier-merge). Data-only, no asserts.
func TestAdaptive_RoundsToConverge(t *testing.T) {
	if testing.Short() {
		t.Skip("adaptive convergence experiment is slow; -short skips it")
	}
	const (
		deckSize  = 40
		maxCopies = 2
		incoming  = 7
	)
	numWorkers := DefaultWorkers()
	t.Logf("numWorkers=%d", numWorkers)

	// Test on a range of random Viserai decks (different seeds → different card mixes →
	// different variance profiles).
	for _, seed := range []int64{1, 2, 3, 42, 123, 999} {
		setupRNG := rand.New(rand.NewSource(seed))
		baseline := deck.Random(hero.Viserai{}, deckSize, maxCopies, setupRNG, nil, registry.Registry{})
		d := baseline.Copy()
		ev := NewEvaluatorParallel(numWorkers)
		stats := ev.EvaluateAdaptive(d, 0.1, Matchup{IncomingDamage: incoming}, rand.New(rand.NewSource(99)))
		t.Logf("random Viserai seed=%d: Runs=%d, mean=%.4f", seed, stats.Runs, stats.Mean())
	}

	// And the high-quality annealed list when available — different convergence profile
	// because the card distribution is tightened.
	if loaded := loadRealDeck(t); loaded != nil {
		d := loaded.Copy()
		ev := NewEvaluatorParallel(numWorkers)
		stats := ev.EvaluateAdaptive(d, 0.1, Matchup{IncomingDamage: incoming}, rand.New(rand.NewSource(99)))
		t.Logf("viserai_v4 (annealed): Runs=%d, mean=%.4f", stats.Runs, stats.Mean())
	}
}
