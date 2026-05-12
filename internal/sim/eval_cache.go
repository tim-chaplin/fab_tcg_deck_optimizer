package sim

// Hand-eval cache. Keyed on the (hand multiset, runechantCarryover, arsenalCardIn, auras)
// tuple. Stores the winning partition's role assignments only; on a hit the chain replays
// against that one partition to rebuild the full TurnSummary, skipping the exponential
// partition search.
//
// Caching is gated on best.Cacheable=true: if any sibling partition read deck or graveyard
// via an accessor, the result depends on hidden state and can't be reused.
//
// Hand order doesn't affect Best's optimal result. The cache key sorts hand IDs into a
// canonical multiset so any ordering hits the same entry; on a hit, the cached role-multiset
// is remapped onto the new hand's ordering.
//
// Auras feed into the key as a sorted multiset of (SelfID, Count) pairs — Handler closures
// aren't hashable, but Handler behaviour is fully determined by SelfID.

import (
	"sync"
	"sync/atomic"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
)

// maxCachedHandSize caps how big a hand the cache will fingerprint. Adult heroes deal up
// to 4 cards plus the arsenal-in slot (5); 10 leaves headroom for any reveal handler that
// pads the hand at start of turn. Hands beyond this size skip the cache (treated as
// always-miss and never stored) — the cache key is the dealt hand, which is bounded by
// hero intelligence, so the cap doesn't apply to mid-chain growth.
const maxCachedHandSize = 10

// maxCachedWeapons caps the weapon slot count the cache fingerprints. Heroes carry at most
// 2 weapons (1H + 1H) plus the occasional off-hand item; 4 leaves headroom.
const maxCachedWeapons = 4

// maxCachedAuras caps how many aura triggers the cache will fingerprint. Real Viserai
// hands rarely have more than 2-3 simultaneous auras (Sigil of Silphidae + Malefic
// Incantation + the occasional charge counter); 8 leaves headroom for archetypes that
// stack more.
const maxCachedAuras = 8

// maxCachedItems caps how many items the cache will fingerprint. Matches maxCachedAuras
// since both fingerprint the same persistent-in-play surface; 8 leaves headroom for an
// archetype that stacks multiple token types alongside any card-based items.
const maxCachedItems = 8

// persistentCacheKey is one fingerprinted entry for a persistent-in-play permanent —
// shared by both Aura and Item entries since they have the same identifying surface.
// CardID identifies the originating card (or token kind — token IDs come from the
// engine's reserved range, see ids.RunechantTokenID etc.). Count is the per-entry
// counter (token copies, charges, fires remaining).
//
// Aura's TriggerType / OncePerTurn / FiredThisTurn aren't captured: no production card
// registers multiple triggers from the same source with different types or gates.
type persistentCacheKey struct {
	CardID ids.CardID
	Count  int
}

// evalCacheKey is the comparable map key for the hand-eval cache. Fixed-size arrays are
// zero-padded with explicit length fields so a shorter input can't collide with a longer one
// on its prefix. heroID and weaponIDs key the loadout (different loadouts produce different
// optimal partitions for the same hand).
//
// Matchup is intentionally NOT in the key — an Evaluator's lifetime spans calls at a constant
// Matchup in production. Tests that mix matchup values must use NewEvaluatorWithoutCache.
type evalCacheKey struct {
	handIDs        [maxCachedHandSize]ids.CardID
	weaponIDs      [maxCachedWeapons]ids.WeaponID
	auras          [maxCachedAuras]persistentCacheKey
	items          [maxCachedItems]persistentCacheKey
	handLen        int
	weaponLen      int
	auraLen        int
	itemLen        int
	heroID         ids.HeroID
	arsenalID      ids.CardID
	opponentMarked bool
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
// for this call." Hands arrive pre-sorted by Card.ID() — the deck-eval pipeline calls
// sortHandByID before each Best so the key is multiset-invariant by construction; we
// just copy the IDs into the fixed-size array. Aura entries are sorted by (SelfID,
// Count) so equivalent aura sets produce the same key regardless of trigger registration
// order. Weapon IDs are NOT sorted because the weapon order is stable across calls (same
// loadout, same slice header) and bestAttackWithWeapons enumerates weapon masks in slice
// order; reordering would still produce the same Value but the cached BestLine's
// swung-weapon names would drift, so we just preserve the input order. Matchup is
// omitted — see evalCacheKey doc.
func makeCacheKey(
	weapons []Weapon, hand []card.Card,
	prior Prior,
) (evalCacheKey, bool) {
	if len(hand) > maxCachedHandSize ||
		len(weapons) > maxCachedWeapons ||
		len(prior.Auras) > maxCachedAuras ||
		len(prior.Items) > maxCachedItems {
		return evalCacheKey{}, false
	}
	var key evalCacheKey
	key.handLen = len(hand)
	for i, c := range hand {
		key.handIDs[i] = c.ID()
	}
	key.weaponLen = len(weapons)
	for i, w := range weapons {
		key.weaponIDs[i] = w.ID()
	}
	// Persistent-in-play entries get insertion-sorted by (CardID, Count) so the cache key
	// stays multiset-invariant across registration order. Token kinds are distinguished
	// via their reserved CardID range (RunechantTokenID, PonderTokenID, …) so no separate
	// TokenType discriminator is needed.
	key.auraLen = len(prior.Auras)
	for i, t := range prior.Auras {
		insertPersistentEntry(key.auras[:i+1], persistentCacheKey{
			CardID: t.CardID(), Count: t.Count(),
		})
	}
	key.itemLen = len(prior.Items)
	for i, it := range prior.Items {
		insertPersistentEntry(key.items[:i+1], persistentCacheKey{
			CardID: it.CardID(), Count: it.Count(),
		})
	}
	if prior.Hero != nil {
		key.heroID = prior.Hero.ID()
	}
	if prior.Arsenal != nil {
		key.arsenalID = prior.Arsenal.ID()
	}
	key.opponentMarked = prior.OpponentMarked
	return key, true
}

// insertPersistentEntry places entry into dst in sorted order, shifting any greater
// elements right by one. dst is a slice over the key array's first (n+1) slots; the
// caller calls this once per source entry with i+1 as the slice length, so the prefix
// dst[:n] is already sorted and dst[n] is the new write position. Sort key:
// (SelfID, TokenType, Count) ascending.
func insertPersistentEntry(dst []persistentCacheKey, entry persistentCacheKey) {
	j := len(dst) - 1
	for j > 0 && persistentEntryLess(entry, dst[j-1]) {
		dst[j] = dst[j-1]
		j--
	}
	dst[j] = entry
}

// persistentEntryLess orders persistentCacheKey entries by CardID, then Count. Used by
// makeCacheKey's insertion sort.
func persistentEntryLess(a, b persistentCacheKey) bool {
	if a.CardID != b.CardID {
		return a.CardID < b.CardID
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
