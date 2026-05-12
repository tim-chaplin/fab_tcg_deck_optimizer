package sim

import (
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/gameengine"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/triggertype"
)

// init wires sim's concrete Aura / Trigger / Item / token builders into gameengine's
// factory slots so the engine's card-facing Create*Aura / AddXxxTrigger / Create*Token
// methods can construct entries without importing sim.
//
// The wrappers adapt sim's concrete return types (*Aura, *Trigger, *Item) to the engine's
// interface-return factory signatures — every concrete type satisfies the matching engine
// interface, so the assignment is a no-op box per call.
func init() {
	gameengine.BuildCardAura = func(self *card.CardState, tt triggertype.Type, h card.AuraHandler, count int, oncePerTurn bool) gameengine.Aura {
		return NewCardAura(self, tt, h, count, oncePerTurn)
	}
	gameengine.BuildCardTrigger = func(self *card.CardState, tt triggertype.Type, h card.TriggerHandler, filter func(card.TypeSet) bool) gameengine.Trigger {
		return NewCardTrigger(self, tt, h, filter)
	}
	gameengine.BuildRunechantAura = func(n int) gameengine.Aura { return NewRunechantAura(n) }
	gameengine.BuildPonderAura = func(n int) gameengine.Aura { return NewPonderAura(n) }
	gameengine.BuildGoldItem = func(n int) gameengine.Item { return NewGoldItem(n) }
	gameengine.BuildSilverItem = func(n int) gameengine.Item { return NewSilverItem(n) }
	gameengine.BuildCopperItem = func(n int) gameengine.Item { return NewCopperItem(n) }
}
