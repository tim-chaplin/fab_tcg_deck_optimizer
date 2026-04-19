package deck

import (
	"math/rand"
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/hand"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/hero"
)

// BenchmarkEvaluate drives the single-threaded hot path: generate a random Viserai deck and
// run a fixed number of shuffles through EvaluateWith. Deterministic (fixed seed, fixed evaluator)
// and fast enough (~2s per op at 500 shuffles) to support 10-count benchstat comparisons during
// optimization loops.
func BenchmarkEvaluate(b *testing.B) {
	const (
		deckSize  = 40
		maxCopies = 2
		shuffles  = 500
		incoming  = 0
	)
	setupRNG := rand.New(rand.NewSource(42))
	d := Random(hero.Viserai{}, deckSize, maxCopies, setupRNG, nil)
	ev := hand.NewEvaluator()
	evalRNG := rand.New(rand.NewSource(42))

	b.ReportAllocs()
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		b.StopTimer()
		hand.ClearMemo()
		d.Stats = Stats{}
		evalRNG = rand.New(rand.NewSource(42))
		b.StartTimer()
		d.EvaluateWith(shuffles, incoming, evalRNG, ev)
	}
}
