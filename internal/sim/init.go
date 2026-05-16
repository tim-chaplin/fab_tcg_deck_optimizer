package sim

import (
	"github.com/tim-chaplin/fab-deck-optimizer/v2/aura"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/gameengine"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/token"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/trigger"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/triggertype"
)

// init wires the concrete Aura / Trigger / Item / token builders living in v2/aura /
// v2/trigger / v2/token into gameengine's factory slots so the engine's card-facing
// Create*Aura / AddXxxTrigger / Create*Token methods can construct entries without
// importing the concrete types directly. Each concrete type satisfies its matching engine
// interface structurally, so the assignment is a no-op box per call.
//
// The aura / trigger builders also wrap the user-supplied typed handlers (card.AuraHandler,
// card.TriggerHandler) into the leaf packages' type-erased Handler closures — keeping
// v2/aura and v2/trigger free of any v2/card dependency.
func init() {
	gameengine.BuildCardAura = func(self *card.CardState, tt triggertype.Type, h card.AuraHandler, count int, oncePerTurn bool) gameengine.Aura {
		return aura.NewFromCard(self.Card, tt, wrapAuraHandler(h), count, oncePerTurn)
	}
	gameengine.BuildCardTrigger = func(self *card.CardState, tt triggertype.Type, h card.TriggerHandler, filter func(card.TypeSet) bool) gameengine.Trigger {
		return trigger.NewCard(self.Card, tt, wrapTriggerHandler(h), wrapTypeFilter(filter))
	}
	gameengine.BuildRunechantAura = func(n int) gameengine.Aura { return token.NewRunechant(n) }
	gameengine.BuildPonderAura = func(n int) gameengine.Aura { return token.NewPonder(n) }
	gameengine.BuildGoldItem = func(n int) gameengine.Item { return token.NewGold(n) }
	gameengine.BuildSilverItem = func(n int) gameengine.Item { return token.NewSilver(n) }
	gameengine.BuildCopperItem = func(n int) gameengine.Item { return token.NewCopper(n) }
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
	return func(engine, logger any, ctx trigger.Ctx) {
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
