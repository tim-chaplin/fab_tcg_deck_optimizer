package sim

import (
	"math/rand"
	"testing"

	fabdeck "github.com/tim-chaplin/fab-deck-optimizer/internal/deck"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/hero"
)

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
