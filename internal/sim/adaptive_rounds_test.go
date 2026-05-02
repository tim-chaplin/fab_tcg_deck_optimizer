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

	"github.com/tim-chaplin/fab-deck-optimizer/internal/heroes"
)

// TestAdaptive_RoundsToConverge reports Stats.Runs (the actual shuffle count adaptive
// stopped at) for several deck shapes. Results vs the per-fan-out chunk size of
// numWorkers × adaptiveCheckInterval tell us whether the barrier-merge ever fires:
//   - Stats.Runs <= chunkSize → adaptive converged inside the first chunk, no barrier
//   - Stats.Runs > chunkSize  → at least one barrier fired
//
// Run with: `go test -run TestAdaptive_RoundsToConverge -v`. Logged output is the
// deliverable — no assertion, just data.
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
		baseline := Random(heroes.Viserai{}, deckSize, maxCopies, setupRNG, nil)
		d := New(baseline.Hero, baseline.Weapons, baseline.Cards)
		ev := NewEvaluatorParallel(numWorkers)
		stats := d.EvaluateAdaptiveWith(incoming, 0, rand.New(rand.NewSource(99)), ev)
		t.Logf("random Viserai seed=%d: Runs=%d, mean=%.4f", seed, stats.Runs, stats.Mean())
	}

	// And the high-quality annealed list when available — different convergence profile
	// because the card distribution is tightened.
	if loaded := loadRealDeck(t); loaded != nil {
		d := New(loaded.Hero, loaded.Weapons, loaded.Cards)
		ev := NewEvaluatorParallel(numWorkers)
		stats := d.EvaluateAdaptiveWith(incoming, 0, rand.New(rand.NewSource(99)), ev)
		t.Logf("viserai_v4 (annealed): Runs=%d, mean=%.4f", stats.Runs, stats.Mean())
	}
}
