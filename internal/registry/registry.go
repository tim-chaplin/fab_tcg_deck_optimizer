package registry

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/deck"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/format"
)

// Registry is the deck-construction view of the production card / weapon roster: the pools
// the deck builder picks from with NotImplemented / Unplayable entries and format-banned
// cards filtered out. Satisfies internal/deck.Registry. Zero-sized; the underlying slices are
// package-level vars, so all instances are interchangeable. Callers pass Registry{} to
// deck.Random.
type Registry struct{}

// LegalCards returns every registered card eligible for deck construction: not flagged
// NotImplemented or Unplayable and legal in Silver Age (the only format today). Returns
// []deck.Card so deck.Random picks from the slice without a per-call conversion. Allocates
// O(legal card count) per call; deck.Random calls it once per sample run.
func (Registry) LegalCards() []deck.Card {
	out := make([]deck.Card, 0, len(cardsByID))
	for _, c := range cardsByID {
		if c == nil || isExcludedFromPool(c) || !format.SilverAge.IsCardLegal(c) {
			continue
		}
		out = append(out, c.(deck.Card))
	}
	return out
}

// LegalWeapons returns every registered weapon not flagged NotImplemented or Unplayable and
// legal in Silver Age (the only format today).
func (Registry) LegalWeapons() []deck.Weapon {
	out := make([]deck.Weapon, 0, len(AllWeapons))
	for _, w := range AllWeapons {
		if isExcludedFromPool(w) || !format.SilverAge.IsCardLegal(w) {
			continue
		}
		out = append(out, w.(deck.Weapon))
	}
	return out
}

// isExcludedFromPool gates a card or weapon out of the deck-construction pool via the
// optional marker interfaces NotImplemented and Unplayable. The markers are method-set tags
// with no type relation to Card / Weapon, so a single any-typed scan handles both inputs.
// Format legality is applied separately by the callers (it needs only the card's name, not a
// marker).
func isExcludedFromPool(x any) bool {
	if _, ok := x.(NotImplemented); ok {
		return true
	}
	if _, ok := x.(Unplayable); ok {
		return true
	}
	return false
}
