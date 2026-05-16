package token

import (
	"github.com/tim-chaplin/fab-deck-optimizer/v2/gameengine"
)

// init wires the per-token-flavor factories into gameengine's Build* slots so the
// engine's card-facing CreateRunechants / CreatePonders / CreateGold / CreateSilver /
// CreateCopper methods can construct entries without importing the concrete types
// directly. Each concrete return value satisfies its matching engine interface
// structurally, so the assignment is a no-op box per call.
//
// Importing v2/token (production or tests) is what enables the engine to actually
// produce tokens. The card-handler wrapping for triggers and auras stays in
// internal/sim/init.go because those require v2/card type assertions inside the wrap.
func init() {
	gameengine.BuildRunechantAura = func(n int) gameengine.Aura { return NewRunechant(n) }
	gameengine.BuildPonderAura = func(n int) gameengine.Aura { return NewPonder(n) }
	gameengine.BuildGoldItem = func(n int) gameengine.Item { return NewGold(n) }
	gameengine.BuildSilverItem = func(n int) gameengine.Item { return NewSilver(n) }
	gameengine.BuildCopperItem = func(n int) gameengine.Item { return NewCopper(n) }
}
