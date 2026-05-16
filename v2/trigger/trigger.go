// Package trigger owns the concrete one-shot Trigger type — a deferred handler that fires
// on the next matching event (typically triggertype.Hit or triggertype.EndOfTurn) and is
// then removed from the engine's trigger queue.
//
// The package defines its own narrow interfaces (Card in interfaces.go) and does not
// import v2/card or v2/gameengine. Stored handlers are typed as Handler — engine, logger,
// and the per-fire ctx all flow through as `any` so trigger never depends on the
// consumer's broader surfaces. Each consumer defines its own interface (e.g. v2/card.Trigger
// for the per-fire ctx) and asserts at the handler boundary; the factory layer
// (sim/init.go) wraps typed user handlers into Handler closures.
package trigger

import (
	"github.com/tim-chaplin/fab-deck-optimizer/v2/triggertype"
)

// Handler is the per-trigger stored handler. Every argument is typed as `any` — engine,
// logger, and the per-fire ctx — so trigger doesn't depend on the consumer's interfaces.
// The factory layer wraps typed user handlers into closures of this shape; handler bodies
// type-assert to whatever richer interfaces they need.
type Handler func(engine, logger, ctx any)

// TypeFilter narrows the firing site to a card-type predicate. The filter receives the
// triggering card's type set as `any`; the wrapping closure asserts back to the concrete
// TypeSet type. nil means any matching event qualifies.
type TypeFilter func(types any) bool

// Trigger is one-shot — the engine fires it on the next matching event, then drops it from
// the queue. typeFilter narrows the firing site when present (e.g. Mauvrion Skies's hit
// trigger "only on attack-action hits"); nil means any matching event qualifies.
type Trigger struct {
	triggerType triggertype.Type
	fire        Handler
	source      Card
	typeFilter  TypeFilter
}

// NewFromCard builds a one-shot trigger whose source is the supplied card. typeFilter
// narrows the firing site (currently used only by triggertype.Hit); pass nil for no filter.
func NewFromCard(source Card, tt triggertype.Type, fire Handler, typeFilter TypeFilter) *Trigger {
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
// The engine passes its card.TypeSet through as `any`; the wrapped filter asserts back.
func (t *Trigger) Matches(types any) bool {
	if t.typeFilter == nil {
		return true
	}
	return t.typeFilter(types)
}

// Fire invokes the stored handler. engine, logger, and the per-fire ctx pass through as
// `any` — the handler (typically wrapped at the factory layer) asserts them to its
// preferred typed interfaces.
func (t *Trigger) Fire(engine, logger any) {
	t.fire(engine, logger, &ctx{t: t})
}

// ctx is the concrete per-fire value Fire passes to the handler as `any`. Consumers
// define their own interface (e.g. v2/card.Trigger) and assert this value to it; the
// method set here is what every such interface can rely on.
type ctx struct {
	t *Trigger
}

func (c *ctx) CardName() string { return c.t.CardName() }
