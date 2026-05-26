package gameengine

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/aura"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/item"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/trigger"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/triggertype"
)

// Card-facing aura / item / trigger registration. Cards call these via the
// internal/card.GameEngine interface; each method constructs the concrete entry from the
// supplied primitives and appends it onto the engine's hook lists. Tests that need to
// install a pre-built entry (token auras, carry-over from a prior turn, sim's pooled DR
// probe aura) reach for StateBuilder.AddAura or set the slice directly.

// CreateAura registers a card-sourced aura. Takes the next free slot from auraPool
// (the pre-allocated [auraPoolSize]aura.Aura backing on this GameState), overwrites
// it in place, and appends the pool-slot pointer to gs.auras. Panics if the pool is
// exhausted (>auraPoolSize simultaneously live card auras).
func (gs *GameState) CreateAura(source card.Card, tt triggertype.Type, handler func(card.GameEngine, card.Logger, card.Aura), count int, oncePerTurn bool, filter func(card.TypeSet) bool) {
	slot := gs.nextFreeAuraSlot()
	slot.SetFromCard(source, tt, handler, count, oncePerTurn, filter)
	gs.auras = append(gs.auras, slot)
	gs.auraCreated = true
}

// nextFreeAuraSlot returns the first auraPool entry with SourceCard() == nil — the
// canonical "dead slot" test. Panics on pool exhaustion: real decks run 2-5 live card
// auras at a peak; hitting auraPoolSize is a regression worth failing on rather than
// silently slow-pathing onto the heap.
func (gs *GameState) nextFreeAuraSlot() *aura.Aura {
	for i := range gs.auraPool {
		if gs.auraPool[i].SourceCard() == nil {
			return &gs.auraPool[i]
		}
	}
	panic("CreateAura: gameengine auraPool exhausted — increase auraPoolSize or investigate why this GameState carries so many simultaneous card auras")
}

// CreateItem registers a card-sourced triggered item.
func (gs *GameState) CreateItem(source card.Card, tt triggertype.Type, handler func(card.GameEngine, card.Logger, card.Item), oncePerTurn bool, filter func(card.TypeSet) bool) {
	gs.items = append(gs.items, item.NewFromCard(source, tt, handler, oncePerTurn, filter))
}

// CreateTrigger registers a card-sourced one-shot ephemeral trigger.
func (gs *GameState) CreateTrigger(source card.Card, tt triggertype.Type, handler func(card.GameEngine, card.Logger, card.EphemeralTrigger), filter func(card.TypeSet) bool) {
	gs.triggers = append(gs.triggers, trigger.NewEphemeralTrigger(source, tt, handler, filter))
}
