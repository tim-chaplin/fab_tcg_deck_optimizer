// Package optimizations memoises hot-path lookups whose results depend only on a card's
// stable identity. The card package exposes the bare builders (BuildDisplayName,
// BuildChainStepText) plus function-variable hooks (DisplayName, ChainStepText); this
// package installs cached implementations at init time and provides Warm* helpers the
// registry calls during startup so the runtime hot path is pure cache reads.
package optimizations

import (
	"sync/atomic"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
)

// bareDisplayName is the uncached card.DisplayName captured before init swaps in the
// memoised version. The slow path delegates to it so the build logic stays in one place
// (card.go).
var bareDisplayName = card.DisplayName

func init() {
	card.DisplayName = cachedDisplayName
}

// cachedDisplayName is the memoised DisplayName installed over card.DisplayName at init.
// Result depends only on Card.ID(), so it's keyed in a per-ID table sized for the full
// uint16 ID space. Invalid (id == 0) is the slot test stubs and ad-hoc fakes share, so we
// skip the cache there — distinct stubs with the same zero ID would otherwise return each
// other's strings.
func cachedDisplayName(c card.Card) string {
	id := c.ID()
	if id == ids.InvalidCard {
		return bareDisplayName(c)
	}
	if s := displayNameCache[id].Load(); s != nil {
		return *s
	}
	return displayNameSlow(c, id)
}

// displayNameCache memoises DisplayName results keyed by Card.ID. WarmDisplayNameCache
// fills every entry once at registry init; sized for the full uint16 ID space so lookups
// are direct bounds-checked array reads.
var displayNameCache [1 << 16]atomic.Pointer[string]

// displayNameSlow computes the DisplayName string and stores it in displayNameCache at
// id. Used both by WarmDisplayNameCache (to populate) and by cachedDisplayName's miss
// path (to backfill on first sighting of an unregistered card — fakes / test stubs).
//
// Multiple goroutines computing the same entry race-safely converge on the first writer's
// string — every writer produces the same value, so reads after a race still match spec.
func displayNameSlow(c card.Card, id ids.CardID) string {
	out := bareDisplayName(c)
	displayNameCache[id].Store(&out)
	return out
}

// WarmDisplayNameCache populates the DisplayName cache for every non-nil card in cards.
// Idempotent. The registry package's init calls this once with the registry slice so the
// runtime hot path is pure cache reads — fakes/test stubs created without registration
// still work via cachedDisplayName's lazy backfill.
func WarmDisplayNameCache(cards []card.Card) {
	for _, c := range cards {
		if c == nil {
			continue
		}
		displayNameSlow(c, c.ID())
	}
}
