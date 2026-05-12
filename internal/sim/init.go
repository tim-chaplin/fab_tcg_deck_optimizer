package sim

import "github.com/tim-chaplin/fab-deck-optimizer/v2/gameengine"

// init wires sim's concrete Aura / Trigger / Item / token builders into gameengine's
// factory slots so the engine's card-facing Create*Aura / AddXxxTrigger / Create*Token
// methods can construct entries without importing sim.
//
// Importing sim transitively registers everything — gameengine standalone (or with a
// different builder set) is supported but every engine that goes through sim's chain
// runner picks up the canonical FaB token impls here.
func init() {
	gameengine.BuildCardAura = NewCardAura
	gameengine.BuildCardTrigger = NewCardTrigger
	gameengine.BuildRunechantAura = NewRunechantAura
	gameengine.BuildPonderAura = NewPonderAura
	gameengine.BuildGoldItem = NewGoldItem
	gameengine.BuildSilverItem = NewSilverItem
	gameengine.BuildCopperItem = NewCopperItem
}
