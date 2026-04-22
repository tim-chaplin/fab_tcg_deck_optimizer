package deck

import (
	"context"
	"math/rand"
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/hand"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/hero"
)

// BenchmarkIterateImprovements mimics iterate mode's inner loop: starting from a random Viserai
// deck, screen each mutation at shallow-shuffles, deep-confirm shallow passers, adopt the first
// confirmed improvement and restart. Stops once targetImprovements have been adopted, so the
// benchmark covers both the high-volume shallow screen path and the rare deep-confirm path in
// realistic proportions.
//
// Variance-control:
//   - Every b.N iteration starts from a cold memo via hand.ClearMemo() so each sample measures
//     the same work; iteration 1+ otherwise inherit a warm cache and skew the distribution.
//   - targetImprovements is sized so each iteration does ~5 full rounds, making cold-cache
//     startup a small fraction of total time and averaging out per-round scheduling / GC blips.
//   - Shuffle counts (shallow=100, deep=5000) are a compressed-but-realistic version of
//     production defaults (100 / 10000). Each iteration runs in single-digit seconds so
//     `-benchtime=5x -count=5` gives a usable sample in about a minute.
//   - Seed is fixed so every b.N iteration walks the same mutation-pick sequence.
//
// Recommended invocation: `go test -bench=BenchmarkIterateImprovements -benchtime=5x -count=5`.
func BenchmarkIterateImprovements(b *testing.B) {
	const (
		deckSize           = 40
		maxCopies          = 2
		shallowShuffles    = 100
		deepShuffles       = 5000
		// Non-zero incoming so the benchmark exercises the Defend-role partition branches.
		// incoming=0 is rare in production use; 7 tracks a mid-game opponent swing that's
		// typical of iterate sessions against the classic archetype.
		incoming           = 7
		targetImprovements = 5
	)

	setupRNG := rand.New(rand.NewSource(42))
	baseline := Random(hero.Viserai{}, deckSize, maxCopies, setupRNG, nil)
	baselineAvg := baseline.Evaluate(shallowShuffles, incoming, setupRNG).Mean()

	for n := 0; n < b.N; n++ {
		b.StopTimer()
		hand.ClearMemo()
		iterRNG := rand.New(rand.NewSource(42))
		best := baseline
		bestAvg := baselineAvg
		b.StartTimer()

		for improvements := 0; improvements < targetImprovements; improvements++ {
			mutations := AllMutations(best, maxCopies, nil)
			d, avg, _, found := IterateParallel(
				context.Background(), mutations, bestAvg, 0,
				shallowShuffles, deepShuffles, incoming, 0,
				iterRNG.Int63(), nil, nil,
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
