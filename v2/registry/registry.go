package registry

import "github.com/tim-chaplin/fab-deck-optimizer/v2/deck"

// Registry is the deck-construction view of the production card / weapon roster: the pools
// the deck builder picks from with NotImplemented / Unplayable / NotSilverAgeLegal entries
// filtered out. Satisfies v2/deck.Registry. Zero-sized; the underlying slices are package-
// level vars, so all instances are interchangeable. Callers pass Registry{} to deck.Random.
type Registry struct{}

// LegalCards returns every registered card eligible for deck construction. Returns
// []deck.Card so deck.Random picks from the slice without a per-call conversion. Allocates
// O(legal card count) per call; deck.Random calls it once per sample run.
func (Registry) LegalCards() []deck.Card {
	out := make([]deck.Card, 0, len(cardsByID))
	for _, c := range cardsByID {
		if c == nil || isExcludedFromPool(c) {
			continue
		}
		out = append(out, c.(deck.Card))
	}
	return out
}

// LegalWeapons returns every registered weapon not flagged NotImplemented, Unplayable,
// or NotSilverAgeLegal.
func (Registry) LegalWeapons() []deck.Weapon {
	out := make([]deck.Weapon, 0, len(AllWeapons))
	for _, w := range AllWeapons {
		if isExcludedFromPool(w) {
			continue
		}
		out = append(out, w.(deck.Weapon))
	}
	return out
}

// isExcludedFromPool gates a card or weapon out of the deck-construction pool via the
// optional marker interfaces NotImplemented, Unplayable, and NotSilverAgeLegal. The
// markers are method-set tags with no type relation to Card / Weapon, so a single any-typed
// scan handles both inputs.
func isExcludedFromPool(x any) bool {
	if _, ok := x.(NotImplemented); ok {
		return true
	}
	if _, ok := x.(Unplayable); ok {
		return true
	}
	if _, ok := x.(NotSilverAgeLegal); ok {
		return true
	}
	return false
}
