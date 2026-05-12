package turntests

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/cards"
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
)

// Tests that High Striker queues a TriggerHit Trigger so a later attack hit creates
// Copper tokens. Per-pitch Copper count (6/4/2 for R/Y/B) is covered by turntests since
// it requires the chain runner to actually fire.
func TestHighStriker_QueuesTriggerHit(t *testing.T) {
	g := sim.NewTurnStateFromCards(nil, nil)
	sim.ResolveChainStep(g, g.Logger(), &card.CardState{Card: cards.HighStrikerRed{}})
	if got := triggerHitCount(g); got != 1 {
		t.Fatalf("TriggerHit triggers = %d, want 1 (registered the rider)", got)
	}
	if g.CopperCount() != 0 {
		t.Fatalf("Copper = %d before any attack hits, want 0", g.CopperCount())
	}
}

// triggerHitCount returns the number of queued TriggerHit triggers on g. Reaches
// past GameEngine to the concrete *sim.TurnState — Triggers is sim-owned and the
// engine interface stays sim-free.
func triggerHitCount(g card.GameEngine) int {
	n := 0
	for _, t := range g.(*sim.TurnState).Triggers() {
		if t.TriggerType == sim.TriggerHit {
			n++
		}
	}
	return n
}
