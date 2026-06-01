package sim

import (
	"math/rand"
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/deck"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/format"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/hero/heroes"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry"
)

// Measures the hit-rate delta from keeping the hand-eval cache warm across anneal
// rounds vs resetting it per round.
func TestEvalCache_PersistenceAcrossRounds(t *testing.T) {
	if testing.Short() {
		t.Skip("persistence measurement runs production shuffle counts; -short skips it")
	}
	const (
		deckSize          = 40
		maxCopies         = 2
		incoming          = 7
		shuffles          = 100
		rounds            = 5
		mutationsPerRound = 4
	)
	var baseline *deck.Deck
	if d := loadRealDeck(t); d != nil {
		baseline = d
		t.Log("using mydecks/viserai_v4.json")
	} else {
		setupRNG := rand.New(rand.NewSource(123))
		baseline = deck.Random(heroes.Viserai, format.SilverAge, deckSize, maxCopies, setupRNG, registry.Registry{})
		t.Log("using random Viserai deck")
	}

	// Both runs see the same per-round candidate decks; only the cache lifetime differs.
	roundDecks := makeRoundDecks(baseline, rounds, mutationsPerRound)

	freshTotal := evalCounters{}
	freshPerRound := make([]evalCounters, rounds)
	for r := 0; r < rounds; r++ {
		ev := NewEvaluator()
		for _, d := range roundDecks[r] {
			ev.Evaluate(d.Copy(), shuffles, Matchup{IncomingPhysicalDamage: incoming}, rand.New(rand.NewSource(99)))
		}
		stats := ev.CacheStats()
		freshPerRound[r] = evalCounters{hits: stats.Hits, misses: stats.Misses}
		freshTotal.hits += stats.Hits
		freshTotal.misses += stats.Misses
	}

	warmTotal := evalCounters{}
	warmPerRound := make([]evalCounters, rounds)
	warmEv := NewEvaluator()
	var priorHits, priorMisses int
	for r := 0; r < rounds; r++ {
		for _, d := range roundDecks[r] {
			warmEv.Evaluate(d.Copy(), shuffles, Matchup{IncomingPhysicalDamage: incoming}, rand.New(rand.NewSource(99)))
		}
		stats := warmEv.CacheStats()
		warmPerRound[r] = evalCounters{
			hits:   stats.Hits - priorHits,
			misses: stats.Misses - priorMisses,
		}
		warmTotal = evalCounters{hits: stats.Hits, misses: stats.Misses}
		priorHits, priorMisses = stats.Hits, stats.Misses
	}

	t.Logf("Per-round hit rates (fresh / warm):")
	for r := 0; r < rounds; r++ {
		t.Logf("  round %d: fresh=%.1f%% (%d/%d) warm=%.1f%% (%d/%d)",
			r+1,
			100*hitRate(freshPerRound[r]), freshPerRound[r].hits, freshPerRound[r].hits+freshPerRound[r].misses,
			100*hitRate(warmPerRound[r]), warmPerRound[r].hits, warmPerRound[r].hits+warmPerRound[r].misses,
		)
	}
	t.Logf("Cumulative: fresh=%.1f%% (%d/%d) warm=%.1f%% (%d/%d)",
		100*hitRate(freshTotal), freshTotal.hits, freshTotal.hits+freshTotal.misses,
		100*hitRate(warmTotal), warmTotal.hits, warmTotal.hits+warmTotal.misses,
	)
	t.Logf("Entries in warm cache at end: %d", warmEv.CacheStats().Entries)
}

// Tests that the bounded cache stays at capacity and increments evictions per
// displaced entry.
func TestEvalCacheBounded_EvictsAtCapacity(t *testing.T) {
	const capacity = 4
	c := newEvalCacheBounded(capacity)
	for i := 0; i < capacity*3; i++ {
		c.store(evalCacheKey{heroID: ids.HeroID(i + 1)}, evalCacheEntry{})
	}
	c.mu.RLock()
	size := len(c.entries)
	c.mu.RUnlock()
	if size != capacity {
		t.Errorf("entries = %d, want %d", size, capacity)
	}
	if got := c.evictions.Load(); got != int64(capacity*2) {
		t.Errorf("evictions = %d, want %d", got, capacity*2)
	}
}

type evalCounters struct {
	hits, misses int
}

func hitRate(c evalCounters) float64 {
	total := c.hits + c.misses
	if total == 0 {
		return 0
	}
	return float64(c.hits) / float64(total)
}

// makeRoundDecks builds the per-round mutation candidates the way anneal would: each
// round picks a fresh candidate (1-card swap off the round's base) and the accepted
// candidate threads into the next round's base. mutationsPerRound candidates per round
// approximate the worker-pool's per-round drain count.
func makeRoundDecks(baseline *deck.Deck, rounds, mutationsPerRound int) [][]*deck.Deck {
	out := make([][]*deck.Deck, rounds)
	mutGen := rand.New(rand.NewSource(7))
	current := baseline
	for r := 0; r < rounds; r++ {
		muts := deck.AllMutations(current, 2, false, registry.Registry{})
		mutGen.Shuffle(len(muts), func(i, j int) { muts[i], muts[j] = muts[j], muts[i] })
		picks := make([]*deck.Deck, 0, mutationsPerRound)
		for i := 0; i < mutationsPerRound && i < len(muts); i++ {
			picks = append(picks, muts[i].Deck())
		}
		out[r] = picks
		// Thread the first candidate into the next round's base — simulates anneal
		// accepting one mutation per round so each round's mutation pool centres on a
		// new baseline.
		if len(picks) > 0 {
			current = picks[0]
		}
	}
	return out
}
