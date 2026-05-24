package sim

import (
	"math/rand"
	"testing"
	"unsafe"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/deck"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry"
)

// Measures the in-RAM cost of the hand-eval cache: fill it via the anneal-shape
// workload, walk the live entries, sum struct + slice-payload sizes. Reports bytes
// per entry. Lets us reason about whether annealCacheCapacity (currently 200k) is
// appropriately sized for available memory.
//
// We size statically rather than diffing runtime.ReadMemStats because the per-Best
// scratch bufs + transient deck allocations dwarf the cache itself and any small
// counting error makes the delta noisy or negative. Walking entries gives a clean,
// deterministic byte count for the cache's own heap footprint.
func TestEvalCache_MemoryPerEntry(t *testing.T) {
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
	rng := rand.New(rand.NewSource(7))
	rng.Shuffle(len(mutations), func(i, j int) {
		mutations[i], mutations[j] = mutations[j], mutations[i]
	})
	mutations = mutations[:mutationsSampled]

	ev := NewEvaluator()
	for _, mut := range mutations {
		evalRNG := rand.New(rand.NewSource(99))
		ev.Evaluate(mut.Deck.Copy(), shufflesPerEval, Matchup{IncomingDamage: incoming}, evalRNG)
	}

	ev.cache.mu.RLock()
	defer ev.cache.mu.RUnlock()

	keyStaticBytes := int(unsafe.Sizeof(evalCacheKey{}))
	entryStaticBytes := int(unsafe.Sizeof(evalCacheEntry{}))
	// Go map per-entry overhead (bucket + hash + pointer chases) is roughly 48-96
	// bytes; we use 80 as a representative midpoint for a comparable-key map. Reported
	// separately so the static/dynamic split is visible.
	const mapOverheadPerEntry = 80

	totalKey, totalEntry, totalSlicePayload, totalStringPayload := 0, 0, 0, 0
	maxLine, maxAttack, maxPitch, maxDef, maxSwung := 0, 0, 0, 0, 0
	for k, e := range ev.cache.entries {
		_ = k
		totalKey += keyStaticBytes
		totalEntry += entryStaticBytes
		totalSlicePayload += len(e.line) * int(unsafe.Sizeof(card.CardAssignment{}))
		totalSlicePayload += len(e.attackOrder) * int(unsafe.Sizeof(playedCard{}))
		totalSlicePayload += len(e.defenders) * int(unsafe.Sizeof(playedCard{}))
		// pitchOrder + line slot for swungWeapons string-header.
		totalSlicePayload += len(e.pitchOrder) * int(unsafe.Sizeof((card.Card)(nil)))
		totalSlicePayload += len(e.swungWeapons) * int(unsafe.Sizeof(""))
		for _, s := range e.swungWeapons {
			totalStringPayload += len(s)
		}
		if len(e.line) > maxLine {
			maxLine = len(e.line)
		}
		if len(e.attackOrder) > maxAttack {
			maxAttack = len(e.attackOrder)
		}
		if len(e.pitchOrder) > maxPitch {
			maxPitch = len(e.pitchOrder)
		}
		if len(e.defenders) > maxDef {
			maxDef = len(e.defenders)
		}
		if len(e.swungWeapons) > maxSwung {
			maxSwung = len(e.swungWeapons)
		}
	}
	entries := len(ev.cache.entries)
	mapOverhead := entries * mapOverheadPerEntry
	totalBytes := totalKey + totalEntry + totalSlicePayload + totalStringPayload + mapOverhead
	bytesPerEntry := 0.0
	if entries > 0 {
		bytesPerEntry = float64(totalBytes) / float64(entries)
	}

	t.Logf("Cache memory measurement (%d entries):", entries)
	t.Logf("  key bytes (static):       %d (%.0f / entry)", totalKey, float64(totalKey)/float64(entries))
	t.Logf("  entry bytes (static):     %d (%.0f / entry)", totalEntry, float64(totalEntry)/float64(entries))
	t.Logf("  slice payload bytes:      %d (%.0f / entry)", totalSlicePayload, float64(totalSlicePayload)/float64(entries))
	t.Logf("  string payload bytes:     %d (%.0f / entry)", totalStringPayload, float64(totalStringPayload)/float64(entries))
	t.Logf("  map overhead (est ~80/e): %d (%.0f / entry)", mapOverhead, float64(mapOverhead)/float64(entries))
	t.Logf("  ------------------------")
	t.Logf("  total:                    %.1f MB (%.0f / entry)", float64(totalBytes)/(1024*1024), bytesPerEntry)
	t.Logf("  max line=%d attack=%d pitch=%d def=%d swung=%d", maxLine, maxAttack, maxPitch, maxDef, maxSwung)
	t.Logf("  extrapolated 200k cap:    %.1f MB", bytesPerEntry*200_000/(1024*1024))
	t.Logf("  extrapolated 1M cap:      %.1f MB", bytesPerEntry*1_000_000/(1024*1024))
	t.Logf("  extrapolated 5M cap:      %.1f MB", bytesPerEntry*5_000_000/(1024*1024))
}
