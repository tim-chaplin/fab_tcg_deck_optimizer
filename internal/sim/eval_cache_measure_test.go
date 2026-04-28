package sim_test

// Cache hit-rate measurement. Runs a fixed-seed Viserai deck through Evaluate at
// production shuffle counts and prints the per-Evaluator cache stats. Produces the data
// the cacheable refactor was building toward — actual hit rate, plus the projected gain
// if priorAuraTriggers support were added.
//
// Run with: `go test -run TestEvalCache_HitRateMeasurement -v`. Skipped in short mode so
// it doesn't bloat normal `go test` runs.

import (
	"math/rand"
	"testing"

	. "github.com/tim-chaplin/fab-deck-optimizer/internal/sim"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/cards"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/heroes"
)

// TestEvalCache_HitRateMeasurement runs a fixed-shape Viserai deck through Evaluate and
// prints the cache stats. Not an assertion test — the t.Logf output is the deliverable.
func TestEvalCache_HitRateMeasurement(t *testing.T) {
	if testing.Short() {
		t.Skip("hit-rate measurement uses production shuffle counts; -short skips it")
	}
	const (
		deckSize  = 40
		maxCopies = 2
		incoming  = 7
		shuffles  = 5000
	)
	setupRNG := rand.New(rand.NewSource(42))
	baseline := Random(heroes.Viserai{}, deckSize, maxCopies, setupRNG, nil)

	// Wire a dedicated Evaluator into Evaluate via EvaluateWith so we can read its cache
	// stats after the run.
	ev := NewEvaluator()
	rng := rand.New(rand.NewSource(42))
	baseline.EvaluateWith(shuffles, incoming, rng, ev)

	stats := ev.CacheStats()
	total := stats.Hits + stats.Misses + stats.SkipsTriggers
	t.Logf("cache stats over %d shuffles × ~10 turns/shuffle (~%d Best calls):", shuffles, total)
	t.Logf("  hits:           %d (%.1f%%)", stats.Hits, 100*stats.HitRate())
	t.Logf("  misses:         %d (%.1f%%)", stats.Misses, 100*float64(stats.Misses)/float64(max1(total)))
	t.Logf("  skips-triggers: %d (%.1f%%)", stats.SkipsTriggers, 100*float64(stats.SkipsTriggers)/float64(max1(total)))
	t.Logf("  uncacheable:    %d (%.1f%% of misses)", stats.Uncacheable, 100*float64(stats.Uncacheable)/float64(max1(stats.Misses)))
	t.Logf("  entries:        %d", stats.Entries)
	t.Logf("  potential hit rate with trigger support: %.1f%% (current %.1f%% + skips %.1f%%)",
		100*stats.PotentialHitRateWithTriggers(),
		100*stats.HitRate(),
		100*float64(stats.SkipsTriggers)/float64(max1(total)))
}

func max1(n int) int {
	if n == 0 {
		return 1
	}
	return n
}

// TestEvalCache_PerHandEquivalence pins that for the same hand inputs, the cache-replay
// path produces a TurnSummary whose Value matches a from-scratch search. Walks several
// runs of a fixed-shape Viserai hand sequence, asserting Value equality on every Best
// call. This is the unit-level equivalence — the deck-eval-loop integration test below
// catches the same drift through the aggregate Stats.
func TestEvalCache_PerHandEquivalence(t *testing.T) {
	hands := [][]Card{
		{cards.SkyFireLanternsRed{}, cards.MaleficIncantationBlue{}},
		{cards.MoonWishYellow{}, cards.FlyingHighRed{}},
		{cards.RavenousRabbleRed{}, cards.RavenousRabbleRed{}},
	}
	deck := []Card{cards.MaleficIncantationBlue{}, cards.SunKissRed{}}
	cachedEv := NewEvaluator()
	freshEv := NewEvaluatorWithoutCache()
	for _, h := range hands {
		// Run twice to exercise cache hit on the second invocation.
		for i := 0; i < 2; i++ {
			cached := cachedEv.Best(heroes.Viserai{}, nil, h, 0, deck, 0, nil)
			fresh := freshEv.Best(heroes.Viserai{}, nil, h, 0, deck, 0, nil)
			if cached.Value != fresh.Value {
				t.Errorf("hand=%v iter=%d: cached.Value=%d fresh.Value=%d", h, i, cached.Value, fresh.Value)
			}
		}
	}
}

