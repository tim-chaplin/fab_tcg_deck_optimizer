package sim

// Hand-eval cache. Keyed on the (hand multiset, runechantCarryover, arsenalCardIn, auras)
// tuple — the inputs that fully determine Best's output WHEN the chain hasn't read hidden
// state. Cache stores only the winning partition's role assignments (BestLine +
// SwungWeapons); on a hit, the chain is replayed against that one partition to recover
// full TurnSummary state. That replay skips the partition search (the dominant cost —
// exponential in hand size) but still runs the per-leaf chain dispatcher to rebuild
// State / Log / Value, so the result matches a from-scratch Best call exactly.
//
// Caching is gated on best.Cacheable=true at end of search. If any sibling partition read
// deck or graveyard via an accessor, the result depends on hidden shuffle / prior-turn-
// graveyard contents and the cache can't safely reuse it.
//
// Hand order doesn't affect Best's optimal result (the search is exhaustive over role
// assignments). The cache key sorts hand IDs into a canonical multiset so the same hand
// in any order hits the same entry. On a hit, the cached BestLine's role-multiset is
// remapped onto the new hand's ordering so downstream consumers (deck-eval loop, printout)
// see roles attached to the right slice positions.
//
// Aura entries (priorAuraTriggers carrying in from the previous turn) feed into the key
// as a sorted multiset of (SelfID, Count) pairs — the trigger Handler closures aren't
// hashable, but Handler behaviour is fully determined by SelfID (each card type's Play
// always registers the same handler logic), so SelfID + Count captures everything that
// affects chain output. See auraCacheKey for caveats.

