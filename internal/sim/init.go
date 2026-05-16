package sim

import (
	"github.com/tim-chaplin/fab-deck-optimizer/v2/aura"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/gameengine"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/trigger"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/triggertype"
)

// init wires the card-handler-shaped builders (BuildCardAura / BuildCardTrigger) into
// gameengine's factory slots. These two stay in sim because they wrap user-supplied typed
// handlers (card.AuraHandler, card.TriggerHandler) into the leaf packages' type-erased
// Handler closures — that wrapping requires v2/card, which v2/aura and v2/trigger
// deliberately don't import. The simpler item / aura-flavored token slots
// (BuildGoldItem etc.) wire themselves from v2/token/init.go.
func init() {
	gameengine.BuildCardAura = func(self *card.CardState, tt triggertype.Type, h card.AuraHandler, count int, oncePerTurn bool) gameengine.Aura {
		return aura.NewFromCard(self.Card, tt, wrapAuraHandler(h), count, oncePerTurn)
	}
	gameengine.BuildCardTrigger = func(self *card.CardState, tt triggertype.Type, h card.TriggerHandler, filter func(card.TypeSet) bool) gameengine.Trigger {
		return trigger.NewFromCard(self.Card, tt, wrapTriggerHandler(h), wrapTypeFilter(filter))
	}
}

// wrapAuraHandler adapts a typed card.AuraHandler into the type-erased aura.Handler the
// leaf package stores. The engine/logger values flowing through aura.Fire are *GameEngine
// / *TurnLogger at runtime; the type assertions succeed by construction.
func wrapAuraHandler(h card.AuraHandler) aura.Handler {
	return func(engine, logger, ctx any) {
		h(engine.(card.GameEngine), logger.(card.Logger), ctx.(card.Aura))
	}
}

// wrapTriggerHandler is the trigger.Handler counterpart of wrapAuraHandler.
func wrapTriggerHandler(h card.TriggerHandler) trigger.Handler {
	return func(engine, logger, ctx any) {
		h(engine.(card.GameEngine), logger.(card.Logger), ctx.(card.Trigger))
	}
}

// wrapTypeFilter adapts a card.TypeSet predicate into trigger.TypeFilter (operates on
// `any`). nil filters propagate unchanged so trigger.Matches's nil-fast-path still hits.
func wrapTypeFilter(f func(card.TypeSet) bool) trigger.TypeFilter {
	if f == nil {
		return nil
	}
	return func(types any) bool { return f(types.(card.TypeSet)) }
}
