// Package trigger owns the concrete one-shot Trigger type — a deferred handler that fires
// on the next matching event (typically triggertype.Hit or triggertype.EndOfTurn) and is
// then removed from the engine's trigger queue. Satisfies gameengine.Trigger structurally
// without importing gameengine; Fire takes a card.GameEngine, which the engine satisfies.
package trigger

import (
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/triggertype"
)

// Trigger is one-shot — the engine fires it on the next matching event, then drops it from
// the queue. typeFilter narrows the firing site when present (e.g. Mauvrion Skies's hit
// trigger "only on attack-action hits"); nil means any matching event qualifies.
type Trigger struct {
	triggerType triggertype.Type
	handler     card.TriggerHandler
	source      card.Card
	typeFilter  func(card.TypeSet) bool
}

// NewCard builds a one-shot trigger whose source is self.Card. typeFilter narrows the
// firing site (currently used only by triggertype.Hit); pass nil for no filter.
func NewCard(self *card.CardState, tt triggertype.Type, handler card.TriggerHandler, typeFilter func(card.TypeSet) bool) *Trigger {
	return &Trigger{
		triggerType: tt,
		handler:     handler,
		source:      self.Card,
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

func (t *Trigger) Matches(types card.TypeSet) bool {
	if t.typeFilter == nil {
		return true
	}
	return t.typeFilter(types)
}

// Fire invokes the trigger's registered handler with a card.Trigger context.
func (t *Trigger) Fire(ge card.GameEngine, l card.Logger) {
	ctx := &triggerCtx{t: t}
	t.handler(ge, l, ctx)
}

// triggerCtx is the adapter handed to each trigger handler — satisfies card.Trigger.
type triggerCtx struct {
	t *Trigger
}

func (c *triggerCtx) CardName() string { return c.t.CardName() }
