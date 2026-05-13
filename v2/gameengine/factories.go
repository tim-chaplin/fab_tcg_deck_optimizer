package gameengine

import (
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/triggertype"
)

// Concrete Aura / Trigger / Item types live outside gameengine (in v2/aura, v2/trigger,
// v2/token). The engine's card-facing Create*Aura / AddXxxTrigger / Create*Token methods
// construct entries by calling these builder funcs; sim registers them via init.
//
// The factory layer is also where typed card.AuraHandler / card.TriggerHandler values are
// wrapped into the leaf packages' type-erased Handler closures — keeping aura / trigger
// free of any v2/card dependency.
//
// Engines built before the builders are registered crash on the cards-facing construction
// methods — by design: a runtime with no concrete types behind the engine can't build
// anything, so silently no-oping would hide a setup bug.
var (
	// BuildCardAura constructs a card-backed aura whose source is self.Card.
	BuildCardAura func(self *card.CardState, tt triggertype.Type, handler card.AuraHandler, count int, oncePerTurn bool) Aura

	// BuildCardTrigger constructs a one-shot trigger whose source is self.Card.
	BuildCardTrigger func(self *card.CardState, tt triggertype.Type, handler card.TriggerHandler, typeFilter func(card.TypeSet) bool) Trigger

	// BuildRunechantAura returns a runechant token aura at count n.
	BuildRunechantAura func(n int) Aura

	// BuildPonderAura returns a ponder token aura at count n.
	BuildPonderAura func(n int) Aura

	// BuildGoldItem / BuildSilverItem / BuildCopperItem return fresh token items at count n.
	BuildGoldItem   func(n int) Item
	BuildSilverItem func(n int) Item
	BuildCopperItem func(n int) Item
)
