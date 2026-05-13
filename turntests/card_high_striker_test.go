package turntests

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/cards"

	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/gameengine"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/triggertype"
)

// Tests that High Striker queues a TriggerHit Trigger so a later attack hit creates
// Copper tokens. Per-pitch Copper count (6/4/2 for R/Y/B) is covered by turntests since
// it requires the chain runner to actually fire.
func TestHighStriker_QueuesTriggerHit(t *testing.T) {
	g := &gameengine.GameEngine{GameState: gameengine.GameStateBuilder().Build()}
	g.ResolveChainStep(g.Logger(), &card.CardState{Card: cards.HighStrikerRed{}})
	if got := triggerHitCount(g); got != 1 {
		t.Fatalf("TriggerHit triggers = %d, want 1 (registered the rider)", got)
	}
	if g.CopperCount() != 0 {
		t.Fatalf("Copper = %d before any attack hits, want 0", g.CopperCount())
	}
}

// triggerHitCount returns the number of queued TriggerHit triggers on g. Reaches
// past GameEngine to the concrete *gameengine.GameEngine — Triggers is sim-owned and the
// engine interface stays sim-free.
func triggerHitCount(g card.GameEngine) int {
	n := 0
	for _, t := range g.(*gameengine.GameEngine).Triggers() {
		if t.TriggerType() == triggertype.Hit {
			n++
		}
	}
	return n
}
