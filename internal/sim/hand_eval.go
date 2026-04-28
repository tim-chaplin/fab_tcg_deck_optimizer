package sim

// Entry points for hand evaluation. Best / BestWithTriggers compute the optimal turn line
// for a given hand against an opponent attacking for incomingDamage. The Evaluator type is
// a no-op wrapper kept so concurrent callers can construct per-goroutine instances; every
// call allocates fresh scratch state.

import ()

// Best returns the optimal TurnSummary for the given hand against an opponent that will
// attack for incomingDamage on their next turn. Equipped weapons may be swung for their Cost
// if resources allow.
//
// Cards partition into five roles: Pitch (resource), Attack (played, may extend chain),
// Defend (blocks plus DR Plays), Held (stays in hand for next turn), Arsenal (moves to or
// stays in the arsenal slot at end of turn). Pitch resources split across attack / defense
// phases since resources don't carry between turns.
//
// arsenalCardIn is the card sitting in the arsenal slot at start of turn (nil if empty).
// runechantCarryover is the Runechant token count carrying in from the previous turn.
// TurnSummary.State.Runechants is the count at end of the chosen chain; feed it back as
// the next turn's carryover.
func Best(hero Hero, weapons []Weapon, hand []Card, incomingDamage int, deck []Card, runechantCarryover int, arsenalCardIn Card) TurnSummary {
	return sharedEvaluator.BestWithTriggers(hero, weapons, hand, incomingDamage, deck, runechantCarryover, arsenalCardIn, nil)
}

// BestWithTriggers is Best plus an explicit priorAuraTriggers input — the AuraTriggers
// carrying in from the previous turn. Mid-chain triggers (Malefic Incantation's
// TriggerAttackAction rune, etc.) may fire and contribute damage to this turn's Value.
func BestWithTriggers(hero Hero, weapons []Weapon, hand []Card, incomingDamage int, deck []Card, runechantCarryover int, arsenalCardIn Card, priorAuraTriggers []AuraTrigger) TurnSummary {
	return sharedEvaluator.BestWithTriggers(hero, weapons, hand, incomingDamage, deck, runechantCarryover, arsenalCardIn, priorAuraTriggers)
}

// Best is the method form of the package-level Best.
func (e *Evaluator) Best(hero Hero, weapons []Weapon, hand []Card, incomingDamage int, deck []Card, runechantCarryover int, arsenalCardIn Card) TurnSummary {
	return e.BestWithTriggers(hero, weapons, hand, incomingDamage, deck, runechantCarryover, arsenalCardIn, nil)
}

// BestWithTriggers is the method form of the package-level BestWithTriggers. Returns a
// TurnSummary with State.Log fully populated.
func (e *Evaluator) BestWithTriggers(hero Hero, weapons []Weapon, hand []Card, incomingDamage int, deck []Card, runechantCarryover int, arsenalCardIn Card, priorAuraTriggers []AuraTrigger) TurnSummary {
	return e.findBest(hero, weapons, hand, incomingDamage, deck, runechantCarryover, arsenalCardIn, priorAuraTriggers, false)
}

// BestWithTriggersSkipLog is BestWithTriggers without populating State.Log. Same Value and
// non-Log carry-state fields; State.Log comes back empty. The deck-eval loop uses this for
// every turn to skip the per-chain Log slice copy that dominates allocation bytes; only
// turns that become the new deck-best are replayed via BestWithTriggers to recover Log.
func (e *Evaluator) BestWithTriggersSkipLog(hero Hero, weapons []Weapon, hand []Card, incomingDamage int, deck []Card, runechantCarryover int, arsenalCardIn Card, priorAuraTriggers []AuraTrigger) TurnSummary {
	return e.findBest(hero, weapons, hand, incomingDamage, deck, runechantCarryover, arsenalCardIn, priorAuraTriggers, true)
}

