package gameengine

import (
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
)

// triggerEntry is the canonical engine-owned Trigger impl. Card-driven AddXxxTrigger
// methods on *GameEngine build one of these and append it to g.triggers.
type triggerEntry struct {
	triggerType TriggerType
	handler     card.TriggerHandler
	source      card.Card // nil when no source card (start-of-turn-style triggers)
	typeFilter  func(card.TypeSet) bool
}

// NewCardTrigger builds a one-shot trigger whose source is self.Card. typeFilter narrows
// the firing site (currently used only by TriggerHit); pass nil for no filter.
func NewCardTrigger(self *card.CardState, tt TriggerType, handler card.TriggerHandler, typeFilter func(card.TypeSet) bool) Trigger {
	return &triggerEntry{
		triggerType: tt,
		handler:     handler,
		source:      self.Card,
		typeFilter:  typeFilter,
	}
}

func (t *triggerEntry) TriggerType() TriggerType { return t.triggerType }
func (t *triggerEntry) CardName() string {
	if t.source != nil {
		return t.source.DisplayName()
	}
	return ""
}
func (t *triggerEntry) Matches(types card.TypeSet) bool {
	if t.typeFilter == nil {
		return true
	}
	return t.typeFilter(types)
}

func (t *triggerEntry) Fire(g *GameEngine, l card.Logger) {
	ctx := &triggerCtx{t: t}
	t.handler(g, l, ctx)
}

func (t *triggerEntry) Clone() Trigger {
	out := *t
	return &out
}

// triggerCtx is the engine-built adapter handed to each trigger handler — satisfies the
// card.Trigger interface so cards interact through a stable surface.
type triggerCtx struct {
	t *triggerEntry
}

func (c *triggerCtx) CardName() string { return c.t.CardName() }

// === Card-facing trigger registration on GameEngine ===

// AddHitTrigger registers a one-shot TriggerHit listener. filter narrows the qualifying
// hits to a card-type predicate; nil = any hit qualifies.
func (g *GameEngine) AddHitTrigger(self *card.CardState, handler card.TriggerHandler, filter func(card.TypeSet) bool) {
	g.AppendTrigger(NewCardTrigger(self, TriggerHit, handler, filter))
}

// AddEndOfTurnTrigger registers a one-shot TriggerEndOfTurn listener — fires after the
// chain finishes resolving but before the carry-state snapshot.
func (g *GameEngine) AddEndOfTurnTrigger(self *card.CardState, handler card.TriggerHandler) {
	g.AppendTrigger(NewCardTrigger(self, TriggerEndOfTurn, handler, nil))
}
