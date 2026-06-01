package main

import (
	"context"
	"math/rand"
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/deck"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/format"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/gameengine"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/hero/heroes"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

// BenchmarkAnnealFromScratch replays the opening phase of a real `fabsim anneal` run: start
// from a random legal deck and drive several accept-and-cool rounds through the production
// helpers (buildRoundMutations + RunMutationRound + coolDown) at a low shuffle budget.
//
// Unlike BenchmarkAnnealRoundOnViseraiV4 — which drains a fixed, converged deck against an
// unreachable baseline and so only ever exercises one narrow card mix — this bench lets the
// deck churn: each accepted mutation swaps in cards from across the legal pool, so the eval
// hot paths see a wide variety of cards (auras, modal cards, weapons, on-hit riders) the way
// an early real run does. That makes it the representative target for profiling and PGO: code
// tuned only to viserai_v4's card mix doesn't necessarily help a from-scratch search.
//
// This is the workload PGO regenerates from. To refresh cmd/fabsim/default.pgo (the profile
// fabsim builds with via -pgo=auto):
//
//	go test -bench=BenchmarkAnnealFromScratch$ -benchtime=3x -count=1 -run=^$ \
//	    -cpuprofile=cmd/fabsim/default.pgo ./cmd/fabsim/
//
// Fixed seed + the production (mutationWorkers=1) shape make the accept trajectory
// deterministic on a given machine. -shuffles is pinned low (100) so a round resolves fast
// and the bench covers many distinct decks within its budget. Needs no saved deck file.
func BenchmarkAnnealFromScratch(b *testing.B) {
	const (
		incoming       = 7
		deckSize       = 40
		maxCopies      = 2
		shuffles       = 100
		rounds         = 25
		startTemp      = 1.0
		tempDecay      = 0.95
		minTemp        = 0.0
		minImprovement = 0.1
	)
	// Size the pool slots' graveyard/banished backings before any Evaluator is built.
	gameengine.MaxDeckSize = deckSize
	cfg := annealConfig{
		shuffles:  shuffles,
		matchup:   sim.Matchup{IncomingPhysicalDamage: incoming},
		deckSize:  deckSize,
		maxCopies: maxCopies,
	}

	b.ReportAllocs()
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		b.StopTimer()
		rng := rand.New(rand.NewSource(42))
		current := deck.Random(heroes.Viserai, format.SilverAge, deckSize, maxCopies, rng, registry.Registry{})
		currentAvg := sim.NewEvaluatorParallel(sim.DefaultWorkers()).
			Evaluate(current, shuffles, cfg.matchup, rng).Mean()
		cache := sim.NewCacheBounded(annealCacheCapacity)
		ranking := registry.NewRanking()
		temperature := startTemp
		b.StartTimer()

		for round := 0; round < rounds; round++ {
			mutations := buildRoundMutations(cfg, current, ranking, temperature, rng)
			d, _, avg, _, found := sim.RunMutationRound(
				context.Background(), mutations, currentAvg, temperature, minImprovement,
				cfg.shuffles, cfg.matchup, 0, 0,
				rng.Int63(), nil, cache,
			)
			if found {
				current = d
				currentAvg = avg
				temperature = coolDown(temperature, tempDecay, minTemp)
			}
		}
	}
}
