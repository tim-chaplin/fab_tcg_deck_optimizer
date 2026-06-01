package main

import (
	"math/rand"
	"os"
	"path/filepath"
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/deck"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/format"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/gameengine"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/hero/heroes"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/textio"
)

// BenchmarkEvalRealDeck mimics what `fabsim eval <deck> -incoming 5 -deep-shuffles 10000`
// does end-to-end: loads a saved 40-card Viserai deck and runs 10,000 shuffles through
// Evaluate. Sized to match real-eval workloads — internal/sim.BenchmarkEvaluate runs only
// 500 shuffles and uses Random() output, whose card mix and cache-hit profile can drift
// from production. Use this benchmark when reasoning about end-to-end fabsim eval speedups.
//
// Skips when mydecks/viserai_v4.json is absent so go test ./... still passes on a fresh
// checkout that doesn't carry the saved deck.
func BenchmarkEvalRealDeck(b *testing.B) {
	const (
		shuffles = 10000
		incoming = 5
	)
	path := findRepoFile(b, filepath.Join("mydecks", "viserai_v4.json"))
	if path == "" {
		b.Skip("mydecks/viserai_v4.json not found — saved decks are needed to run this bench")
	}
	data, err := os.ReadFile(path)
	if err != nil {
		b.Fatalf("read deck: %v", err)
	}
	loaded, _, err := textio.UnmarshalDeck(data)
	if err != nil {
		b.Fatalf("unmarshal deck: %v", err)
	}
	ev := sim.NewEvaluator()

	b.ReportAllocs()
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		b.StopTimer()
		d := loaded.Copy()
		evalRNG := rand.New(rand.NewSource(42))
		b.StartTimer()
		ev.Evaluate(d, shuffles, sim.Matchup{IncomingPhysicalDamage: incoming}, evalRNG)
	}
}

// BenchmarkEvalViseraiV5 mirrors `fabsim eval viserai_v5 -shuffles 5000 -incoming 7` end-to-end
// at half the production shuffle count to keep iterations bounded. Sized for direct comparison
// with BenchmarkEvalRealDeck (v4) so the gap between converged decks is visible on the same
// scale, and with BenchmarkEvalRandomDeck so we can read off which hot paths are deck-specific
// versus inherent to evaluation.
func BenchmarkEvalViseraiV5(b *testing.B) {
	const (
		shuffles = 5000
		incoming = 7
	)
	path := findRepoFile(b, filepath.Join("mydecks", "viserai_v5.json"))
	if path == "" {
		b.Skip("mydecks/viserai_v5.json not found")
	}
	data, err := os.ReadFile(path)
	if err != nil {
		b.Fatalf("read deck: %v", err)
	}
	loaded, _, err := textio.UnmarshalDeck(data)
	if err != nil {
		b.Fatalf("unmarshal deck: %v", err)
	}
	ev := sim.NewEvaluator()

	b.ReportAllocs()
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		b.StopTimer()
		d := loaded.Copy()
		evalRNG := rand.New(rand.NewSource(42))
		b.StartTimer()
		ev.Evaluate(d, shuffles, sim.Matchup{IncomingPhysicalDamage: incoming}, evalRNG)
	}
}

// BenchmarkEvalRandomDeck evaluates a random legal Viserai deck under the same shuffles /
// incoming as the v5 bench. Comparing the two profiles (which functions dominate where, how
// the alloc / time mix shifts) is the cleanest way to tell which hot spots are inherent to
// 40-card-Viserai evaluation versus deck-mix-specific symptoms of whatever makes v5 slow.
// Fixed seed for reproducibility.
func BenchmarkEvalRandomDeck(b *testing.B) {
	const (
		shuffles  = 5000
		incoming  = 7
		deckSize  = 40
		maxCopies = 2
		seed      = 42
	)
	gameengine.MaxDeckSize = deckSize
	rng := rand.New(rand.NewSource(seed))
	loaded := deck.Random(heroes.Viserai, format.SilverAge, deckSize, maxCopies, rng, registry.Registry{})
	ev := sim.NewEvaluator()

	b.ReportAllocs()
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		b.StopTimer()
		d := loaded.Copy()
		evalRNG := rand.New(rand.NewSource(seed))
		b.StartTimer()
		ev.Evaluate(d, shuffles, sim.Matchup{IncomingPhysicalDamage: incoming}, evalRNG)
	}
}

// findRepoFile walks up from the test's working directory looking for a relative path that
// exists, returning the first match or "" if none found within 5 parent hops. Lets the
// benchmark run from cmd/fabsim without hard-coding the repo root.
func findRepoFile(b *testing.B, rel string) string {
	dir, err := os.Getwd()
	if err != nil {
		b.Fatalf("getwd: %v", err)
	}
	for i := 0; i < 5; i++ {
		candidate := filepath.Join(dir, rel)
		if _, err := os.Stat(candidate); err == nil {
			return candidate
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return ""
		}
		dir = parent
	}
	return ""
}
