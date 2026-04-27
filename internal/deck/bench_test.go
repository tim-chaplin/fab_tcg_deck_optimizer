package deck

import (
	"context"
	"math/rand"
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/hero"
)

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
	baseline := Random(hero.Viserai{}, deckSize, maxCopies, setupRNG, nil)
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
