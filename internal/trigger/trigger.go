// Package trigger owns the concrete one-shot Trigger type — a deferred handler that fires
// on the next matching event (typically triggertype.Hit or triggertype.EndOfTurn) and is
// then removed from the engine's trigger queue.
//
// *Trigger itself satisfies card.Trigger, so Fire passes the trigger directly to the
// handler with zero allocation.
package trigger

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/triggertype"
)

// Handler is the typed trigger handler signature. Lives in internal/trigger (not internal/card)
// because the shape is meaningful only to trigger: it's the signature stored on each
// Trigger and called at every Fire. card.GameEngine.AddHitTrigger et al. inline the
// function type in their parameter declarations so internal/card doesn't need to import
// internal/trigger.
type Handler func(card.GameEngine, card.Logger, card.Trigger)

// TypeFilter narrows the firing site to a card-type predicate. nil means any matching
// event qualifies.
type TypeFilter func(card.TypeSet) bool

// Trigger is one-shot — the engine fires it on the next matching event, then drops it from
// the queue. typeFilter narrows the firing site when present (e.g. Mauvrion Skies's hit
// trigger "only on attack-action hits"); nil means any matching event qualifies.
type Trigger struct {
	triggerType triggertype.Type
	fire        Handler
	source      card.Card
	typeFilter  TypeFilter
}

// NewFromCard builds a one-shot trigger whose source is the supplied card. typeFilter
// narrows the firing site (currently used only by triggertype.Hit); pass nil for no filter.
func NewFromCard(source card.Card, tt triggertype.Type, fire Handler, typeFilter TypeFilter) *Trigger {
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

// Fire invokes the stored handler. *Trigger itself satisfies card.Trigger via its
// CardName method, so the handler receives the typed surface directly with no per-fire
// wrapper allocation.
func (t *Trigger) Fire(engine card.GameEngine, logger card.Logger) {
	t.fire(engine, logger, t)
}

// Compile-time check that *Trigger satisfies card.Trigger.
var _ card.Trigger = (*Trigger)(nil)
