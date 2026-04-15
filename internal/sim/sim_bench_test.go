package sim

import (
	"math/rand"
	"testing"
)

func BenchmarkRun(b *testing.B) {
	deck := mixedDeck()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rng := rand.New(rand.NewSource(int64(i)))
		_ = Run(deck, 100, 4, rng)
	}
}
