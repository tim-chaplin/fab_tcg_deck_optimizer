package sim

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/deck"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/gameengine"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/weapon"
)

// Entry points for hand evaluation. Best computes the optimal turn line for a given hand
// against the supplied Matchup. The Evaluator type caches per-goroutine scratch state and
// the optional hand-eval cache; concurrent callers construct one per goroutine.

// best returns the optimal TurnSummary for the given hand against the matchup mp.
// Equipped weapons may be swung for their Cost if resources allow.
//
// Cards partition into five roles: Pitch (resource), Attack (played, may extend chain),
// Defend (blocks plus DR Plays), Held (stays in hand for next turn), Arsenal (moves to or
// stays in the arsenal slot at end of turn). Pitch resources split across attack / defense
// phases since resources don't carry between turns.
//
// master holds start-of-turn carryover from the previous turn — Hero, Arsenal, Auras,
// Items, Banished, Graveyard, OpponentMarked, plus matchup-derived IncomingDamage /
// ArcaneIncomingDamage. Per-turn ephemerals (value, hand, pitched, log) are ignored.
// TurnSummary.State is the post-chain GameState; the next turn's master comes from
// calling PrepareNextTurn on it.
//
// Package-private so callers route through EvalOneTurnForTesting and exercise the same
// per-turn pipeline production uses.
func best(weapons []weapon.Weapon, hand []card.Card, d *deck.Deck, master *gameengine.GameState) TurnSummary {
	return sharedEvaluator.Best(weapons, hand, d, master)
}

// Best is the method form of the package-level Best.
func (e *Evaluator) Best(weapons []weapon.Weapon, hand []card.Card, d *deck.Deck, master *gameengine.GameState) TurnSummary {
	return e.findBest(weapons, hand, d, master)
}

// Evaluator caches per-goroutine scratch state across Best calls. The first call allocates
// an attackBufs sized for (handSize, weapons); subsequent calls with the same shape reuse
// it, saving ~12% of bytes on a 10k-shuffle eval. Different shapes invalidate the cache.
// Not safe for concurrent use; concurrent callers construct one Evaluator per goroutine.
//
// The hand-eval cache (cache field) memoizes the optimal partition per evalCacheKey. On a
// hit, Best skips the partition search and replays the chain against the cached BestLine;
// on a miss, the search runs and the result is stored when the chain didn't depend on
// hidden state. nil disables caching.
type Evaluator struct {
	cachedBufs     *attackBufs
	cachedHandSize int
	cachedWeapons  []weapon.Weapon
	cache          *evalCache
	// numWorkers controls Evaluate's shuffle-loop fan-out. 0 or 1 runs sequentially in
	// the calling goroutine; >1 spawns N workers that share the (RWMutex-protected) cache
	// but carry their own attackBufs scratch.
	numWorkers int
}

// NewEvaluator returns a fresh Evaluator with its own private cache and the shuffle loop
// running single-threaded. Safe for concurrent use across goroutines as long as each
// goroutine uses its own instance — internal scratch state is not synchronised.
func NewEvaluator() *Evaluator {
	return &Evaluator{cache: newEvalCache()}
}

// NewEvaluatorParallel returns an Evaluator that fans the shuffle loop across numWorkers
// goroutines, each carrying its own attackBufs scratch and sharing the Evaluator's
// private cache.
func NewEvaluatorParallel(numWorkers int) *Evaluator {
	return &Evaluator{cache: newEvalCache(), numWorkers: numWorkers}
}

// NewEvaluatorWithCache returns an Evaluator pointing at an existing shared Cache so a
// pool of workers can pool lookup work. numWorkers is 0 (single-threaded shuffle); set
// the field directly on the returned pointer to layer shuffle parallelism on top.
func NewEvaluatorWithCache(c *Cache) *Evaluator {
	return &Evaluator{cache: c}
}

// NewEvaluatorWithoutCache returns a fresh Evaluator with the hand-eval cache disabled.
func NewEvaluatorWithoutCache() *Evaluator {
	return &Evaluator{}
}

// Cache is the thread-safe hand-eval cache shared across multiple Evaluators. Lookups take
// a read lock (concurrent readers don't serialise); store and reset take the write lock.
type Cache = evalCache

// NewCache returns a fresh unbounded shared cache.
func NewCache() *Cache { return newEvalCache() }

// NewCacheBounded returns a shared cache that evicts a random entry on store once
// capacity is reached. capacity <= 0 disables eviction.
func NewCacheBounded(capacity int) *Cache { return newEvalCacheBounded(capacity) }

// ResetCache drops the cached entries while preserving the stats counters.
// No-op when caching is disabled.
func (e *Evaluator) ResetCache() {
	if e.cache != nil {
		e.cache.reset()
	}
}

// CacheStats returns a snapshot of the Evaluator's cache counters. Returns a zero-valued
// CacheStats when the Evaluator was constructed without a cache. Reads atomic counters
// without taking the entries lock; the entries-count read takes the read lock briefly to
// avoid racing a concurrent reset.
func (e *Evaluator) CacheStats() CacheStats {
	if e.cache == nil {
		return CacheStats{}
	}
	e.cache.mu.RLock()
	entries := len(e.cache.entries)
	e.cache.mu.RUnlock()
	return CacheStats{
		Hits:        int(e.cache.hits.Load()),
		Misses:      int(e.cache.misses.Load()),
		Uncacheable: int(e.cache.uncacheable.Load()),
		Evictions:   int(e.cache.evictions.Load()),
		Entries:     entries,
	}
}

// sharedEvaluator backs the package-level best — caching is OFF because the cache key
// assumes a constant Matchup per Evaluator and shared callers exercise many matchups
// against this entry point. Tests that want cache behaviour construct their own Evaluator.
var sharedEvaluator = NewEvaluatorWithoutCache()
