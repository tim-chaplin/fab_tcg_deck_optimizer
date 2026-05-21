package registry

import "github.com/tim-chaplin/fab-deck-optimizer/internal/ids"

// CardID aliases ids.CardID so callers don't need two imports to hold IDs.
type CardID = ids.CardID

// cardsByName maps DisplayName → CardID. Keyed on DisplayName (not bare Name) so each pitch
// variant gets a distinct entry — the three printings of one card share Name but have
// distinct DisplayNames.
var cardsByName = func() map[string]CardID {
	m := make(map[string]CardID, len(cardsByID)-1)
	for id, c := range cardsByID {
		if c == nil {
			continue
		}
		m[c.DisplayName()] = CardID(id)
	}
	return m
}()

// GetCard returns the card for the given ID, or nil when the ID has no registered card.
// Panics on the Invalid sentinel or an out-of-range ID. Callers iterating IDs should
// null-check the result.
func GetCard(id CardID) Card {
	if id == ids.InvalidCard || int(id) >= len(cardsByID) {
		panic("registry: invalid card ID")
	}
	return cardsByID[id]
}

// CardByName looks up an ID by the card's DisplayName ("Aether Slash [R]"). Returns
// (Invalid, false) if no such card.
func CardByName(name string) (CardID, bool) {
	id, ok := cardsByName[name]
	return id, ok
}

// AllCards returns every valid card ID in registration order. Freshly allocated; safe to
// mutate (e.g. shuffle for random deck generation).
func AllCards() []CardID {
	out := make([]CardID, 0, len(cardsByID)-1)
	for id := 1; id < len(cardsByID); id++ {
		if cardsByID[id] == nil {
			continue
		}
		out = append(out, CardID(id))
	}
	return out
}
