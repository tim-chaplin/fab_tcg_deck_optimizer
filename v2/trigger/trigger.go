// Package trigger owns the concrete one-shot Trigger type — a deferred handler that fires
// on the next matching event (typically triggertype.Hit or triggertype.EndOfTurn) and is
// then removed from the engine's trigger queue.
//
// The package defines its own narrow interfaces (Card, Ctx, Handler) and does not import
// v2/card or v2/gameengine. Stored handlers are typed as Handler — engine/logger params
// are `any` so trigger never depends on the consumer's broader engine surface. Callers
// wrap their typed handlers into Handler closures at the factory layer (sim/init.go).
package trigger

import (
	"github.com/tim-chaplin/fab-deck-optimizer/v2/triggertype"
)

// Card is the minimal source-card surface trigger needs (display name). The richer
// v2/card.Card satisfies it structurally so consumers can pass card values through.
type Card interface {
	DisplayName() string
}

// Handler is the per-trigger stored handler. Engine and logger params are typed as `any`
// so trigger doesn't depend on the consumer's interfaces — the factory layer wraps typed
// user handlers into closures of this shape.
type Handler func(engine, logger any, ctx Ctx)

// Ctx is the per-fire surface trigger handlers see on the firing trigger.
type Ctx interface {
	CardName() string
}

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

// NewCard builds a one-shot trigger whose source is the supplied card. typeFilter narrows
// the firing site (currently used only by triggertype.Hit); pass nil for no filter.
func NewCard(source Card, tt triggertype.Type, fire Handler, typeFilter TypeFilter) *Trigger {
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

// Fire invokes the stored handler. engine/logger are passed through as `any` — the handler
// (typically wrapped at the factory layer) asserts them to its preferred typed interfaces.
func (t *Trigger) Fire(engine, logger any) {
	t.fire(engine, logger, &ctx{t: t})
}

// ctx is the adapter handed to each trigger handler — satisfies the Ctx interface.
type ctx struct {
	t *Trigger
}

func (c *ctx) CardName() string { return c.t.CardName() }
