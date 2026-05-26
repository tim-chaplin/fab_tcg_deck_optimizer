package turntests

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card/cards"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/gameengine"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/triggertype"
)

// Tests that High Striker queues a TriggerHit Trigger so a later attack hit creates
// Copper tokens. Per-pitch Copper count (6/4/2 for R/Y/B) is covered by turntests since
// it requires the attack-turn runner to actually fire.
func TestHighStriker_QueuesTriggerHit(t *testing.T) {
	ge := gameengine.New()
	ge.ResolveAttackStep(ge.Logger(), &card.CardState{Card: cards.HighStrikerRed{}})
	if got := triggerHitCount(ge); got != 1 {
		t.Fatalf("TriggerHit triggers = %d, want 1 (registered the rider)", got)
	}
	if ge.CopperCount() != 0 {
		t.Fatalf("Copper = %d before any attack hits, want 0", ge.CopperCount())
	}
}

// triggerHitCount returns the number of queued TriggerHit triggers on g. Reaches
// past GameEngine to the concrete *gameengine.GameEngine — Triggers is sim-owned and the
// engine interface stays sim-free.
func triggerHitCount(ge card.GameEngine) int {
	n := 0
	for _, t := range ge.(*gameengine.GameEngine).Triggers() {
		if t.TriggerType() == triggertype.Hit {
			n++
		}
	}
	return n
}
