package cards

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

// Tests that High Striker queues a TriggerHit Trigger so a later attack hit creates
// Copper tokens. Per-pitch Copper count (6/4/2 for R/Y/B) is covered by turntests since
// it requires the chain runner to actually fire.
func TestHighStriker_QueuesTriggerHit(t *testing.T) {
	s := sim.NewTurnStateFromCards(nil, nil)
	sim.ResolveChainStep(s, s.Logger(), &sim.CardState{Card: HighStrikerRed{}})
	if got := triggerHitCount(s); got != 1 {
		t.Fatalf("TriggerHit triggers = %d, want 1 (registered the rider)", got)
	}
	if s.Copper() != 0 {
		t.Fatalf("Copper = %d before any attack hits, want 0", s.Copper())
	}
}

// triggerHitCount returns the number of queued TriggerHit triggers on s. Reaches
// past GameEngine to the concrete *sim.TurnState — Triggers is sim-owned and the
// engine interface stays sim-free.
func triggerHitCount(s sim.GameEngine) int {
	n := 0
	for _, t := range s.(*sim.TurnState).Triggers() {
		if t.TriggerType == sim.TriggerHit {
			n++
		}
	}
	return n
}
