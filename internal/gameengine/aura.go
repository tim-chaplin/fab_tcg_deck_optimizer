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

// CreateAura registers a card-sourced aura.
func (gs *GameState) CreateAura(source card.Card, tt triggertype.Type, handler func(card.GameEngine, card.Logger, card.Aura), count int, oncePerTurn bool, filter func(card.TypeSet) bool) {
	gs.auras = append(gs.auras, aura.NewFromCard(source, tt, handler, count, oncePerTurn, filter))
	gs.auraCreated = true
}

// CreateItem registers a card-sourced triggered item.
func (gs *GameState) CreateItem(source card.Card, tt triggertype.Type, handler func(card.GameEngine, card.Logger, card.Item), oncePerTurn bool, filter func(card.TypeSet) bool) {
	gs.items = append(gs.items, item.NewFromCard(source, tt, handler, oncePerTurn, filter))
}

// CreateTrigger registers a card-sourced one-shot ephemeral trigger.
func (gs *GameState) CreateTrigger(source card.Card, tt triggertype.Type, handler func(card.GameEngine, card.Logger, card.EphemeralTrigger), filter func(card.TypeSet) bool) {
	gs.triggers = append(gs.triggers, trigger.NewEphemeralTrigger(source, tt, handler, filter))
}