// Evaluator caches per-goroutine scratch state across Best calls. The first call allocates
// an attackBufs sized for (handSize, weapons); subsequent calls with the same shape reuse it,
// avoiding ~12% of total bytes for a 10k-shuffle eval (newAttackBufs was the second-biggest
// allocator after the eval-time slice copies). Different shapes invalidate the cache and
// allocate fresh — fine for normal use because a single deck eval reuses one shape across
// every shuffle. Not safe for concurrent use; concurrent callers construct one Evaluator per
// goroutine (the parallel-shuffle path inside evaluateImpl already does this — each worker
// in the fan-out builds its own private Evaluator pointing to ev.cache).
//
// The hand-eval cache (cache field) memoizes the optimal partition for each unique
// (handMultiset, incomingDamage, runechantCarryover, arsenalCardIn) tuple seen during the
// Evaluator's lifetime. On a hit, Best skips the partition search and replays the chain
// against the cached BestLine; on a miss, the search runs and the result is stored when
// the chain didn't depend on hidden state. nil disables caching — used by benchmark and
// test helpers that want the from-scratch path.
type Evaluator struct {
	cachedBufs     *attackBufs
	cachedHandSize int
	cachedWeapons  []Weapon
	cache          *evalCache
	// numWorkers tells EvaluateWith how many goroutines to fan the shuffle loop across.
	// 0 or 1 runs sequentially in the calling goroutine, reusing cachedBufs as the per-
	// call scratch — the original single-threaded behaviour and the right default for
	// tests that want deterministic single-RNG runs. > 1 spawns N workers that share the
	// cache (which is RWMutex-protected) but each carry their own attackBufs scratch;
	// fabsim eval / anneal / compare use this path.
	numWorkers int
}

// NewEvaluator returns a fresh Evaluator with the hand-eval cache enabled and the shuffle
// loop running single-threaded. Safe for concurrent use across goroutines as long as each
// goroutine uses its own instance — internal scratch state is not synchronised.
func NewEvaluator() *Evaluator {
	return &Evaluator{cache: newEvalCache()}
}

// NewEvaluatorParallel returns an Evaluator that fans the shuffle loop across numWorkers
// goroutines. Each worker carries its own attackBufs scratch but they all share one
// thread-safe cache. fabsim eval / anneal / compare use this path; tests that want a
// deterministic single-RNG run construct via NewEvaluator instead.
func NewEvaluatorParallel(numWorkers int) *Evaluator {
	return &Evaluator{cache: newEvalCache(), numWorkers: numWorkers}
}

// NewEvaluatorWithoutCache returns a fresh Evaluator with the hand-eval cache disabled.
// Used for the from-scratch path in benchmarks and equivalence tests; production callers
// route through NewEvaluator / NewEvaluatorParallel.
func NewEvaluatorWithoutCache() *Evaluator {
	return &Evaluator{}
}

// ResetCache drops the cached entries while preserving the stats counters. Use between
// distinct decks when reusing one Evaluator across many of them (the iterate-mode worker
// pool's per-mutation loop): entries from one deck rarely help another — different card
// sets produce different hand multisets — so dropping them at deck boundaries caps memory
// at one-deck's-worth of entries. No-op when caching is disabled. Routes through the
// cache's reset method so the write lock guards against concurrent lookups in a parallel-
// shuffle worker pool.
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
		Entries:     entries,
	}
}

// sharedEvaluator backs the package-level Best — single-threaded callers don't need to
// construct their own. Caching is OFF here: the cache key omits incomingDamage on the
// premise that an Evaluator's lifetime spans calls at constant incomingDamage, which is
// true for production callers (each constructs their own Evaluator at a fixed incoming)
// but NOT for the test suite, which exercises Best at many different incoming values
// against the same package-level entry point. Tests that want cache behaviour construct
// their own Evaluator via NewEvaluator.
var sharedEvaluator = NewEvaluatorWithoutCache()