import (
	"sync"
	"sync/atomic"

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

// maxCachedAuras caps how many aura triggers the cache will fingerprint. Real Viserai
// hands rarely have more than 2-3 simultaneous auras (Sigil of Silphidae + Malefic
// Incantation + the occasional charge counter); 8 leaves headroom for archetypes that
// stack more.
const maxCachedAuras = 8

// auraCacheKey is one fingerprinted entry in the evalCacheKey.auras array. SelfID is the
// CardID of the aura the trigger belongs to — Handler closures aren't comparable, but
// each card type's Play always registers the same handler logic, so SelfID determines
// the per-fire behaviour. Count is the remaining-fires counter.
//
// Caveats this minimal shape doesn't capture today: Type (TriggerStartOfTurn vs
// TriggerAttackAction), OncePerTurn / FiredThisTurn flags. No production card registers
// multiple triggers from the same Self with different types or gates, so the simpler
// (SelfID, Count) tuple is enough — but if a future card needs to disambiguate, this is
// the place to extend the key.
type auraCacheKey struct {
	SelfID ids.CardID
	Count  int
}

// evalCacheKey is the comparable map key for the hand-eval cache. handIDs is the sorted
// multiset of hand-card IDs zero-padded to maxCachedHandSize so the array is fixed-size
// (slices aren't comparable). handLen records the actual count so a 3-card hand and a
// 4-card hand whose first 3 IDs match can't collide on the trailing zero pad. heroID and
// weaponIDs[..weaponLen] capture the player's loadout — different heroes / weapon sets
// can produce different optimal partitions for the same hand, so they must key the cache.
// auras[..auraLen] is the sorted multiset of (SelfID, Count) tuples for the priorAura-
// Triggers passed in — same fixed-size-array trick.
//
// incomingDamage and arcaneIncomingDamage are intentionally NOT in the key. An Evaluator's
// lifetime spans calls at constant matchup parameters in production (the iterate-mode worker
// pool, fabsim eval, and fabsim compare each fix incoming up front and reuse one Evaluator
// across many decks). Adding either would just bloat the key for no real-world hit-rate gain.
// Tests that mix matchup values across calls must use NewEvaluatorWithoutCache (sharedEvaluator
// already does for that reason) — otherwise a cached entry from one matchup would silently
// apply to a query under another.
type evalCacheKey struct {
	handIDs            [maxCachedHandSize]ids.CardID
	weaponIDs          [maxCachedWeapons]ids.WeaponID
	auras              [maxCachedAuras]auraCacheKey
	handLen            int
	weaponLen          int
	auraLen            int
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

// evalCache holds cached Best results plus the running stats counters the debug printout
// reads. Thread-safe: mu guards entries (map writes use Lock; reads use RLock); hits /
// misses / uncacheable are atomic so the lookup hot path bumps them without contending
// with the entries lock. A single evalCache can be shared across multiple Evaluators —
// each Evaluator's per-call scratch (attackBufs) is goroutine-local but the cache lookup
// and store are concurrency-safe, which lets a shuffle-parallel worker pool reuse the
// same Cache across all workers.
type evalCache struct {
	mu      sync.RWMutex
	entries map[evalCacheKey]evalCacheEntry
	// hits / misses count cache lookup outcomes (every Best call increments exactly one).
	// uncacheable counts misses where the search ran but Cacheable=false at the end so
	// the result wasn't stored — useful for quantifying how much hidden-state reading
	// we'd need to remove to bump the hit rate further. Atomic so the lookup path bumps
	// them without taking the map lock.
	hits, misses, uncacheable atomic.Int64
}

// CacheStats is the public snapshot of an Evaluator's cache counters, returned by
// Evaluator.CacheStats. Hits + Misses is the total Best-call count; Uncacheable is a
// subset of Misses (the searches that ran but produced uncacheable results so weren't
// stored).
type CacheStats struct {
	Hits        int
	Misses      int
	Uncacheable int
	Entries     int
}

// HitRate returns hits / (hits+misses) as a fraction in [0, 1]. Returns 0 when no calls
// have been made.
func (s CacheStats) HitRate() float64 {
	total := s.Hits + s.Misses
	if total == 0 {
		return 0
	}
	return float64(s.Hits) / float64(total)
}

// newEvalCache returns a fresh cache. Entries grows lazily on first store.
func newEvalCache() *evalCache {
	return &evalCache{}
}

// makeCacheKey builds the comparable cache key from the inputs to Best. Returns ok=false
// when the hand exceeds maxCachedHandSize, the weapon slot count exceeds maxCachedWeapons,
// or the carryover-aura count exceeds maxCachedAuras; callers treat that as "skip caching
// for this call." The hand IDs are sorted ascending so the key is multiset-invariant —
// same hand in any draw order hits the same entry. Aura entries are sorted the same way
// (by SelfID then Count) so equivalent aura sets produce the same key regardless of
// trigger registration order. Weapon IDs are NOT sorted because the weapon order is
// stable across calls (same loadout, same slice header) and bestAttackWithWeapons
// enumerates weapon masks in slice order; reordering would still produce the same Value
// but the cached BestLine's swung-weapon names would drift, so we just preserve the
// input order. incomingDamage and arcaneIncomingDamage are omitted from the key by the
// Evaluator-lifetime-constant assumption (see evalCacheKey doc).
func makeCacheKey(
	hero Hero, weapons []Weapon, hand []Card,
	runechantCarryover int, arsenalCardIn Card,
	priorAuraTriggers []AuraTrigger,
) (evalCacheKey, bool) {
	if len(hand) > maxCachedHandSize ||
		len(weapons) > maxCachedWeapons ||
		len(priorAuraTriggers) > maxCachedAuras {
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
	// Aura entries: insertion-sort by (SelfID, Count) so the key is multiset-invariant
	// across trigger registration order. The aura set is small (typically 0-3) so the
	// O(n^2) cost is negligible.
	key.auraLen = len(priorAuraTriggers)
	for i, t := range priorAuraTriggers {
		entry := auraCacheKey{SelfID: t.Self.ID(), Count: t.Count}
		j := i
		for j > 0 && auraEntryLess(entry, key.auras[j-1]) {
			key.auras[j] = key.auras[j-1]
			j--
		}
		key.auras[j] = entry
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

// auraEntryLess orders auraCacheKey entries by SelfID first, then Count. Used by
// makeCacheKey's insertion sort so the auras[..auraLen] prefix is canonically ordered.
func auraEntryLess(a, b auraCacheKey) bool {
	if a.SelfID != b.SelfID {
		return a.SelfID < b.SelfID
	}
	return a.Count < b.Count
}

// lookup returns the cached entry for key, or (zero, false) on miss. Doesn't bump the
// stats counters — the caller does after confirming a hit / miss. Holds the read lock for
// the map access only; the lock is released before the caller bumps counters or runs any
// further work, which keeps lookup contention minimal under a parallel-shuffle worker
// pool reading the same cache.
func (c *evalCache) lookup(key evalCacheKey) (evalCacheEntry, bool) {
	c.mu.RLock()
	e, ok := c.entries[key]
	c.mu.RUnlock()
	return e, ok
}

// store inserts entry under key, lazily allocating the backing map. Takes the write lock;
// concurrent miss-then-store from sibling workers may both store the same key (the second
// write overwrites identical data) but never observes a partially-constructed map.
func (c *evalCache) store(key evalCacheKey, entry evalCacheEntry) {
	c.mu.Lock()
	if c.entries == nil {
		c.entries = make(map[evalCacheKey]evalCacheEntry)
	}
	c.entries[key] = entry
	c.mu.Unlock()
}

// reset drops the entries map (lazy realloc on next store) under the write lock so a
// concurrent reader/writer can never see a half-cleared state. Stats counters survive
// the reset — see Evaluator.ResetCache for the rationale.
func (c *evalCache) reset() {
	c.mu.Lock()
	c.entries = nil
	c.mu.Unlock()
}
