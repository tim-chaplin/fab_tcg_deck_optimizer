package sim

import (
	"math/rand"
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/deck"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry"
)

// Measures cache hit / uncacheable rates on an anneal-shape workload: every single-card
// mutation of viserai_v4 evaluated against a shared cache. Sweeping the full mutation
// pool brings in cards beyond the base list, exposing realistic cacheable=false rates.
func TestEvalCache_AnnealShapeMeasurement(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	const (
		incoming         = 7
		shufflesPerEval  = 50
		mutationsSampled = 60
	)
	loaded := loadRealDeck(t)
	if loaded == nil {
		t.Skip("mydecks/viserai_v4.json not found")
	}
	mutations := deck.AllMutations(loaded, 2, true, registry.Registry{})
	if len(mutations) < mutationsSampled {
		t.Fatalf("mutation pool size %d < sample size %d", len(mutations), mutationsSampled)
	}
	// Shuffle deterministically so the same subset runs across baseline / branch.
	rng := rand.New(rand.NewSource(7))
	rng.Shuffle(len(mutations), func(i, j int) {
		mutations[i], mutations[j] = mutations[j], mutations[i]
	})
	mutations = mutations[:mutationsSampled]

	ev := NewEvaluator()
	for _, mut := range mutations {
		evalRNG := rand.New(rand.NewSource(99))
		ev.Evaluate(mut.Deck(), shufflesPerEval, Matchup{IncomingPhysicalDamage: incoming}, evalRNG)
	}

	cs := ev.CacheStats()
	total := cs.Hits + cs.Misses
	t.Logf("Anneal-shape measurement: %d mutations × %d shuffles each ≈ %d Best calls",
		mutationsSampled, shufflesPerEval, total)
	t.Logf("  hits:        %d (%.1f%%)", cs.Hits, 100*cs.HitRate())
	t.Logf("  misses:      %d (%.1f%%)", cs.Misses, 100*float64(cs.Misses)/float64(max1(total)))
	t.Logf("  uncacheable: %d (%.1f%% of misses, %.1f%% of total)",
		cs.Uncacheable, 100*float64(cs.Uncacheable)/float64(max1(cs.Misses)), 100*float64(cs.Uncacheable)/float64(max1(total)))
	t.Logf("  entries:     %d", cs.Entries)
}
