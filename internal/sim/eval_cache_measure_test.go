package sim

// Cache hit-rate measurement. Loads mydecks/viserai_v4 and runs it through Evaluate at
// production shuffle counts, printing per-Evaluator cache stats. A real annealed list gives
// a realistic hit rate because Viserai's archetype is trigger-driven (cross-turn auras).
//
// Run with: `go test -run TestEvalCache_HitRateMeasurement -v`. Skipped in short mode and
// when mydecks/viserai_v4.json is absent.

import (
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card/cards"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/deck"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/gameengine"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/hero/heroes"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/textio"
)

// loadRealDeck reads mydecks/viserai_v4.json by walking up the directory tree. Returns nil
// when not found so callers can b.Skip / t.Skip cleanly.
func loadRealDeck(tb testing.TB) *deck.Deck {
	tb.Helper()
	dir, err := os.Getwd()
	if err != nil {
		tb.Fatalf("getwd: %v", err)
	}
	rel := filepath.Join("mydecks", "viserai_v4.json")
	for i := 0; i < 5; i++ {
		candidate := filepath.Join(dir, rel)
		if data, err := os.ReadFile(candidate); err == nil {
			loaded, _, err := textio.UnmarshalDeck(data)
			if err != nil {
				tb.Fatalf("unmarshal %s: %v", candidate, err)
			}
			return loaded
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	return nil
}

// Runs viserai_v4 through Evaluate at 10k shuffles and logs the cache stats. Output-only;
// no assertions.
func TestEvalCache_HitRateMeasurement(t *testing.T) {
	if testing.Short() {
		t.Skip("hit-rate measurement uses production shuffle counts; -short skips it")
	}
	const (
		incoming = 7
		shuffles = 10000
	)
	loaded := loadRealDeck(t)
	if loaded == nil {
		t.Skip("mydecks/viserai_v4.json not found — saved deck is needed to measure realistic hit rate")
	}
	baseline := loaded.Copy()

	// Use a dedicated Evaluator so we can read its cache stats after the run.
	ev := NewEvaluator()
	rng := rand.New(rand.NewSource(42))
	ev.Evaluate(baseline, shuffles, Matchup{IncomingPhysicalDamage: incoming}, rng)

	stats := ev.CacheStats()
	total := stats.Hits + stats.Misses
	t.Logf("cache stats over %d shuffles × ~10 turns/shuffle (~%d Best calls):", shuffles, total)
	t.Logf("  hits:        %d (%.1f%%)", stats.Hits, 100*stats.HitRate())
	t.Logf("  misses:      %d (%.1f%%)", stats.Misses, 100*float64(stats.Misses)/float64(max1(total)))
	t.Logf("  uncacheable: %d (%.1f%% of misses)", stats.Uncacheable, 100*float64(stats.Uncacheable)/float64(max1(stats.Misses)))
	t.Logf("  entries:     %d", stats.Entries)
}

func max1(n int) int {
	if n == 0 {
		return 1
	}
	return n
}

// Tests that NewEvaluatorParallel and NewEvaluator return identical Hands / TotalValue /
// Runs when fed RNG streams seeded the same way.
func TestEvalCache_ParallelEquivalentToSequential(t *testing.T) {
	const (
		deckSize   = 40
		maxCopies  = 2
		incoming   = 7
		masterSeed = int64(99)
		numWorkers = 4
	)
	setupRNG := rand.New(rand.NewSource(123))
	baseline := deck.Random(heroes.Viserai, deckSize, maxCopies, setupRNG, registry.Registry{})

	// Parallel pulls a per-worker seed from rng.Int63(); pre-extract the same int64s so the
	// per-worker RNGs are byte-identical and the sequential aggregate matches one parallel
	// chunk exactly. Scoped to one chunk so the seed extraction doesn't repeat.
	mirrorMaster := rand.New(rand.NewSource(masterSeed))
	workerSeeds := make([]int64, numWorkers)
	for i := range workerSeeds {
		workerSeeds[i] = mirrorMaster.Int63()
	}

	var seqStats deck.Stats
	for _, seed := range workerSeeds {
		workerRNG := rand.New(rand.NewSource(seed))
		workerStats := NewEvaluator().Evaluate(baseline.Copy(), 1, Matchup{IncomingPhysicalDamage: incoming}, workerRNG)
		mergeStats(&seqStats, workerStats)
	}

	parRNG := rand.New(rand.NewSource(masterSeed))
	parStats := NewEvaluatorParallel(numWorkers).Evaluate(baseline.Copy(), numWorkers, Matchup{IncomingPhysicalDamage: incoming}, parRNG)

	if seqStats.Hands != parStats.Hands {
		t.Errorf("Hands: seq=%d par=%d", seqStats.Hands, parStats.Hands)
	}
	if seqStats.TotalValue != parStats.TotalValue {
		t.Errorf("TotalValue: seq=%.0f par=%.0f delta=%.0f",
			seqStats.TotalValue, parStats.TotalValue, parStats.TotalValue-seqStats.TotalValue)
	}
	if seqStats.Runs != parStats.Runs {
		t.Errorf("Runs: seq=%d par=%d", seqStats.Runs, parStats.Runs)
	}
}

// mergeStats sums Hands, Runs, and TotalValue from src into dst.
func mergeStats(dst *deck.Stats, src deck.Stats) {
	dst.Hands += src.Hands
	dst.Runs += src.Runs
	dst.TotalValue += src.TotalValue
}

// Tests that ResetCache drops cache entries while preserving the hit/miss counters.
func TestEvalCache_ResetCache(t *testing.T) {
	ev := NewEvaluator()
	hand := []card.Card{cards.MaleficIncantationBlue{}, cards.MaleficIncantationBlue{}}

	// First call populates the cache (miss + store).
	ev.Best(nil, hand, nil, gameengine.GameStateBuilder().SetHero(heroes.Viserai).Build())
	preStats := ev.CacheStats()
	if preStats.Entries == 0 {
		t.Fatalf("expected cache to have an entry after first Best call")
	}

	// Second call hits the cache.
	ev.Best(nil, hand, nil, gameengine.GameStateBuilder().SetHero(heroes.Viserai).Build())
	if got := ev.CacheStats().Hits; got != preStats.Hits+1 {
		t.Errorf("hits = %d, want %d (one new hit on second call)", got, preStats.Hits+1)
	}

	// Reset drops entries; stats counters survive.
	ev.ResetCache()
	post := ev.CacheStats()
	if post.Entries != 0 {
		t.Errorf("Entries = %d after ResetCache, want 0", post.Entries)
	}
	if post.Hits != preStats.Hits+1 || post.Misses != preStats.Misses {
		t.Errorf("stats wiped by ResetCache: pre=%+v post=%+v", preStats, post)
	}

	// Same hand after reset is now a miss — confirms entries are actually gone.
	ev.Best(nil, hand, nil, gameengine.GameStateBuilder().SetHero(heroes.Viserai).Build())
	if got := ev.CacheStats().Misses; got != post.Misses+1 {
		t.Errorf("missed = %d, want %d (one new miss after reset)", got, post.Misses+1)
	}
}

// Tests that cached and uncached Best produce equal Value for the same hand inputs.
func TestEvalCache_PerHandEquivalence(t *testing.T) {
	hands := [][]card.Card{
		{cards.SkyFireLanternsRed{}, cards.MaleficIncantationBlue{}},
		{cards.MoonWishYellow{}, cards.FlyingHighRed{}},
		{cards.RavenousRabbleRed{}, cards.RavenousRabbleRed{}},
	}
	deck := DeckOf(cards.MaleficIncantationBlue{}, cards.SunKissRed{})
	cachedEv := NewEvaluator()
	freshEv := NewEvaluatorWithoutCache()
	for _, h := range hands {
		// Run twice to exercise cache hit on the second invocation.
		for i := 0; i < 2; i++ {
			cached := cachedEv.Best(nil, h, deck, gameengine.GameStateBuilder().SetHero(heroes.Viserai).Build())
			fresh := freshEv.Best(nil, h, deck, gameengine.GameStateBuilder().SetHero(heroes.Viserai).Build())
			if cached.Value != fresh.Value {
				t.Errorf("hand=%v iter=%d: cached.Value=%d fresh.Value=%d", h, i, cached.Value, fresh.Value)
			}
		}
	}
}

// Tests that cache-replay Evaluate and from-scratch Evaluate agree on Hands and
// TotalValue for matching RNG seeds.
func TestEvalCache_EquivalenceWithUncached(t *testing.T) {
	const (
		deckSize  = 40
		maxCopies = 2
		incoming  = 7
		shuffles  = 100
	)
	setupRNG := rand.New(rand.NewSource(123))
	baseline := deck.Random(heroes.Viserai, deckSize, maxCopies, setupRNG, registry.Registry{})

	cached := baseline.Copy()
	cachedStats := NewEvaluator().Evaluate(cached, shuffles, Matchup{IncomingPhysicalDamage: incoming}, rand.New(rand.NewSource(99)))

	uncached := baseline.Copy()
	uncachedStats := NewEvaluatorWithoutCache().Evaluate(uncached, shuffles, Matchup{IncomingPhysicalDamage: incoming}, rand.New(rand.NewSource(99)))

	if cachedStats.Hands != uncachedStats.Hands {
		t.Errorf("Hands: cached=%d uncached=%d", cachedStats.Hands, uncachedStats.Hands)
	}
	if cachedStats.TotalValue != uncachedStats.TotalValue {
		t.Errorf("TotalValue: cached=%.0f uncached=%.0f delta=%.0f", cachedStats.TotalValue, uncachedStats.TotalValue, cachedStats.TotalValue-uncachedStats.TotalValue)
	}
}

// BenchmarkEvalCache_SingleDeck compares one full Evaluate of viserai_v4 with the cache
// enabled vs disabled. Skipped when the saved deck is absent.
func BenchmarkEvalCache_SingleDeck(b *testing.B) {
	const (
		incoming = 7
		shuffles = 1000
	)
	loaded := loadRealDeck(b)
	if loaded == nil {
		b.Skip("mydecks/viserai_v4.json not found — saved deck needed for realistic bench")
	}

	b.Run("with-cache", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for n := 0; n < b.N; n++ {
			b.StopTimer()
			ev := NewEvaluator()
			rng := rand.New(rand.NewSource(42))
			d := loaded.Copy()
			b.StartTimer()
			ev.Evaluate(d, shuffles, Matchup{IncomingPhysicalDamage: incoming}, rng)
		}
	})
	b.Run("without-cache", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for n := 0; n < b.N; n++ {
			b.StopTimer()
			ev := NewEvaluatorWithoutCache()
			rng := rand.New(rand.NewSource(42))
			d := loaded.Copy()
			b.StartTimer()
			ev.Evaluate(d, shuffles, Matchup{IncomingPhysicalDamage: incoming}, rng)
		}
	})
}

// BenchmarkEvalCache_ParallelDeck runs viserai_v4 across NewEvaluatorParallel at varying
// worker counts; compare against BenchmarkEvalCache_SingleDeck's with-cache run.
func BenchmarkEvalCache_ParallelDeck(b *testing.B) {
	const (
		incoming = 7
		shuffles = 1000
	)
	loaded := loadRealDeck(b)
	if loaded == nil {
		b.Skip("mydecks/viserai_v4.json not found — saved deck needed for realistic bench")
	}
	for _, workers := range []int{1, 2, 4, 8} {
		w := workers
		b.Run(fmt.Sprintf("workers=%d", w), func(b *testing.B) {
			b.ReportAllocs()
			b.ResetTimer()
			for n := 0; n < b.N; n++ {
				b.StopTimer()
				ev := NewEvaluatorParallel(w)
				rng := rand.New(rand.NewSource(42))
				d := loaded.Copy()
				b.StartTimer()
				ev.Evaluate(d, shuffles, Matchup{IncomingPhysicalDamage: incoming}, rng)
			}
		})
	}
}
