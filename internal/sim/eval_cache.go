package sim

// Hand-eval cache. Keyed on the (hand multiset, runechantCarryover, arsenalCardIn, auras)
// tuple. Stores the winning partition's role assignments; on a hit, the attack turn replays that
// partition to rebuild the full TurnSummary, skipping the exponential partition search.
//
// Caching is gated on best.Cacheable=true: if any sibling partition read deck or graveyard,
// the result depends on hidden state and can't be reused.
//
// Hand order doesn't affect Best, so the key sorts hand IDs into a canonical multiset; on a
// hit, the cached role-multiset is remapped onto the new hand's ordering. Auras key as a
// sorted (SelfID, Count) multiset — Handler closures aren't hashable, but Handler behaviour
// is fully determined by SelfID.

import (
	"sync"
	"sync/atomic"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/gameengine"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/weapon"
)

// maxCachedHandSize caps the hand size the cache fingerprints. Adult heroes deal 4 cards +
// arsenal-in (5); 10 leaves headroom for reveal handlers that pad the hand at start of turn.
// Hands beyond this skip the cache. The cap is on the dealt hand only — bounded by hero
// intelligence — not mid-attack-turn growth.
const maxCachedHandSize = 10

// maxCachedWeapons caps the weapon slot count the cache fingerprints. Heroes carry at most
// 2 weapons (1H + 1H) plus the occasional off-hand item; 4 leaves headroom.
const maxCachedWeapons = 4

// maxCachedAuras caps fingerprinted aura triggers. Real Viserai hands run 2-3 simultaneous
// auras; 8 leaves headroom for aura-stacking archetypes.
const maxCachedAuras = 8

// maxCachedItems caps fingerprinted items. Matches maxCachedAuras (same persistent-in-play
// surface).
const maxCachedItems = 8

// persistentCacheKey is one fingerprinted persistent-in-play permanent (Aura or Item).
// CardID identifies the source card (or token kind via the engine's reserved CardID range).
// Count is the per-entry counter (token copies, charges, fires remaining).
//
// Aura's TriggerType / OncePerTurn / FiredThisTurn aren't captured: no production card
// registers multiple triggers from the same source with different types or gates.
type persistentCacheKey struct {
	CardID ids.CardID
	Count  int
}

// evalCacheKey is the comparable map key. Fixed-size arrays carry explicit length fields so
// a shorter input can't collide with a longer one on its prefix. heroID + weaponIDs key the
// loadout. Matchup is NOT in the key — an Evaluator's lifetime spans calls at a constant
// Matchup; tests mixing matchups must use NewEvaluatorWithoutCache.
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

// playedCard is one resolved card in a cached solution: the card, its modal Mode, and
// whether it was played from arsenal. A cache replay seeds each attack step directly from
// these.
type playedCard struct {
	card        card.Card
	mode        int8
	fromArsenal bool
}

// cacheSolution is the in-flight winning solution before it stores into an evalCacheEntry.
// Slices alias attackBufs scratch, kept alloc-free across hand-offs.
type cacheSolution struct {
	attack    []playedCard
	pitch     []card.Card
	defenders []playedCard
}

// copyFrom rewrites s to match src, reusing s's backing slices.
func (s *cacheSolution) copyFrom(src *cacheSolution) {
	s.attack = append(s.attack[:0], src.attack...)
	s.pitch = append(s.pitch[:0], src.pitch...)
	s.defenders = append(s.defenders[:0], src.defenders...)
}

// reset empties s, keeping its backing slices for reuse.
func (s *cacheSolution) reset() {
	s.attack = s.attack[:0]
	s.pitch = s.pitch[:0]
	s.defenders = s.defenders[:0]
}

// evalCacheEntry is the cached winning solution — everything needed to rebuild a
// TurnSummary by replaying the winning line verbatim:
//
//   - line: BestLine roles (Card + Role + FromArsenal).
//   - swungWeapons: swung-weapon names, for TurnSummary display only.
//   - attackOrder: winning attacker permutation in resolution order, each with its chosen
//     modal Mode. Includes weapon / item ability cards at their attack-turn positions.
//   - pitchOrder: winning attack-phase pitch ordering — the pitch-pool pop sequence (also
//     identifying which pitch cards funded the attack vs defense phase).
//   - defenders: winning defender list, each with its blocker Mode.
//   - cardsRemovedFromDeck: cards the cached attack turn pulled off the deck. Replay refuses to
//     run when the caller's deck has fewer cards left.
type evalCacheEntry struct {
	line                 []card.CardAssignment
	swungWeapons         []string
	attackOrder          []playedCard
	pitchOrder           []card.Card
	defenders            []playedCard
	cardsRemovedFromDeck int
}

// evalCache holds cached Best results plus running stats. Thread-safe: mu guards entries
// (writes Lock, reads RLock); hits / misses / uncacheable are atomic so the lookup path
// bumps them without contending. A single cache can be shared across Evaluators — per-call
// scratch is goroutine-local, lookup and store are concurrency-safe.
//
// capacity is the entry cap before eviction kicks in (0 = unbounded). At cap, store evicts
// one random entry (Go map iteration order, no bookkeeping) so lookup stays RLock-only.
type evalCache struct {
	mu       sync.RWMutex
	entries  map[evalCacheKey]evalCacheEntry
	capacity int
	// hits / misses count lookup outcomes (one increment per Best). uncacheable counts
	// misses where the search ran but Cacheable=false, so the result wasn't stored —
	// useful for quantifying hidden-state reads. evictions counts entries dropped by the
	// capacity policy. Atomic so the lookup path bumps them without the map lock.
	hits, misses, uncacheable, evictions atomic.Int64
}

