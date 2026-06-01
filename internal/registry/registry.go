package registry

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/deck"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/format"
)

// Registry is the deck-construction view of the production card / weapon roster: the pools
// the deck builder picks from with NotImplemented / Unplayable entries, format-banned cards,
// and cards illegal for the hero's class / talents filtered out. Satisfies
// internal/deck.Registry. Zero-sized; the underlying slices are package-level vars, so all
// instances are interchangeable. Callers pass Registry{} to deck.Random.
type Registry struct{}

// LegalCardsFor returns every registered card eligible for deck construction by hero h in
// format f: not flagged NotImplemented or Unplayable, legal in f (the banlist), and playable
// by h (class / talent — see heroCanPlay). Returns []deck.Card so deck.Random picks from the
// slice without a per-call conversion. Allocates O(legal card count) per call; deck.Random
// calls it once per sample run.
func (Registry) LegalCardsFor(f format.Format, h deck.Hero) []deck.Card {
	hc := asCardHero(h)
	base := legalCardsForFormat(f)
	kept := base[:0]
	for _, c := range base {
		if heroCanPlay(hc, c) {
			kept = append(kept, c)
		}
	}
	return kept
}

// LegalWeaponsFor returns every registered weapon not flagged NotImplemented or Unplayable,
// legal in f, and playable by h's class / talents.
func (Registry) LegalWeaponsFor(f format.Format, h deck.Hero) []deck.Weapon {
	hc := asCardHero(h)
	out := make([]deck.Weapon, 0, len(AllWeapons))
	for _, w := range AllWeapons {
		if isExcludedFromPool(w) || !f.IsCardLegal(w) {
			continue
		}
		// Weapons are card.Cards with a type-line, so the hero class / talent check applies.
		if wc, ok := w.(card.Card); ok && !heroCanPlay(hc, wc) {
			continue
		}
		out = append(out, w.(deck.Weapon))
	}
	return out
}

// legalCardsForFormat returns the registered cards passing the format-level filters only —
// markers plus the banlist, with no hero restriction. Backs LegalCardsFor (which then applies
// the hero check) and the card ranking (which is hero-agnostic).
func legalCardsForFormat(f format.Format) []deck.Card {
	out := make([]deck.Card, 0, len(cardsByID))
	for _, c := range cardsByID {
		if c == nil || isExcludedFromPool(c) || !f.IsCardLegal(c) {
			continue
		}
		out = append(out, c.(deck.Card))
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
