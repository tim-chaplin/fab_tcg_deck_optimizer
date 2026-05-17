package sim

// Cache hit-rate measurement. Loads a high-quality saved Viserai deck (mydecks/viserai_v4)
// and runs it through Evaluate at production shuffle counts, printing per-Evaluator cache
// stats. Annealed decks have far higher trigger-skip rates than randomly-generated shapes
// because Viserai's archetype is trigger-driven (Sigil of Silphidae, Malefic Incantation,
// etc. carry across turns); using a real annealed list gives a realistic picture of what
// the cache buys in production.
//
// Run with: `go test -run TestEvalCache_HitRateMeasurement -v`. Skipped in short mode so
// it doesn't bloat normal `go test` runs. Skipped when mydecks/viserai_v4.json is absent
// so go test ./... still passes on a fresh checkout that doesn't carry the saved deck.

import (
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/v2/deck"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/registry"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/textio"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/cards"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/gameengine"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/hero"
)

// loadRealDeck reads mydecks/viserai_v4.json from somewhere up the directory tree.
// Returns nil when the file isn't found so callers can b.Skip / t.Skip cleanly. Mirrors
// cmd/fabsim/eval_realdeck_bench_test.go's findRepoFile helper but specialised to the
// load path so the cache tests have one self-contained loader.
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

// TestEvalCache_HitRateMeasurement loads viserai_v4 (a high-quality annealed list) and
// runs it through Evaluate at the production 10k shuffle count, printing the cache stats.
// Not an assertion test — the t.Logf output is the deliverable.
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
	ev.Evaluate(baseline, shuffles, Matchup{IncomingDamage: incoming}, rng)

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

// Tests that the parallel-shuffle and sequential paths agree on per-turn-mean within a
// small tolerance — parallel derives per-worker seeds, so the streams differ but converge.
func TestEvalCache_ParallelEquivalentToSequential(t *testing.T) {
	const (
		deckSize  = 40
		maxCopies = 2
		incoming  = 7
		shuffles  = 1000
		// drift bound: per-turn-mean across 1k shuffles is empirically ±0.1 between
		// parallel and sequential RNG streams. 0.5 leaves plenty of headroom while
		// catching a real correctness regression (which would shift by ≥1 unit).
		driftTolerance = 0.5
	)
	setupRNG := rand.New(rand.NewSource(123))
	baseline := deck.Random(hero.Viserai{}, deckSize, maxCopies, setupRNG, nil, registry.Registry{})

	seq := baseline.Copy()
	seqStats := NewEvaluator().Evaluate(seq, shuffles, Matchup{IncomingDamage: incoming}, rand.New(rand.NewSource(99)))

	par := baseline.Copy()
	parStats := NewEvaluatorParallel(4).Evaluate(par, shuffles, Matchup{IncomingDamage: incoming}, rand.New(rand.NewSource(99)))

	// Hands counts can differ slightly because parallel and sequential consume different
	// per-shuffle RNG streams: a shuffle that runs out of deck cards on hand 7 in one
	// stream might deal hand 8 in the other. Empirically <1% drift on 1k-shuffle runs.
	handsRatio := float64(seqStats.Hands-parStats.Hands) / float64(seqStats.Hands)
	if handsRatio < -0.02 || handsRatio > 0.02 {
		t.Errorf("Hands count drift %.4f exceeds 2%% (seq=%d par=%d)", handsRatio, seqStats.Hands, parStats.Hands)
	}
	drift := seqStats.Mean() - parStats.Mean()
	if drift < -driftTolerance || drift > driftTolerance {
		t.Errorf("mean drift %.6f exceeds tolerance %.6f (seq=%.6f par=%.6f)",
			drift, driftTolerance, seqStats.Mean(), parStats.Mean())
	}
}

// Tests that ResetCache drops cache entries while preserving the hit/miss counters.
func TestEvalCache_ResetCache(t *testing.T) {
	ev := NewEvaluator()
	hand := []card.Card{cards.MaleficIncantationBlue{}, cards.MaleficIncantationBlue{}}

	// First call populates the cache (miss + store).
	ev.Best(nil, hand, Matchup{IncomingDamage: 0}, nil, gameengine.GameStateBuilder().SetHero(hero.Viserai{}).Build())
	preStats := ev.CacheStats()
	if preStats.Entries == 0 {
		t.Fatalf("expected cache to have an entry after first Best call")
	}

	// Second call hits the cache.
	ev.Best(nil, hand, Matchup{IncomingDamage: 0}, nil, gameengine.GameStateBuilder().SetHero(hero.Viserai{}).Build())
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

	// Same hand after reset is now a miss — confirms entries are actually gone, not just
	// the count reading wrong.
	ev.Best(nil, hand, Matchup{IncomingDamage: 0}, nil, gameengine.GameStateBuilder().SetHero(hero.Viserai{}).Build())
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
			cached := cachedEv.Best(nil, h, Matchup{IncomingDamage: 0}, deck, gameengine.GameStateBuilder().SetHero(hero.Viserai{}).Build())
			fresh := freshEv.Best(nil, h, Matchup{IncomingDamage: 0}, deck, gameengine.GameStateBuilder().SetHero(hero.Viserai{}).Build())
			if cached.Value != fresh.Value {
				t.Errorf("hand=%v iter=%d: cached.Value=%d fresh.Value=%d", h, i, cached.Value, fresh.Value)
			}
		}
	}
}

// Tests that the cache-replay path produces summary numbers equal (within driftTolerance)
// to a from-scratch search at the deck-eval-loop level.
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
	baseline := deck.Random(hero.Viserai{}, deckSize, maxCopies, setupRNG, nil, registry.Registry{})

	cached := baseline.Copy()
	cachedStats := NewEvaluator().Evaluate(cached, shuffles, Matchup{IncomingDamage: incoming}, rand.New(rand.NewSource(99)))

	uncached := baseline.Copy()
	uncachedStats := NewEvaluatorWithoutCache().Evaluate(uncached, shuffles, Matchup{IncomingDamage: incoming}, rand.New(rand.NewSource(99)))

	if cachedStats.Hands != uncachedStats.Hands {
		t.Errorf("Hands: cached=%d uncached=%d", cachedStats.Hands, uncachedStats.Hands)
	}
	drift := cachedStats.Mean() - uncachedStats.Mean()
	if drift < -driftTolerance || drift > driftTolerance {
		t.Errorf("mean drift %.6f exceeds tolerance %.6f (cached=%.6f uncached=%.6f)",
			drift, driftTolerance, cachedStats.Mean(), uncachedStats.Mean())
	}
}

// BenchmarkEvalCache_SingleDeck compares one full Evaluate of viserai_v4 (a high-quality
// annealed Viserai list) with the cache enabled vs disabled. Real annealed decks are the
// production target — random Viserai shapes have a different cache-hit profile because
// they don't carry the trigger-heavy archetype synergies that drive Viserai's actual
// gameplay. Skipped when the saved deck is absent.
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
			ev.Evaluate(d, shuffles, Matchup{IncomingDamage: incoming}, rng)
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
			ev.Evaluate(d, shuffles, Matchup{IncomingDamage: incoming}, rng)
		}
	})
}

// BenchmarkEvalCache_ParallelDeck mirrors BenchmarkEvalCache_SingleDeck but runs the
// shuffle loop across multiple workers via NewEvaluatorParallel. Compares to the existing
// single-threaded with-cache run as the baseline.
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
				ev.Evaluate(d, shuffles, Matchup{IncomingDamage: incoming}, rng)
			}
		})
	}
}
