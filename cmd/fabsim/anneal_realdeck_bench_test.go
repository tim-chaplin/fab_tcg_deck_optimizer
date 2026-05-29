package main

import (
	"context"
	"math/rand"
	"os"
	"path/filepath"
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/deck"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/textio"
)

// BenchmarkAnnealRoundOnViseraiV4 mimics anneal's per-round workload anchored on
// mydecks/viserai_v4.json: build the mutation pool, run the first sampleSize mutations
// through RunMutationRound against an unreachable baseline so the worker pool drains every
// sampled mutation end-to-end. This is the gold-standard anneal bench — the workload PGO
// profiles target and the canonical reference for measuring per-mutation-eval changes.
//
// To refresh cmd/fabsim/default.pgo (the profile fabsim builds with via -pgo=auto):
//
//	go test -bench=BenchmarkAnnealRoundOnViseraiV4$ -benchtime=3x -count=1 -run=^$ \
//	    -cpuprofile=cmd/fabsim/default.pgo ./cmd/fabsim/
//
// Skips when mydecks/viserai_v4.json is absent so go test ./... still passes on a fresh
// checkout that doesn't carry the saved deck.
//
// For PR-validation benchstat, use the much faster BenchmarkAnnealRoundOnViseraiV4_Quick
// sibling.
func BenchmarkAnnealRoundOnViseraiV4(b *testing.B) {
	const (
		maxCopies           = 2
		incoming            = 7 // mid-game opponent swing
		unreachableBaseline = 1_000_000.0
		sampleSize          = 8
		shuffles            = 10000
	)
	path := findRepoFile(b, filepath.Join("mydecks", "viserai_v4.json"))
	if path == "" {
		b.Skip("mydecks/viserai_v4.json not found — saved deck is needed to run this bench")
	}
	data, err := os.ReadFile(path)
	if err != nil {
		b.Fatalf("read deck: %v", err)
	}
	loaded, _, err := textio.UnmarshalDeck(data)
	if err != nil {
		b.Fatalf("unmarshal deck: %v", err)
	}
	baseline := loaded.Copy()
	all := deck.AllMutations(baseline, maxCopies, true, registry.Registry{})
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
		_, _, _, _, found := sim.RunMutationRound(
			context.Background(), mutations, unreachableBaseline, 0, 0,
			shuffles, sim.Matchup{IncomingPhysicalDamage: incoming}, 0, 0,
			iterRNG.Int63(), nil, nil,
		)
		if found {
			b.Fatalf("iter %d: unreachable baseline was beaten — bench setup is wrong", n)
		}
	}
}

// Measures the gain from a persistent hand-eval cache across consecutive anneal rounds
// on viserai_v4.
func BenchmarkAnnealMultiRoundOnViseraiV4(b *testing.B) {
	const (
		maxCopies           = 2
		incoming            = 7
		unreachableBaseline = 1_000_000.0
		sampleSize          = 2
		shuffles            = 100
		rounds              = 5
		capacity            = annealCacheCapacity
	)
	path := findRepoFile(b, filepath.Join("mydecks", "viserai_v4.json"))
	if path == "" {
		b.Skip("mydecks/viserai_v4.json not found")
	}
	data, err := os.ReadFile(path)
	if err != nil {
		b.Fatalf("read deck: %v", err)
	}
	loaded, _, err := textio.UnmarshalDeck(data)
	if err != nil {
		b.Fatalf("unmarshal deck: %v", err)
	}
	all := deck.AllMutations(loaded.Copy(), maxCopies, true, registry.Registry{})
	if len(all) < sampleSize {
		b.Fatalf("mutation pool size %d < sample size %d", len(all), sampleSize)
	}
	mutations := all[:sampleSize]

	runRound := func(cache *sim.Cache, seed int64) {
		_, _, _, _, _ = sim.RunMutationRound(
			context.Background(), mutations, unreachableBaseline, 0, 0,
			shuffles, sim.Matchup{IncomingPhysicalDamage: incoming}, 0, 0,
			seed, nil, cache,
		)
	}

	b.Run("fresh", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for n := 0; n < b.N; n++ {
			b.StopTimer()
			iterRNG := rand.New(rand.NewSource(42))
			b.StartTimer()
			for r := 0; r < rounds; r++ {
				runRound(nil, iterRNG.Int63())
			}
		}
	})
	b.Run("persistent", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for n := 0; n < b.N; n++ {
			b.StopTimer()
			iterRNG := rand.New(rand.NewSource(42))
			cache := sim.NewCacheBounded(capacity)
			b.StartTimer()
			for r := 0; r < rounds; r++ {
				runRound(cache, iterRNG.Int63())
			}
		}
	})
}

// BenchmarkAnnealRoundOnViseraiV4_Quick is the PR-validation sibling of
// BenchmarkAnnealRoundOnViseraiV4: same anneal-round workload on viserai_v4 but with a
// 2-mutation sample and a fixed 200-shuffle Evaluate, so a -count=10 -benchtime=3x benchstat
// finishes in well under a minute. Hits the same hot paths
// (attack-turn runner, partition search, cache) — a regression in either shows up here just as
// reliably as in the full bench, just with a smaller per-iteration footprint.
func BenchmarkAnnealRoundOnViseraiV4_Quick(b *testing.B) {
	const (
		maxCopies           = 2
		incoming            = 7
		unreachableBaseline = 1_000_000.0
		sampleSize          = 2
		shuffles            = 200
	)
	path := findRepoFile(b, filepath.Join("mydecks", "viserai_v4.json"))
	if path == "" {
		b.Skip("mydecks/viserai_v4.json not found — saved deck is needed to run this bench")
	}
	data, err := os.ReadFile(path)
	if err != nil {
		b.Fatalf("read deck: %v", err)
	}
	loaded, _, err := textio.UnmarshalDeck(data)
	if err != nil {
		b.Fatalf("unmarshal deck: %v", err)
	}
	baseline := loaded.Copy()
	all := deck.AllMutations(baseline, maxCopies, true, registry.Registry{})
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
		_, _, _, _, found := sim.RunMutationRound(
			context.Background(), mutations, unreachableBaseline, 0, 0,
			shuffles, sim.Matchup{IncomingPhysicalDamage: incoming}, 0, 0,
			iterRNG.Int63(), nil, nil,
		)
		if found {
			b.Fatalf("iter %d: unreachable baseline was beaten — bench setup is wrong", n)
		}
	}
}
