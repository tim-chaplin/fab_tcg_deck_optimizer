package card

import (
	"sync/atomic"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
)

// DisplayName returns the human-readable identifier "<Name> [R/Y/B]" used in log lines,
// deck listings, and debug printouts. Pitch values outside 1-3 (weapons, items, hero
// cards, tokens with no printed pitch) fall through to the bare Name() — no color suffix
// to disambiguate. Result depends only on Card.ID(), so it's memoised in a per-ID table
// and the per-call string concat disappears on the hot path.
//
// The cache implementation lives here; callers across the codebase reach only DisplayName.
func DisplayName(c Card) string {
	id := c.ID()
	// Invalid (id == 0) is the slot test stubs and ad-hoc fakes share — skip the cache so
	// distinct stubs with the same zero ID don't return each other's strings.
	if id == ids.InvalidCard {
		return buildDisplayName(c)
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
// idx. Used both by WarmDisplayNameCache (to populate) and by DisplayName's miss path
// (to backfill on first sighting of an unregistered card — fakes / test stubs).
//
// Multiple goroutines computing the same entry race-safely converge on the first writer's
// string — every writer produces the same value, so reads after a race still match spec.
func displayNameSlow(c Card, id ids.CardID) string {
	out := buildDisplayName(c)
	displayNameCache[id].Store(&out)
	return out
}

// buildDisplayName produces the "<Name> [R/Y/B]" string. Pitch values outside 1-3 fall
// through to the bare Name() — no color suffix to disambiguate.
func buildDisplayName(c Card) string {
	switch c.Pitch() {
	case 1:
		return c.Name() + " [R]"
	case 2:
		return c.Name() + " [Y]"
	case 3:
		return c.Name() + " [B]"
	}
	return c.Name()
}

// WarmDisplayNameCache populates the DisplayName cache for every non-nil card in cards.
// Idempotent. The cards package init calls this once with the registry slice so the
// runtime hot path is pure cache reads — fakes/test stubs created without registration
// still work via DisplayName's lazy backfill.
func WarmDisplayNameCache(cards []Card) {
	for _, c := range cards {
		if c == nil {
			continue
		}
		displayNameSlow(c, c.ID())
	}
}
