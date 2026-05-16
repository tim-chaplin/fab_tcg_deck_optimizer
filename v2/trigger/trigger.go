// Package trigger owns the concrete one-shot Trigger type — a deferred handler that fires
// on the next matching event (typically triggertype.Hit or triggertype.EndOfTurn) and is
// then removed from the engine's trigger queue.
//
// Handlers are typed against v2/card directly (card.TriggerHandler taking
// card.GameEngine / card.Logger / card.Trigger); the package's own ctx struct implements
// card.Trigger so handler bodies receive the typed surface without an outer wrapping layer.
package trigger

import (
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/triggertype"
)

// TypeFilter narrows the firing site to a card-type predicate. nil means any matching
// event qualifies.
type TypeFilter func(card.TypeSet) bool

// Trigger is one-shot — the engine fires it on the next matching event, then drops it from
// the queue. typeFilter narrows the firing site when present (e.g. Mauvrion Skies's hit
// trigger "only on attack-action hits"); nil means any matching event qualifies.
type Trigger struct {
	triggerType triggertype.Type
	fire        card.TriggerHandler
	source      card.Card
	typeFilter  TypeFilter
}

// NewFromCard builds a one-shot trigger whose source is the supplied card. typeFilter
// narrows the firing site (currently used only by triggertype.Hit); pass nil for no filter.
func NewFromCard(source card.Card, tt triggertype.Type, fire card.TriggerHandler, typeFilter TypeFilter) *Trigger {
	return &Trigger{
		triggerType: tt,
		fire:        fire,
		source:      source,
		typeFilter:  typeFilter,
	}
}

func (t *Trigger) TriggerType() triggertype.Type { return t.triggerType }

func (t *Trigger) CardName() string {
	if t.source != nil {
		return t.source.DisplayName()
	}
	return ""
}

// Matches reports whether the trigger's type filter accepts the firing event's type set.
func (t *Trigger) Matches(types card.TypeSet) bool {
	if t.typeFilter == nil {
		return true
	}
	return t.typeFilter(types)
}

// Fire invokes the stored handler. The handler signature is typed against v2/card so the
// engine and logger arguments pass straight through without a per-fire assertion.
func (t *Trigger) Fire(engine card.GameEngine, logger card.Logger) {
	t.fire(engine, logger, &ctx{t: t})
}

// ctx is the per-fire value Fire passes to the handler. Implements card.Trigger.
type ctx struct {
	t *Trigger
}

// Compile-time check that ctx satisfies card.Trigger.
var _ card.Trigger = (*ctx)(nil)

func (c *ctx) CardName() string { return c.t.CardName() }
