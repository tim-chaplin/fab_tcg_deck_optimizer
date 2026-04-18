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
		deckSize           = 40
		maxCopies          = 2
		shallowShuffles    = 100
		deepShuffles       = 1000
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
				rng.Int63(), rng, nil,
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

func BenchmarkRun(b *testing.B) {
	deck := mixedDeck()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rng := rand.New(rand.NewSource(int64(i)))
		_ = Run(deck, 100, 4, rng)
	}
}

// BenchmarkRunRealDeck exercises the simulator with a real random Viserai deck, so the profile
// reflects actual card implementations (Runechants, aura checks, Mauvrion grants, etc.) rather
// than the stubbed fake.Red/Blue attacks BenchmarkRun uses.
func BenchmarkRunRealDeck(b *testing.B) {
	rng := rand.New(rand.NewSource(42))
	d := fabdeck.Random(hero.Viserai{}, 40, 2, rng)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r := rand.New(rand.NewSource(int64(i)))
		_ = Run(d.Cards, 100, 4, r)
	}
}
