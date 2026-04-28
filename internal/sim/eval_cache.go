package sim

// Hand-eval cache. Keyed on the (hand multiset, incomingDamage, runechantCarryover,
// arsenalCardIn) tuple — the inputs that fully determine Best's output WHEN the chain
// hasn't read hidden state. Cache stores only the winning partition's role assignments
// (BestLine + SwungWeapons); on a hit, the chain is replayed against that one partition
// to recover full TurnSummary state. That replay skips the partition search (the dominant
// cost — exponential in hand size) but still runs the per-leaf chain dispatcher to
// rebuild State / Log / Value, so the result matches a from-scratch Best call exactly.
//
// Two seed conditions for caching:
//   1. priorAuraTriggers must be empty. Carryover triggers add a hidden state input the
//      cache key doesn't model; serializing AuraTrigger.Handler closures isn't worth the
//      complexity, so we just skip caching when triggers are in play.
//   2. The chain must report Cacheable=true at end of search. If any sibling partition
//      read deck or graveyard via an accessor, the result depends on hidden shuffle /
//      prior-turn-graveyard contents and the cache can't safely reuse it.
//
// Hand order doesn't affect Best's optimal result (the search is exhaustive over role
// assignments). The cache key sorts hand IDs into a canonical multiset so the same hand
// in any order hits the same entry. On a hit, the cached BestLine's role-multiset is
// remapped onto the new hand's ordering so downstream consumers (deck-eval loop, printout)
// see roles attached to the right slice positions.

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
)

// maxCachedHandSize caps how big a hand the cache will fingerprint. Adult heroes deal up
// to 4 cards plus the arsenal-in slot (5) plus mid-turn-draw extensions; 24 covers the
// fattest plausible mid-chain growth. Hands beyond this size skip the cache (treated as
// always-miss and never stored) — should never happen in practice because the per-leaf
// hand size (the partition input) is bounded by the dealt hand size, which never exceeds
// hero intelligence.
const maxCachedHandSize = 24

// maxCachedWeapons caps the weapon slot count the cache fingerprints. Heroes carry at most
// 2 weapons (1H + 1H) plus the occasional off-hand item; 4 leaves headroom.
const maxCachedWeapons = 4

// evalCacheKey is the comparable map key for the hand-eval cache. handIDs is the sorted
// multiset of hand-card IDs zero-padded to maxCachedHandSize so the array is fixed-size
// (slices aren't comparable). handLen records the actual count so a 3-card hand and a
// 4-card hand whose first 3 IDs match can't collide on the trailing zero pad. heroID and
// weaponIDs[..weaponLen] capture the player's loadout — different heroes / weapon sets
// can produce different optimal partitions for the same hand, so they must key the cache.
//
// incomingDamage is intentionally NOT in the key. An Evaluator's lifetime spans calls at
// a constant incomingDamage in production (the iterate-mode worker pool, fabsim eval, and
// fabsim compare each fix incoming up front and reuse one Evaluator across many decks).
// Adding it would just bloat the key for no real-world hit-rate gain. Tests that mix
// incoming values across calls must use NewEvaluatorWithoutCache (sharedEvaluator already
// does for that reason) — otherwise a cached entry from incoming=A would silently apply
// to a query at incoming=B.
type evalCacheKey struct {
	handIDs            [maxCachedHandSize]ids.CardID
	weaponIDs          [maxCachedWeapons]ids.WeaponID
	handLen            int
	weaponLen          int
	runechantCarryover int
	heroID             ids.HeroID
	arsenalID          ids.CardID
}

// evalCacheEntry is the cached winning-partition shape. Stores only what's needed to
// replay: the BestLine roles (each Card paired with its Role + FromArsenal flag) and the
// list of swung weapon names. Value, State, and Log come from re-running the chain
// against the cached partition.
type evalCacheEntry struct {
	line         []CardAssignment
	swungWeapons []string
}

// evalCache holds per-Evaluator cached Best results plus the running stats counters the
// debug printout reads.
type evalCache struct {
	entries map[evalCacheKey]evalCacheEntry
	// hits / misses count cache lookup outcomes. skipsTriggers counts calls where
	// priorAuraTriggers was non-empty so the cache was bypassed entirely. uncacheable
	// counts misses where the search ran but Cacheable=false at the end so the result
	// wasn't stored — useful for quantifying how much hidden-state reading we'd need to
	// remove to bump the hit rate.
	hits, misses, skipsTriggers, uncacheable int
}

