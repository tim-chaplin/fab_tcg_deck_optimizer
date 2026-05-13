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
func init() {
	gameengine.BuildCardAura = func(self *card.CardState, tt triggertype.Type, h card.AuraHandler, count int, oncePerTurn bool) gameengine.Aura {
		return aura.NewCard(self, tt, h, count, oncePerTurn)
	}
	gameengine.BuildCardTrigger = func(self *card.CardState, tt triggertype.Type, h card.TriggerHandler, filter func(card.TypeSet) bool) gameengine.Trigger {
		return trigger.NewCard(self, tt, h, filter)
	}
	gameengine.BuildRunechantAura = func(n int) gameengine.Aura { return aura.NewRunechant(n) }
	gameengine.BuildPonderAura = func(n int) gameengine.Aura { return aura.NewPonder(n) }
	gameengine.BuildGoldItem = func(n int) gameengine.Item { return token.NewGold(n) }
	gameengine.BuildSilverItem = func(n int) gameengine.Item { return token.NewSilver(n) }
	gameengine.BuildCopperItem = func(n int) gameengine.Item { return token.NewCopper(n) }
}
