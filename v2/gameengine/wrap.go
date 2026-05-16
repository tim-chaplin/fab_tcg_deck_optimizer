package gameengine

import (
	"github.com/tim-chaplin/fab-deck-optimizer/v2/aura"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/trigger"
)

// Bridges between v2/card's typed handler signatures and the type-erased Handler shapes
// v2/aura and v2/trigger store. The leaf packages don't import v2/card by design (so a
// card mistake can't break aura / trigger compilation); the conversion happens here at
// the engine layer, where v2/card is already in scope. Engine and logger values flowing
// through the wrapped closures are *GameEngine and *turnlogger.TurnLogger at runtime, so
// the type assertions succeed by construction.

// WrapAuraHandler adapts a typed card.AuraHandler into the type-erased aura.Handler the
// leaf package stores.
func WrapAuraHandler(h card.AuraHandler) aura.Handler {
	return func(engine, logger, ctx any) {
		h(engine.(card.GameEngine), logger.(card.Logger), ctx.(card.Aura))
	}
}

// WrapTriggerHandler is the trigger.Handler counterpart of WrapAuraHandler.
func WrapTriggerHandler(h card.TriggerHandler) trigger.Handler {
	return func(engine, logger, ctx any) {
		h(engine.(card.GameEngine), logger.(card.Logger), ctx.(card.Trigger))
	}
}

// WrapTypeFilter adapts a card.TypeSet predicate into trigger.TypeFilter (operates on
// `any`). nil filters propagate unchanged so trigger.Matches's nil-fast-path still hits.
func WrapTypeFilter(f func(card.TypeSet) bool) trigger.TypeFilter {
	if f == nil {
		return nil
	}
	return func(types any) bool { return f(types.(card.TypeSet)) }
}
