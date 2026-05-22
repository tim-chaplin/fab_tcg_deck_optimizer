package trigger

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/triggertype"
)

// EphemeralTrigger is the one-shot trigger: the engine fires it on the next matching event
// and drops it from the queue. It carries nothing beyond the shared core.
type EphemeralTrigger struct {
	Trigger[card.EphemeralTrigger]
}

// NewEphemeralTrigger builds a one-shot card-sourced ephemeral trigger. typeFilter narrows
// the firing site (currently used only by triggertype.Hit); pass nil for no filter.
func NewEphemeralTrigger(source card.Card, tt triggertype.Type, fire func(card.GameEngine, card.Logger, card.EphemeralTrigger), typeFilter TypeFilter) *EphemeralTrigger {
	return &EphemeralTrigger{Trigger: FromCard[card.EphemeralTrigger](source, tt, fire, false, typeFilter)}
}

// Fire invokes the stored handler with this trigger as the typed receiver.
func (e *EphemeralTrigger) Fire(engine card.GameEngine, logger card.Logger) {
	e.Invoke(engine, logger, e)
}

// Compile-time check that *EphemeralTrigger satisfies card.EphemeralTrigger.
var _ card.EphemeralTrigger = (*EphemeralTrigger)(nil)
