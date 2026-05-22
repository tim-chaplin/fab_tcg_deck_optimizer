// Package item owns the concrete Item type — an in-play permanent. A token item carries an
// activated ability the chain runner enqueues as a 1-AP playable; a card-sourced item
// carries a trigger that fires on a scheduled event. Both embed the shared trigger.Trigger
// core; token items leave it zero-valued so trigger type 0 never matches a firing event.
package item

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/trigger"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/triggertype"
)

// Item is the concrete entry the engine stores in its persistent item list. The embedded
// trigger.Trigger holds the source / token identity and, for card-sourced items, the
// firing event and handler; Item adds the stack count and the optional activated ability.
//
// activeEngine is set by Fire so a handler-side Destroy routes back through
// engine.DestroyItem without allocating a per-fire wrapper. Cleared after the handler
// returns.
type Item struct {
	trigger.Trigger[card.Item]
	activeEngine card.GameEngine
	ability      any
	count        int
}

// NewFromToken builds a token-sourced item with the supplied name, identifier,
// activated-ability card, and initial count. Production callers reach for internal/token's
// NewGold / NewSilver / NewCopper.
func NewFromToken(name string, tokenID ids.CardID, ability any, count int) *Item {
	return &Item{
		Trigger: trigger.FromToken[card.Item](name, tokenID, 0, nil, false, nil),
		ability: ability,
		count:   count,
	}
}

// NewFromCard builds a card-sourced triggered item — an in-play permanent whose handler
// fires when an event of type tt resolves. count starts at 1; oncePerTurn caps it to one
// fire per turn.
func NewFromCard(source card.Card, tt triggertype.Type, fire func(card.GameEngine, card.Logger, card.Item), oncePerTurn bool) *Item {
	return &Item{
		Trigger: trigger.FromCard[card.Item](source, tt, fire, oncePerTurn, nil),
		count:   1,
	}
}

func (i *Item) Count() int     { return i.count }
func (i *Item) SetCount(n int) { i.count = n }

// Ability returns the activated-ability card boxed as any, or nil for a card-sourced item
// with no activated ability. Callers type-assert to their richer card type when invoking.
func (i *Item) Ability() any { return i.ability }

// Fire invokes the stored handler with this item as the typed receiver. activeEngine is set
// for the handler's duration and cleared afterward.
func (i *Item) Fire(engine card.GameEngine, logger card.Logger) {
	i.activeEngine = engine
	i.Invoke(engine, logger, i)
	i.activeEngine = nil
}

// Destroy ends the item currently being fired, routed through activeEngine. Calling it
// outside a Fire window panics deterministically rather than silently no-oping.
func (i *Item) Destroy(addToGraveyard bool) {
	i.activeEngine.DestroyItem(addToGraveyard)
}

// Copy returns a deep copy boxed as any so the gameengine.Item interface can avoid
// referencing the concrete type. activeEngine is cleared on the copy.
func (i *Item) Copy() any {
	out := *i
	out.activeEngine = nil
	return &out
}

// Compile-time check that *Item satisfies card.Item.
var _ card.Item = (*Item)(nil)
