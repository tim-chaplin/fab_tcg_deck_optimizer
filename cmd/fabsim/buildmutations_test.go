package main

import (
	"math/rand"
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/deck"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/format"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/hero/heroes"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry"
)

// TestBuildRoundMutations_TemperatureWidensShuffle checks the phase-2 behavior: at T=0 the round's
// mutation order follows the ranking (non-decreasing rank), and at T>=1 it's fully shuffled.
func TestBuildRoundMutations_TemperatureWidensShuffle(t *testing.T) {
	rng := rand.New(rand.NewSource(1))
	d := deck.Random(heroes.Viserai, format.SilverAge, 40, 2, rng, registry.Registry{})
	cfg := annealConfig{maxCopies: 2}
	ranking := registry.NewRanking()

	ranked := func(ms []deck.Mutation) bool {
		for i := 1; i < len(ms); i++ {
			if mutationRank(ms[i], ranking) < mutationRank(ms[i-1], ranking) {
				return false
			}
		}
		return true
	}

	if cold := buildRoundMutations(cfg, d, ranking, 0, rand.New(rand.NewSource(2))); !ranked(cold) {
		t.Error("T=0 should keep strict ranked order (non-decreasing rank)")
	}
	if hot := buildRoundMutations(cfg, d, ranking, 1, rand.New(rand.NewSource(2))); ranked(hot) {
		t.Error("T=1 should fully shuffle — ranked order is astronomically unlikely")
	}
}
