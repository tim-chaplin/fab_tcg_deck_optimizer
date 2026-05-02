package main

import (
	"math/rand"
	"os"
	"path/filepath"
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/deckio"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
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
	loaded, err := deckio.Unmarshal(data)
	if err != nil {
		b.Fatalf("unmarshal deck: %v", err)
	}
	ev := sim.NewEvaluator()

	b.ReportAllocs()
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		b.StopTimer()
		d := sim.New(loaded.Hero, loaded.Weapons, loaded.Cards)
		evalRNG := rand.New(rand.NewSource(42))
		b.StartTimer()
		d.EvaluateWith(shuffles, incoming, 0, evalRNG, ev)
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
