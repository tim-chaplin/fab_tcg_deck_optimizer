package sim_test

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

	. "github.com/tim-chaplin/fab-deck-optimizer/internal/sim"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/cards"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/deckio"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/heroes"
)

// loadRealDeck reads mydecks/viserai_v4.json from somewhere up the directory tree.
// Returns nil when the file isn't found so callers can b.Skip / t.Skip cleanly. Mirrors
// cmd/fabsim/eval_realdeck_bench_test.go's findRepoFile helper but specialised to the
// load path so the cache tests have one self-contained loader.
func loadRealDeck(tb testing.TB) *Deck {
	tb.Helper()
	dir, err := os.Getwd()
	if err != nil {
		tb.Fatalf("getwd: %v", err)
	}
	rel := filepath.Join("mydecks", "viserai_v4.json")
	for i := 0; i < 5; i++ {
		candidate := filepath.Join(dir, rel)
		if data, err := os.ReadFile(candidate); err == nil {
			loaded, err := deckio.Unmarshal(data)
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
	baseline := New(loaded.Hero, loaded.Weapons, loaded.Cards)

	// Wire a dedicated Evaluator into Evaluate via EvaluateWith so we can read its cache
	// stats after the run.
	ev := NewEvaluator()
	rng := rand.New(rand.NewSource(42))
	baseline.EvaluateWith(shuffles, Matchup{IncomingDamage: incoming}, rng, ev)

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

// TestEvalCache_ParallelEquivalentToSequential pins that the parallel-shuffle path
// produces the same per-turn-mean within a small tolerance as the single-threaded path.
// Same fixed-seed deck on both, run once with NewEvaluatorParallel(N) and once with
// NewEvaluator(); the per-turn-mean should agree at the same shuffle count modulo the
// expected RNG-distribution drift (parallel path derives per-worker seeds from the input
// rng, so the actual sequence of shuffled decks differs from sequential — but the mean
// over enough shuffles converges to the same number).
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
	baseline := Random(heroes.Viserai{}, deckSize, maxCopies, setupRNG, nil)

	seq := New(baseline.Hero, baseline.Weapons, baseline.Cards)
	seq.EvaluateWith(shuffles, Matchup{IncomingDamage: incoming}, rand.New(rand.NewSource(99)), NewEvaluator())

	par := New(baseline.Hero, baseline.Weapons, baseline.Cards)
	par.EvaluateWith(shuffles, Matchup{IncomingDamage: incoming}, rand.New(rand.NewSource(99)), NewEvaluatorParallel(4))

	// Hands counts can differ slightly because parallel and sequential consume different
	// per-shuffle RNG streams: a shuffle that runs out of deck cards on hand 7 in one
	// stream might deal hand 8 in the other. Empirically <1% drift on 1k-shuffle runs.
	handsRatio := float64(seq.Stats.Hands-par.Stats.Hands) / float64(seq.Stats.Hands)
	if handsRatio < -0.02 || handsRatio > 0.02 {
		t.Errorf("Hands count drift %.4f exceeds 2%% (seq=%d par=%d)", handsRatio, seq.Stats.Hands, par.Stats.Hands)
	}
	drift := seq.Stats.Mean() - par.Stats.Mean()
	if drift < -driftTolerance || drift > driftTolerance {
		t.Errorf("mean drift %.6f exceeds tolerance %.6f (seq=%.6f par=%.6f)",
			drift, driftTolerance, seq.Stats.Mean(), par.Stats.Mean())
	}
}

// TestEvalCache_ResetCache pins that ResetCache drops cached entries while leaving the
// stats counters intact, so the iterate-mode worker pool can clear the cache between
// mutations without losing the running hit/miss tally for diagnostics.
func TestEvalCache_ResetCache(t *testing.T) {
	ev := NewEvaluator()
	hand := []Card{cards.MaleficIncantationBlue{}, cards.MaleficIncantationBlue{}}

	// First call populates the cache (miss + store).
	ev.Best(heroes.Viserai{}, nil, hand, Matchup{IncomingDamage: 0}, nil, 0, nil)
	preStats := ev.CacheStats()
	if preStats.Entries == 0 {
		t.Fatalf("expected cache to have an entry after first Best call")
	}

	// Second call hits the cache.
	ev.Best(heroes.Viserai{}, nil, hand, Matchup{IncomingDamage: 0}, nil, 0, nil)
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
	ev.Best(heroes.Viserai{}, nil, hand, Matchup{IncomingDamage: 0}, nil, 0, nil)
	if got := ev.CacheStats().Misses; got != post.Misses+1 {
		t.Errorf("missed = %d, want %d (one new miss after reset)", got, post.Misses+1)
	}
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
			cached := cachedEv.Best(heroes.Viserai{}, nil, h, Matchup{IncomingDamage: 0}, deck, 0, nil)
			fresh := freshEv.Best(heroes.Viserai{}, nil, h, Matchup{IncomingDamage: 0}, deck, 0, nil)
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
	cached.EvaluateWith(shuffles, Matchup{IncomingDamage: incoming}, rand.New(rand.NewSource(99)), NewEvaluator())

	uncached := New(baseline.Hero, baseline.Weapons, baseline.Cards)
	uncached.EvaluateWith(shuffles, Matchup{IncomingDamage: incoming}, rand.New(rand.NewSource(99)), NewEvaluatorWithoutCache())

	if cached.Stats.Hands != uncached.Stats.Hands {
		t.Errorf("Hands: cached=%d uncached=%d", cached.Stats.Hands, uncached.Stats.Hands)
	}
	drift := cached.Stats.Mean() - uncached.Stats.Mean()
	if drift < -driftTolerance || drift > driftTolerance {
		t.Errorf("mean drift %.6f exceeds tolerance %.6f (cached=%.6f uncached=%.6f)",
			drift, driftTolerance, cached.Stats.Mean(), uncached.Stats.Mean())
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
			d := New(loaded.Hero, loaded.Weapons, loaded.Cards)
			b.StartTimer()
			d.EvaluateWith(shuffles, Matchup{IncomingDamage: incoming}, rng, ev)
		}
	})
	b.Run("without-cache", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for n := 0; n < b.N; n++ {
			b.StopTimer()
			ev := NewEvaluatorWithoutCache()
			rng := rand.New(rand.NewSource(42))
			d := New(loaded.Hero, loaded.Weapons, loaded.Cards)
			b.StartTimer()
			d.EvaluateWith(shuffles, Matchup{IncomingDamage: incoming}, rng, ev)
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
				d := New(loaded.Hero, loaded.Weapons, loaded.Cards)
				b.StartTimer()
				d.EvaluateWith(shuffles, Matchup{IncomingDamage: incoming}, rng, ev)
			}
		})
	}
}
