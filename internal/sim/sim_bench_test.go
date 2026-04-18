package sim

import (
	"context"
	"math/rand"
	"testing"

	fabdeck "github.com/tim-chaplin/fab-deck-optimizer/internal/deck"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/hero"
)

// BenchmarkIterateImprovements mimics iterate mode's inner loop: starting from a random Viserai
// deck, screen each mutation at shallow-shuffles, deep-confirm shallow passers, adopt the first
// confirmed improvement and restart. Stops once targetImprovements have been adopted, so the
// benchmark covers both the high-volume shallow screen path and the rare deep-confirm path in
// realistic proportions.
//
// Shallow/deep shuffle counts are smaller than production defaults to keep the benchmark runnable
// in single-digit seconds per iteration — the cost profile is the same, just compressed. Start
// seed is fixed so b.N iterations are comparable; rng threads through the whole hill climb so
// later mutations see different shuffle sequences than earlier ones, matching live iterate.
func BenchmarkIterateImprovements(b *testing.B) {
	const (
		deckSize        = 40
		maxCopies       = 2
		shallowShuffles = 100
		// deepShuffles is closer to the production default (10000) than a cheap smoke value — the
		// parallel-deep design only shows its real advantage when deep is meaningfully more
		// expensive than shallow, which is where iterate actually lives. Keeping it at 1000 hid
		// the regression the user hit in real iterate runs.
		deepShuffles       = 5000
		incoming           = 0
		targetImprovements = 2
	)
	for n := 0; n < b.N; n++ {
		b.StopTimer()
		rng := rand.New(rand.NewSource(42))
		best := fabdeck.Random(hero.Viserai{}, deckSize, maxCopies, rng)
		bestAvg := best.Evaluate(shallowShuffles, incoming, rng).Avg()
		b.StartTimer()

		improvements := 0
		for improvements < targetImprovements {
			mutations := fabdeck.AllMutations(best, maxCopies)
			d, avg, _, found := fabdeck.IterateParallel(
				context.Background(), mutations, bestAvg, shallowShuffles, deepShuffles, incoming, 0,
				rng.Int63(), rng, nil, nil,
			)
			if !found {
				b.Fatalf("local maximum reached before hitting %d improvements", targetImprovements)
			}
			bestAvg = avg
			best = d
			improvements++
		}
	}
}

