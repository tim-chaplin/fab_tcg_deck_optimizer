package sim

import "github.com/tim-chaplin/fab-deck-optimizer/v2/card"

// Trigger is a one-shot deferred handler keyed to a TriggerType. The sim removes
// each fired entry from TurnState.Triggers after invoking Handler. Long-lived
// listeners that survive multiple fires (Sigil-of-X start-of-turn riders, etc.)
// use Aura instead — the two structs share fields by intent but no longer share
// a Go type.
//
// Handlers MUST NOT call AddXxxTrigger from within a fire walk: the firing helper
// snapshots len(s.triggers) before iterating, so a trigger added during fire stays
// queued for the next matching event but isn't seen on the current pass.
type Trigger struct {
	// Source is the originating card, surfaced in the trace via Source.DisplayName().
	Source card.Card
	// TriggerType matches the firing site (TriggerEndOfTurn, TriggerHit, …).
	TriggerType TriggerType
	// TypeFilter optionally narrows TriggerHit to a card-type predicate. nil = no filter.
	TypeFilter func(card.TypeSet) bool
	// Handler runs when TriggerType (and TypeFilter, if set) match.
	Handler card.TriggerHandler
}

// triggerCtx adapts *Trigger to card.Trigger for the firing-pipeline invocation.
type triggerCtx struct {
	t *Trigger
}

func (c *triggerCtx) SourceName() string {
	if c.t.Source != nil {
		return c.t.Source.DisplayName()
	}
	return ""
}