// TestEvalCache_EquivalenceWithUncached pins that the cache-replay path produces summary
// numbers WITHIN A SMALL TOLERANCE of a from-scratch search. The cache stores the winning
// partition's role multiset; replay applies that multiset to the new call. When multiple
// optimal partitions tie on Value/leftoverRunechants/futureValuePlayed/willOccupy, the
// from-scratch search picks the FIRST in iteration order (which depends on input hand
// order). Replay reproduces the cached multiset, which may differ from what fresh search
// would pick on a same-multiset hand in a different positional order. Both are valid
// optima — same Value for that turn — but the BestLine multiset can differ, which
// cascades through the deck-eval loop (different held / arsenal cards into the next turn,
// different next-hand multiset, different next-turn Value). Empirically the drift is
// tiny — order-of-1-Value over 100 shuffles on a Viserai deck — so we tolerate a small
// percentage gap.
func TestEvalCache_EquivalenceWithUncached(t *testing.T) {
	const (
		deckSize  = 40
		maxCopies = 2
		incoming  = 7
		shuffles  = 100
		// driftTolerance bounds the per-turn-mean absolute difference. Empirically <0.001
		// for the workloads we care about; 0.05 leaves substantial headroom while still
		// catching a real correctness regression (which would be orders of magnitude larger).
		driftTolerance = 0.05
	)
	setupRNG := rand.New(rand.NewSource(123))
	baseline := Random(heroes.Viserai{}, deckSize, maxCopies, setupRNG, nil)

	cached := New(baseline.Hero, baseline.Weapons, baseline.Cards)
	cached.EvaluateWith(shuffles, incoming, rand.New(rand.NewSource(99)), NewEvaluator())

	uncached := New(baseline.Hero, baseline.Weapons, baseline.Cards)
	uncached.EvaluateWith(shuffles, incoming, rand.New(rand.NewSource(99)), NewEvaluatorWithoutCache())

	if cached.Stats.Hands != uncached.Stats.Hands {
		t.Errorf("Hands: cached=%d uncached=%d", cached.Stats.Hands, uncached.Stats.Hands)
	}
	drift := cached.Stats.Mean() - uncached.Stats.Mean()
	if drift < -driftTolerance || drift > driftTolerance {
		t.Errorf("mean drift %.6f exceeds tolerance %.6f (cached=%.6f uncached=%.6f)",
			drift, driftTolerance, cached.Stats.Mean(), uncached.Stats.Mean())
	}
}

// BenchmarkEvalCache_SingleDeck compares one full Evaluate of a fixed-shape Viserai deck
// with the cache enabled vs disabled. Hit rate on a single-deck eval is high because the
// same hand multisets recur across shuffles within one deck — the per-deck workload is
// where the cache should pay off.
func BenchmarkEvalCache_SingleDeck(b *testing.B) {
	const (
		deckSize  = 40
		maxCopies = 2
		incoming  = 7
		shuffles  = 1000
	)
	setupRNG := rand.New(rand.NewSource(42))
	baseline := Random(heroes.Viserai{}, deckSize, maxCopies, setupRNG, nil)

	b.Run("with-cache", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for n := 0; n < b.N; n++ {
			b.StopTimer()
			ev := NewEvaluator()
			rng := rand.New(rand.NewSource(42))
			d := New(baseline.Hero, baseline.Weapons, baseline.Cards)
			b.StartTimer()
			d.EvaluateWith(shuffles, incoming, rng, ev)
		}
	})
	b.Run("without-cache", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for n := 0; n < b.N; n++ {
			b.StopTimer()
			ev := NewEvaluatorWithoutCache()
			rng := rand.New(rand.NewSource(42))
			d := New(baseline.Hero, baseline.Weapons, baseline.Cards)
			b.StartTimer()
			d.EvaluateWith(shuffles, incoming, rng, ev)
		}
	})
}
