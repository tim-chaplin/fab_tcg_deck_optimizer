package gameengine

import "github.com/tim-chaplin/fab-deck-optimizer/v2/card"

// Concrete Aura / Trigger / Item types live outside gameengine (sim today, eventually
// their own v2/ packages). The engine's card-facing Create*Aura / AddXxxTrigger /
// Create*Token methods construct entries by calling these builder funcs; the consuming
// package registers them via init.
//
// Engines built before the builders are registered crash on the cards-facing
// construction methods — by design: a runtime with no concrete types behind the engine
// can't build anything, so silently no-oping would hide a setup bug.
var (
	// BuildCardAura constructs a card-backed aura whose source is self.Card. Registered
	// by sim with sim.NewCardAura.
	BuildCardAura func(self *card.CardState, tt TriggerType, handler card.AuraHandler, count int, oncePerTurn bool) Aura

	// BuildCardTrigger constructs a one-shot trigger whose source is self.Card. Registered
	// by sim with sim.NewCardTrigger.
	BuildCardTrigger func(self *card.CardState, tt TriggerType, handler card.TriggerHandler, typeFilter func(card.TypeSet) bool) Trigger

	// BuildRunechantAura returns a runechant token aura at count n. Registered by sim
	// with sim.NewRunechantAura.
	BuildRunechantAura func(n int) Aura

	// BuildPonderAura returns a ponder token aura at count n. Registered by sim with
	// sim.NewPonderAura.
	BuildPonderAura func(n int) Aura

	// BuildGoldItem / BuildSilverItem / BuildCopperItem return fresh token items at
	// count n. Registered by sim with sim.NewGoldItem / NewSilverItem / NewCopperItem.
	BuildGoldItem   func(n int) Item
	BuildSilverItem func(n int) Item
	BuildCopperItem func(n int) Item
)
