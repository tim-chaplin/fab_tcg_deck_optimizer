package deck

import (
	"context"
	"math/rand"
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/heroes"
)

// BenchmarkAnnealRound mimics anneal's per-round workload: build the mutation pool for a
// random Viserai deck, then run a fixed-size SAMPLE of mutations (mutationSampleSize)
// through IterateParallel with an unreachable baseline (so the worker pool drains all
// sampled mutations end-to-end rather than short-circuiting). Reports per-iteration wall
// + allocations — the right granularity for measuring per-mutation-eval optimisations
// (buffer pools, chain pruning, hot-path tightening) on the workload anneal actually sees.
//
// The sample is the FIRST mutationSampleSize entries of AllMutations (deterministic
// across iterations), not a random pick — keeps the bench reproducible while still
// covering the mix of weapon-loadout / single-swap / pair-swap mutations the early head
// of the pool contains.
//
// Adaptive shuffles match the production default; the unreachable baseline forces every
// mutation's eval to run to completion (adaptive stop or cap), which is the worst-case
// throughput a round can pay. ~2.5s per iteration on a 16-thread Ryzen 5800H — fast
// enough for `-benchtime=3x -count=3` dev cycles while still amortising per-Best fixed
// costs across many distinct mutation states.
//
// Recommended invocation: `go test -bench=BenchmarkAnnealRound -benchtime=3x -count=3 -benchmem`.
func BenchmarkAnnealRound(b *testing.B) {
	const (
		deckSize  = 40
		maxCopies = 2
		// Non-zero incoming so the benchmark exercises the Defend-role partition branches.
		// 7 tracks a mid-game opponent swing typical of anneal sessions against Viserai.
		incoming = 7
		// Impossibly-high baseline keeps any mutation from being accepted, so the worker
		// pool drains every sampled mutation. Production rounds short-circuit on the
		// first improvement, but the per-mutation-eval cost is the same; a full-drain
		// sample just exposes that cost cleanly.
		unreachableBaseline = 1_000_000.0
		// mutationSampleSize is small enough to keep iteration time tractable, large
		// enough to amortise per-Best fixed costs across many distinct mutation states.
		mutationSampleSize = 50
	)

	setupRNG := rand.New(rand.NewSource(42))
	baseline := Random(heroes.Viserai{}, deckSize, maxCopies, setupRNG, nil)
	all := AllMutations(baseline, maxCopies, nil)
	if len(all) < mutationSampleSize {
		b.Fatalf("mutation pool size %d < sample size %d; bench setup needs a larger deck", len(all), mutationSampleSize)
	}
	mutations := all[:mutationSampleSize]

	b.ReportAllocs()
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		b.StopTimer()
		iterRNG := rand.New(rand.NewSource(42))
		b.StartTimer()
		// adaptive=true matches the production default; the shuffles arg is ignored in that mode.
		_, _, _, found := IterateParallel(
			context.Background(), mutations, unreachableBaseline, 0, 0,
			0, incoming, 0,
			iterRNG.Int63(), nil, true,
		)
		if found {
			b.Fatalf("iter %d: unreachable baseline was beaten — bench setup is wrong", n)
		}
	}
}

// BenchmarkIterateImprovements mimics iterate mode's inner loop: from a random Viserai deck,
// evaluate each mutation, adopt the first improvement and restart, until
// targetImprovements are adopted. Covers the round-scheduling + per-mutation-eval mix.
//
// Variance-control:
//   - targetImprovements is sized so each iteration does ~5 full rounds, amortising
//     per-round scheduling/GC blips.
//   - Shuffle count (5000) compresses the production default (10000) to keep each
//     iteration in single-digit seconds.
//   - Seed is fixed so every iteration walks the same mutation-pick sequence.
//
// Recommended invocation: `go test -bench=BenchmarkIterateImprovements -benchtime=5x -count=5`.
func BenchmarkIterateImprovements(b *testing.B) {
	const (
		deckSize  = 40
		maxCopies = 2
		shuffles  = 5000
		// Non-zero incoming so the benchmark exercises the Defend-role partition branches.
		// incoming=0 is rare in production use; 7 tracks a mid-game opponent swing that's
		// typical of iterate sessions against the classic archetype.
		incoming           = 7
		targetImprovements = 5
	)

	setupRNG := rand.New(rand.NewSource(42))
	baseline := Random(heroes.Viserai{}, deckSize, maxCopies, setupRNG, nil)
	baselineAvg := baseline.Evaluate(shuffles, incoming, setupRNG).Mean()

	for n := 0; n < b.N; n++ {
		b.StopTimer()
		iterRNG := rand.New(rand.NewSource(42))
		best := baseline
		bestAvg := baselineAvg
		b.StartTimer()

		for improvements := 0; improvements < targetImprovements; improvements++ {
			mutations := AllMutations(best, maxCopies, nil)
			d, avg, _, found := IterateParallel(
				context.Background(), mutations, bestAvg, 0, 0,
				shuffles, incoming, 0,
				iterRNG.Int63(), nil, false,
			)
			if !found {
				b.Fatalf("iter %d: local maximum reached at improvement %d of %d (baseline too good, try different seed)",
					n, improvements, targetImprovements)
			}
			bestAvg = avg
			best = d
		}
	}
}
