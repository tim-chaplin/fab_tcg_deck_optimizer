// Package simreg wires v2/registry into sim's forward-declared hooks. See
// docs/dev-standards.md "Registry / sim split" for the import contract and when to blank-
// import this package.
package simreg

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/optimizations"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/registry"
)

// Pre-asserts each registered card / weapon to its sim.Card / sim.Weapon view at init time
// so the runtime hooks are pure slice / direct reads with no per-call type assertion.
func init() {
	allIDs := registry.AllCards()
	maxID := 0
	for _, id := range allIDs {
		if int(id) > maxID {
			maxID = int(id)
		}
	}
	cardsByID := make([]sim.Card, maxID+1)
	for _, id := range allIDs {
		cardsByID[id] = registry.GetCard(id).(sim.Card)
	}
	sim.GetCard = func(id ids.CardID) sim.Card {
		if id == ids.InvalidCard || int(id) >= len(cardsByID) {
			panic("sim.GetCard: invalid card ID")
		}
		return cardsByID[id]
	}
	sim.DeckableCards = registry.DeckableCards

	sim.AllWeapons = make([]sim.Weapon, len(registry.AllWeapons))
	for i, w := range registry.AllWeapons {
		sim.AllWeapons[i] = w.(sim.Weapon)
	}

	optimizations.WarmChainStepCache(cardsByID)
}
