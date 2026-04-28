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
// goroutine (iterate.go's worker pool already does this).
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
}

// NewEvaluator returns a fresh Evaluator with the hand-eval cache enabled. Safe for
// concurrent use across goroutines as long as each goroutine uses its own instance —
// internal scratch state is not synchronised.
func NewEvaluator() *Evaluator {
	return &Evaluator{cache: newEvalCache()}
}

// NewEvaluatorWithoutCache returns a fresh Evaluator with the hand-eval cache disabled.
// Used for the from-scratch path in benchmarks and equivalence tests; production callers
// route through NewEvaluator.
func NewEvaluatorWithoutCache() *Evaluator {
	return &Evaluator{}
}

// CacheStats returns a snapshot of the Evaluator's cache counters. Returns a zero-valued
// CacheStats when the Evaluator was constructed without a cache.
func (e *Evaluator) CacheStats() CacheStats {
	if e.cache == nil {
		return CacheStats{}
	}
	return CacheStats{
		Hits:          e.cache.hits,
		Misses:        e.cache.misses,
		SkipsTriggers: e.cache.skipsTriggers,
		Uncacheable:   e.cache.uncacheable,
		Entries:       len(e.cache.entries),
	}
}

// sharedEvaluator backs the package-level Best — single-threaded callers don't need to
// construct their own. Has caching enabled via NewEvaluator.
var sharedEvaluator = NewEvaluator()
