package main

import (
	"context"
	"math/rand"
	"os"
	"path/filepath"
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/deckio"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

// BenchmarkAnnealRoundOnViseraiV4 is BenchmarkAnnealRound anchored on a real saved deck
// (mydecks/viserai_v4.json) instead of a Random() Viserai pool. The unreachable baseline
// forces every sampled mutation to drain end-to-end.
//
// To refresh cmd/fabsim/default.pgo (the profile fabsim builds with via -pgo=auto):
//
//	go test -bench=BenchmarkAnnealRoundOnViseraiV4 -benchtime=3x -count=1 -run=^$ \
//	    -cpuprofile=cmd/fabsim/default.pgo ./cmd/fabsim/
//
// Skips when mydecks/viserai_v4.json is absent so go test ./... still passes on a fresh
// checkout that doesn't carry the saved deck.
func BenchmarkAnnealRoundOnViseraiV4(b *testing.B) {
	const (
		maxCopies           = 2
		incoming            = 7 // mid-game opponent swing
		unreachableBaseline = 1_000_000.0
		sampleSize          = 8
	)
	path := findRepoFile(b, filepath.Join("mydecks", "viserai_v4.json"))
	if path == "" {
		b.Skip("mydecks/viserai_v4.json not found — saved deck is needed to run this bench")
	}
	data, err := os.ReadFile(path)
	if err != nil {
		b.Fatalf("read deck: %v", err)
	}
	loaded, err := deckio.Unmarshal(data)
	if err != nil {
		b.Fatalf("unmarshal deck: %v", err)
	}
	baseline := sim.New(loaded.Hero, loaded.Weapons, loaded.Cards)
	all := sim.AllMutations(baseline, maxCopies, nil)
	if len(all) < sampleSize {
		b.Fatalf("mutation pool size %d < sample size %d", len(all), sampleSize)
	}
	mutations := all[:sampleSize]

	b.ReportAllocs()
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		b.StopTimer()
		iterRNG := rand.New(rand.NewSource(42))
		b.StartTimer()
		_, _, _, found := sim.IterateParallel(
			context.Background(), mutations, unreachableBaseline, 0, 0,
			0, incoming, 0, 0,
			iterRNG.Int63(), nil, true,
		)
		if found {
			b.Fatalf("iter %d: unreachable baseline was beaten — bench setup is wrong", n)
		}
	}
}
