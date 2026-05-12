package sim

import (
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/gameengine"
)

// Trigger is the sim concrete impl of gameengine.Trigger. Card-driven AddXxxTrigger methods
// on *gameengine.GameEngine build one of these via the registered builder.
type Trigger struct {
	triggerType gameengine.TriggerType
	handler     card.TriggerHandler
	source      card.Card // nil when no source card (e.g. start-of-turn-style triggers)
	typeFilter  func(card.TypeSet) bool
}

// NewCardTrigger builds a one-shot trigger whose source is self.Card. typeFilter narrows
// the firing site (currently used only by TriggerHit); pass nil for no filter.
func NewCardTrigger(self *card.CardState, tt gameengine.TriggerType, handler card.TriggerHandler, typeFilter func(card.TypeSet) bool) gameengine.Trigger {
	return &Trigger{
		triggerType: tt,
		handler:     handler,
		source:      self.Card,
		typeFilter:  typeFilter,
	}
}

func (t *Trigger) TriggerType() gameengine.TriggerType { return t.triggerType }
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

func (t *Trigger) Fire(g *gameengine.GameEngine, l card.Logger) {
	ctx := &triggerCtx{t: t}
	t.handler(g, l, ctx)
}

func (t *Trigger) Clone() gameengine.Trigger {
	out := *t
	return &out
}

// triggerCtx is the adapter handed to each trigger handler — satisfies card.Trigger.
type triggerCtx struct {
	t *Trigger
}

func (c *triggerCtx) CardName() string { return c.t.CardName() }
