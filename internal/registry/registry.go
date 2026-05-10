package registry

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/deck"
)

// Registry is the deck-construction view of the production card / weapon roster — the
// pools the deck builder picks from, with cards and weapons flagged NotImplemented or
// Unplayable already filtered out. Satisfies v2/deck.Registry directly. Zero-sized and
// stateless because the underlying slices are package-level vars; multiple instances are
// interchangeable. Callers pass Registry{} to deck.Random / deck.SanitizeNotImplemented.
type Registry struct{}

// LegalCards returns every registered card eligible for deck construction — every printing
// in the registry minus those flagged NotImplemented or Unplayable. Already typed as
// deck.Card so deck.Random picks from the slice without a per-call conversion. Allocation
// is O(legal card count) per call; the deck builders that hit this method call it once at
// the start of a Random / SanitizeNotImplemented invocation.
func (Registry) LegalCards() []deck.Card {
	out := make([]deck.Card, 0, len(cardsByID))
	for _, c := range cardsByID {
		if c == nil || isExcludedFromPool(c) {
			continue
		}
		out = append(out, c)
	}
	return out
}

// LegalWeapons returns every registered weapon not flagged NotImplemented or Unplayable.
func (Registry) LegalWeapons() []deck.Weapon {
	out := make([]deck.Weapon, 0, len(AllWeapons))
	for _, w := range AllWeapons {
		if isExcludedWeaponFromPool(w) {
			continue
		}
		out = append(out, w)
	}
	return out
}

// isExcludedFromPool / isExcludedWeaponFromPool gate cards and weapons out of the deck-
// construction pool. Three reasons exclude a printing today, expressed as optional marker
// interfaces:
//   - NotImplemented      — the simulator can't model the printed effect yet (registry-side).
//   - Unplayable          — the effect is so weak the optimizer wouldn't pick it (registry-side).
//   - NotSilverAgeLegal   — the printing is banned in Silver Age, the only constructed
//     format the deck builder targets today (sim-side, since it's a format rule).
func isExcludedFromPool(c sim.Card) bool {
	if _, ok := c.(NotImplemented); ok {
		return true
	}
	if _, ok := c.(Unplayable); ok {
		return true
	}
	if _, ok := c.(sim.NotSilverAgeLegal); ok {
		return true
	}
	return false
}

func isExcludedWeaponFromPool(w sim.Weapon) bool {
	if _, ok := w.(NotImplemented); ok {
		return true
	}
	if _, ok := w.(Unplayable); ok {
		return true
	}
	if _, ok := w.(sim.NotSilverAgeLegal); ok {
		return true
	}
	return false
}
