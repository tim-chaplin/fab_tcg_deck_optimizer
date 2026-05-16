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
// Handler types (aura.Handler, trigger.Handler) are leaf-package typed and pass straight
// through — no wrap layer, since v2/aura and v2/trigger import v2/card and their stored
// handler signatures match what cards write.
func init() {
	gameengine.BuildCardAura = func(self *card.CardState, tt triggertype.Type, h aura.Handler, count int, oncePerTurn bool) gameengine.Aura {
		return aura.NewFromCard(self.Card, tt, h, count, oncePerTurn)
	}
	gameengine.BuildCardTrigger = func(self *card.CardState, tt triggertype.Type, h trigger.Handler, filter trigger.TypeFilter) gameengine.Trigger {
		return trigger.NewFromCard(self.Card, tt, h, filter)
	}
	gameengine.BuildRunechantAura = func(n int) gameengine.Aura { return token.NewRunechant(n) }
	gameengine.BuildPonderAura = func(n int) gameengine.Aura { return token.NewPonder(n) }
	gameengine.BuildGoldItem = func(n int) gameengine.Item { return token.NewGold(n) }
	gameengine.BuildSilverItem = func(n int) gameengine.Item { return token.NewSilver(n) }
	gameengine.BuildCopperItem = func(n int) gameengine.Item { return token.NewCopper(n) }
}