// CacheStats is the public snapshot of an Evaluator's cache counters, returned by
// Evaluator.CacheStats. Hits + Misses + SkipsTriggers is the total Best-call count;
// Uncacheable is a subset of Misses (the searches that ran but produced uncacheable
// results so weren't stored).
type CacheStats struct {
	Hits          int
	Misses        int
	SkipsTriggers int
	Uncacheable   int
	Entries       int
}

// HitRate returns hits / (hits+misses+skipsTriggers) as a fraction in [0, 1]. Returns 0
// when no calls have been made.
func (s CacheStats) HitRate() float64 {
	total := s.Hits + s.Misses + s.SkipsTriggers
	if total == 0 {
		return 0
	}
	return float64(s.Hits) / float64(total)
}

// PotentialHitRateWithTriggers projects what the hit rate would be if carryover-trigger
// support were added — counting SkipsTriggers as hypothetical hits. Useful for the
// "is it worth implementing trigger serialization?" decision.
func (s CacheStats) PotentialHitRateWithTriggers() float64 {
	total := s.Hits + s.Misses + s.SkipsTriggers
	if total == 0 {
		return 0
	}
	return float64(s.Hits+s.SkipsTriggers) / float64(total)
}

// newEvalCache returns a fresh cache. Entries grows lazily on first store.
func newEvalCache() *evalCache {
	return &evalCache{}
}

// makeCacheKey builds the comparable cache key from the inputs to Best. Returns ok=false
// when the hand exceeds maxCachedHandSize or the weapon slot count exceeds maxCachedWeapons;
// callers treat that as "skip caching for this call." The hand IDs are sorted ascending so
// the key is multiset-invariant — same hand in any draw order hits the same entry. Weapon
// IDs are NOT sorted because the weapon order is stable across calls (same loadout, same
// slice header) and bestAttackWithWeapons enumerates weapon masks in slice order; reordering
// would still produce the same Value but the cached BestLine's swung-weapon names would
// drift, so we just preserve the input order. incomingDamage is omitted from the key by
// the Evaluator-lifetime-constant assumption (see evalCacheKey doc).
func makeCacheKey(
	hero Hero, weapons []Weapon, hand []Card,
	runechantCarryover int, arsenalCardIn Card,
) (evalCacheKey, bool) {
	if len(hand) > maxCachedHandSize || len(weapons) > maxCachedWeapons {
		return evalCacheKey{}, false
	}
	var key evalCacheKey
	key.handLen = len(hand)
	// Insertion sort — hand size is small (typical 4-7) so the O(n^2) bound is faster than
	// sort.Slice's reflection-based path, and the loop stays inline-friendly without the
	// closure / interface dispatch sort.Slice introduces.
	for i, c := range hand {
		v := c.ID()
		j := i
		for j > 0 && key.handIDs[j-1] > v {
			key.handIDs[j] = key.handIDs[j-1]
			j--
		}
		key.handIDs[j] = v
	}
	key.weaponLen = len(weapons)
	for i, w := range weapons {
		key.weaponIDs[i] = w.ID()
	}
	key.runechantCarryover = runechantCarryover
	if hero != nil {
		key.heroID = hero.ID()
	}
	if arsenalCardIn != nil {
		key.arsenalID = arsenalCardIn.ID()
	}
	return key, true
}

// lookup returns the cached entry for key, or (zero, false) on miss. Doesn't bump
// counters — the caller does after confirming a hit / miss / skip.
func (c *evalCache) lookup(key evalCacheKey) (evalCacheEntry, bool) {
	if c.entries == nil {
		return evalCacheEntry{}, false
	}
	e, ok := c.entries[key]
	return e, ok
}

// store inserts entry under key, lazily allocating the backing map.
func (c *evalCache) store(key evalCacheKey, entry evalCacheEntry) {
	if c.entries == nil {
		c.entries = make(map[evalCacheKey]evalCacheEntry)
	}
	c.entries[key] = entry
}
