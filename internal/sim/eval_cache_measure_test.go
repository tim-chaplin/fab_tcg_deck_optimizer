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
