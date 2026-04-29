package sim_test

// Worker-count sweeps for the two production parallelism dimensions:
//
//   - BenchmarkEvalWorkerSweep: eval mode (single deck, shuffle parallelism only).
//     Sweeps shuffleWorkers across {1..16}; report the wall-clock per Evaluate so the
//     curve where shuffle parallelism stops paying off is visible.
//   - BenchmarkAnnealWorkerSweep: anneal mode (N-mutation queue, both dimensions).
//     Sweeps a 2D grid of (mutationWorkers, shuffleWorkers) and reports per-round wall.
//
// Both load mydecks/viserai_v4 (the high-quality annealed list reused by the cache
// measurement bench) so the workload matches production. Skipped when the deck file
// isn't present.
//
// Recommended invocation:
//
//	go test -bench='BenchmarkEvalWorkerSweep|BenchmarkAnnealWorkerSweep' \
//	    -run='^$' -benchtime=2x -count=1 ./internal/sim
//
// Each sub-bench runs at b.N=2 by default, balancing noise reduction against total runtime
// (the grid is ~2 minutes total at -benchtime=2x on the 8-physical-core test machine).

import (
	"context"
	"fmt"
	"math/rand"
	"testing"

	. "github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

// BenchmarkEvalWorkerSweep measures one Evaluate call on viserai_v4 across a range of
// shuffle-worker counts. Goal: find the inflection point where adding shuffle workers
// stops cutting wall-clock — past which the chunk-merge barrier and cache contention
// outweigh the added parallelism.
//
// Uses fixed shuffles=1000 (rather than adaptive) so every sub-bench does identical
// total work and the ns/op reads as wall-clock-per-1k-shuffles directly.
func BenchmarkEvalWorkerSweep(b *testing.B) {
	const (
		incoming = 7
		shuffles = 1000
	)
	loaded := loadRealDeck(b)
	if loaded == nil {
		b.Skip("mydecks/viserai_v4.json not found — saved deck needed for realistic bench")
	}
	for _, w := range []int{1, 2, 4, 6, 8, 10, 12, 16} {
		w := w
		b.Run(fmt.Sprintf("shuffle-workers=%d", w), func(b *testing.B) {
			b.ReportAllocs()
			b.ResetTimer()
			for n := 0; n < b.N; n++ {
				b.StopTimer()
				ev := NewEvaluatorParallel(w)
				rng := rand.New(rand.NewSource(42))
				d := New(loaded.Hero, loaded.Weapons, loaded.Cards)
				b.StartTimer()
				d.EvaluateWith(shuffles, incoming, rng, ev)
			}
		})
	}
}

// BenchmarkAnnealWorkerSweep measures one IterateParallel call across a 2D grid of
// (mutationWorkers, shuffleWorkers) on viserai_v4. The mutation list is the first
// mutationSampleSize entries of AllMutations(viserai_v4), evaluated against an unreachable
// baseline so the worker pool drains every sampled mutation (production rounds short-
// circuit on the first improvement, but the per-mutation cost is the same; full-drain
// just exposes that cost cleanly).
//
// Each sub-bench reports ns/op = wall-clock per round. Compare across rows /
// columns to read off:
//   - the row direction (varying shuffleWorkers at fixed mutationWorkers) shows whether
//     layering shuffle parallelism on top of mutation parallelism still helps once
//     mutation workers already saturate the cores.
//   - the column direction (varying mutationWorkers at fixed shuffleWorkers) shows the
//     mutation-level fan-out curve.
//
// Combinations span the (mut, shuf) product space:
//   - (1, W): pure shuffle parallelism — equivalent to fabsim eval's single-deck path
//     applied to N mutations sequentially.
//   - (M, 1): pure mutation parallelism — every mutation runs shuffle-single-threaded.
//   - (M, W): both layers, including some oversubscribed combinations (M*W > 8 cores)
//     to test whether stalling shuffle workers free up cores for sibling mutations.
//
// mutationSampleSize=8 keeps the benchmark bounded under ~2 minutes total at
// -benchtime=2x on an 8-physical-core machine, while still exposing the (mut=8, *)
// row's full parallelism.
func BenchmarkAnnealWorkerSweep(b *testing.B) {
	const (
		incoming = 7
		// unreachableBaseline keeps every mutation from being accepted, so the worker pool
		// drains every sampled mutation. Production rounds short-circuit on the first
		// improvement, but the per-mutation eval cost is the same.
		unreachableBaseline = 1_000_000.0
		// mutationSampleSize is the per-iteration drain depth. Sized to give every
		// mutation-parallel combo enough work to saturate its workers (mut=8 → 1 batch,
		// mut=4 → 2 batches, mut=1 → 8 sequential evals) without blowing past the bench
		// budget.
		mutationSampleSize = 8
	)
	loaded := loadRealDeck(b)
	if loaded == nil {
		b.Skip("mydecks/viserai_v4.json not found — saved deck needed for realistic bench")
	}
	all := AllMutations(loaded, 2, nil)
	if len(all) < mutationSampleSize {
		b.Fatalf("mutation pool size %d < sample size %d", len(all), mutationSampleSize)
	}
	mutations := all[:mutationSampleSize]

	combos := []struct {
		mut, shuf int
	}{
		{1, 1},  // single-threaded baseline (mutationWorkers=1 forces sequential)
		{1, 4},  // sequential mutations, shuffle workers=4
		{1, 8},  // sequential mutations, shuffle workers=8 (eval-mode shape)
		{2, 4},  // 2 mutation workers, 4 shuffle workers each
		{4, 2},  // 4 mutation workers, 2 shuffle workers each
		{4, 4},  // moderately oversubscribed: 16 total workers
		{8, 1},  // pure mutation parallelism (anneal's pre-shuffle-parallelism shape)
		{8, 2},  // mut=8 + shuffle=2, oversubscribed at 16
		{8, 4},  // mut=8 + shuffle=4, oversubscribed at 32
		{8, 8},  // mut=8 + shuffle=8, heavily oversubscribed at 64
		{16, 1}, // 16 mutation workers, single-threaded shuffle (SMT-saturated)
	}
	for _, c := range combos {
		c := c
		b.Run(fmt.Sprintf("mut=%d_shuf=%d", c.mut, c.shuf), func(b *testing.B) {
			b.ReportAllocs()
			b.ResetTimer()
			for n := 0; n < b.N; n++ {
				b.StopTimer()
				iterRNG := rand.New(rand.NewSource(42))
				b.StartTimer()
				_, _, _, found := IterateParallel(
					context.Background(), mutations, unreachableBaseline, 0, 0,
					0, incoming, c.mut, c.shuf,
					iterRNG.Int63(), nil, true,
				)
				if found {
					b.Fatalf("iter %d: unreachable baseline was beaten — bench setup is wrong", n)
				}
			}
		})
	}
}