// CacheStats is the public snapshot of an Evaluator's cache counters. Hits + Misses equals
// total Best calls; Uncacheable is a subset of Misses (search ran but didn't store).
type CacheStats struct {
	Hits        int
	Misses      int
	Uncacheable int
	Evictions   int
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

// newEvalCache returns a fresh unbounded cache. Entries grows lazily on first store.
func newEvalCache() *evalCache {
	return &evalCache{}
}

// newEvalCacheBounded returns a cache that evicts a random entry on store when the
// entry count reaches capacity. capacity <= 0 disables eviction (equivalent to
// newEvalCache).
func newEvalCacheBounded(capacity int) *evalCache {
	if capacity < 0 {
		capacity = 0
	}
	return &evalCache{capacity: capacity}
}

// makeCacheKey builds the comparable cache key. Returns ok=false when hand / weapons / auras
// / items exceed their maxCached* caps — callers skip caching. Hands arrive pre-sorted by
// Card.ID() so the key is multiset-invariant. Aura / item entries insertion-sort by
// (CardID, Count) for the same reason. Weapon IDs are NOT sorted: weapon order is stable
// across calls, and reordering would drift the cached BestLine's swung-weapon names.
func makeCacheKey(
	weapons []weapon.Weapon, hand []card.Card,
	masterState *gameengine.GameState,
) (evalCacheKey, bool) {
	auras := masterState.Auras()
	items := masterState.Items()
	if len(hand) > maxCachedHandSize ||
		len(weapons) > maxCachedWeapons ||
		len(auras) > maxCachedAuras ||
		len(items) > maxCachedItems {
		return evalCacheKey{}, false
	}
	var key evalCacheKey
	// Test fakes return Invalid IDs which would collide on the same key. Bail on any
	// Invalid input — production entities always carry a unique non-zero ID.
	key.handLen = len(hand)
	for i, c := range hand {
		id := c.ID()
		if id == ids.InvalidCard {
			return evalCacheKey{}, false
		}
		key.handIDs[i] = id
	}
	key.weaponLen = len(weapons)
	for i, w := range weapons {
		id := w.ID()
		if id == ids.InvalidWeapon {
			return evalCacheKey{}, false
		}
		key.weaponIDs[i] = id
	}
	// Insertion-sort by (CardID, Count) so the key stays multiset-invariant across
	// registration order. Token kinds are distinguished by their reserved CardID range.
	key.auraLen = len(auras)
	for i, t := range auras {
		if t.CardID() == ids.InvalidCard {
			return evalCacheKey{}, false
		}
		insertPersistentEntry(key.auras[:i+1], persistentCacheKey{
			CardID: t.CardID(), Count: t.Count(),
		})
	}
	key.itemLen = len(items)
	for i, it := range items {
		if it.CardID() == ids.InvalidCard {
			return evalCacheKey{}, false
		}
		insertPersistentEntry(key.items[:i+1], persistentCacheKey{
			CardID: it.CardID(), Count: it.Count(),
		})
	}
	if h := masterState.Hero(); h != nil {
		if h.ID() == ids.InvalidHero {
			return evalCacheKey{}, false
		}
		key.heroID = h.ID()
	}
	if arsenal := masterState.Arsenal(); arsenal != nil {
		if arsenal.ID() == ids.InvalidCard {
			return evalCacheKey{}, false
		}
		key.arsenalID = arsenal.ID()
	}
	key.opponentMarked = masterState.OpponentMarked()
	return key, true
}

// insertPersistentEntry places entry into dst in sorted order, shifting greater elements
// right. dst[:n] is already sorted; dst[n] is the write position. Sort key: (CardID, Count).
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

// lookup returns the cached entry for key, or (zero, false) on miss. Doesn't bump stats —
// the caller does. Holds the read lock for the map access only; released before counter
// work so lookup contention stays minimal under shuffle-parallel readers.
func (c *evalCache) lookup(key evalCacheKey) (evalCacheEntry, bool) {
	c.mu.RLock()
	e, ok := c.entries[key]
	c.mu.RUnlock()
	return e, ok
}

// store inserts entry under key, lazily allocating the backing map under the write lock.
// Sibling workers racing to store the same key write identical data twice (harmless).
// At capacity, evicts one entry via Go's randomised map iteration order.
func (c *evalCache) store(key evalCacheKey, entry evalCacheEntry) {
	c.mu.Lock()
	if c.entries == nil {
		c.entries = make(map[evalCacheKey]evalCacheEntry)
	}
	if c.capacity > 0 && len(c.entries) >= c.capacity {
		if _, exists := c.entries[key]; !exists {
			for k := range c.entries {
				delete(c.entries, k)
				c.evictions.Add(1)
				break
			}
		}
	}
	c.entries[key] = entry
	c.mu.Unlock()
}

// reset drops the entries map under the write lock; stats counters survive so cumulative
// hit-rate measurement spans resets.
func (c *evalCache) reset() {
	c.mu.Lock()
	c.entries = nil
	c.mu.Unlock()
}
