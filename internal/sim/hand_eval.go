package sim

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/deck"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/gameengine"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/weapon"
)

// Entry points for hand evaluation. Best computes the optimal turn line for a given hand
// against a Matchup. Evaluator caches per-goroutine scratch and the optional hand-eval cache.

// best returns the optimal TurnSummary for the given hand against matchup mp. Equipped
// weapons may be swung for their Cost if resources allow.
//
// Cards partition into five roles — Pitch / Attack / Defend / Held / Arsenal. Pitch
// resources split across attack / defense phases (resources don't carry between turns).
//
// master holds the previous turn's carryover: Hero, Arsenal, Auras, Items, Banished,
// Graveyard, OpponentMarked, plus matchup-derived IncomingDamage / ArcaneIncomingDamage.
// Per-turn ephemerals are ignored. TurnSummary.State is the post-chain GameState; the next
// turn's master comes from PrepareNextTurn on it.
//
// Package-private so callers route through EvalOneTurnForTesting.
func best(weapons []weapon.Weapon, hand []card.Card, d *deck.Deck, master *gameengine.GameState) TurnSummary {
	return sharedEvaluator.Best(weapons, hand, d, master)
}

// Best is the method form of the package-level Best.
func (e *Evaluator) Best(weapons []weapon.Weapon, hand []card.Card, d *deck.Deck, master *gameengine.GameState) TurnSummary {
	return e.findBest(weapons, hand, d, master)
}

// Evaluator caches per-goroutine scratch across Best calls. First call allocates attackBufs
// sized for (handSize, weapons); subsequent same-shape calls reuse it. Not safe for
// concurrent use — one Evaluator per goroutine.
//
// cache (nil to disable) memoizes the optimal partition per evalCacheKey. On a hit, Best
// replays the chain against the cached BestLine; on a miss, the search runs and the result
// is stored when the chain didn't depend on hidden state.
type Evaluator struct {
	cachedBufs     *attackBufs
	cachedHandSize int
	cachedWeapons  []weapon.Weapon
	cache          *evalCache
	// numWorkers controls Evaluate's shuffle-loop fan-out. 0 or 1 runs sequentially; >1
	// spawns N workers sharing the cache, each with its own attackBufs scratch.
	numWorkers int
}

// NewEvaluator returns a fresh Evaluator with its own private cache, shuffle loop
// single-threaded. One instance per goroutine.
func NewEvaluator() *Evaluator {
	return &Evaluator{cache: newEvalCache()}
}

// NewEvaluatorParallel returns an Evaluator fanning the shuffle loop across numWorkers
// goroutines, each with its own attackBufs sharing the Evaluator's private cache.
func NewEvaluatorParallel(numWorkers int) *Evaluator {
	return &Evaluator{cache: newEvalCache(), numWorkers: numWorkers}
}

// NewEvaluatorWithCache returns an Evaluator pointing at an existing shared Cache so a
// worker pool can pool lookup work. Single-threaded shuffle; set numWorkers on the result
// to layer shuffle parallelism on top.
func NewEvaluatorWithCache(c *Cache) *Evaluator {
	return &Evaluator{cache: c}
}

// NewEvaluatorWithoutCache returns a fresh Evaluator with the hand-eval cache disabled.
func NewEvaluatorWithoutCache() *Evaluator {
	return &Evaluator{}
}

// Cache is the thread-safe hand-eval cache shareable across Evaluators. Lookups take a
// read lock; store and reset take the write lock.
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

// CacheStats returns a snapshot of the Evaluator's cache counters; zero-valued when no
// cache. Counters read atomically; the entries-count read briefly takes the RLock to avoid
// racing a concurrent reset.
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
// assumes a constant Matchup per Evaluator and shared callers exercise many matchups here.
var sharedEvaluator = NewEvaluatorWithoutCache()
