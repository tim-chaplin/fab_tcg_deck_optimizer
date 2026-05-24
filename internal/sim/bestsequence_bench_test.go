package sim

import (
	"math/rand"
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/deck"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/hero/heroes"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry"
)

func BenchmarkEvaluate_BestSequence(b *testing.B) {
	const (
		deckSize  = 40
		maxCopies = 2
		incoming  = 7
		shuffles  = 200
	)
	setupRNG := rand.New(rand.NewSource(123))
	baseline := deck.Random(heroes.Viserai, deckSize, maxCopies, setupRNG, registry.Registry{})

	b.ReportAllocs()
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		b.StopTimer()
		d := baseline.Copy()
		rng := rand.New(rand.NewSource(99))
		ev := NewEvaluatorWithoutCache()
		b.StartTimer()
		ev.Evaluate(d, shuffles, Matchup{IncomingDamage: incoming}, rng)
	}
}
