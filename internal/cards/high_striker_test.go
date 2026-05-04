package cards

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

// Tests that High Striker queues a NextAttackActionHit trigger so a later attack
// action's hit creates Copper tokens. Per-pitch Copper count (6/4/2 for R/Y/B) is
// covered by the e2e tests since it requires the chain runner to actually fire.
func TestHighStriker_QueuesAttackActionHitTrigger(t *testing.T) {
	s := sim.NewTurnState(nil, nil)
	(HighStrikerRed{}).Play(s, &sim.CardState{Card: HighStrikerRed{}})
	if got := s.PendingNextAttackActionHits(); got != 1 {
		t.Fatalf("PendingNextAttackActionHits = %d, want 1 (registered the rider)", got)
	}
	if s.Copper() != 0 {
		t.Fatalf("Copper = %d before any attack hits, want 0", s.Copper())
	}
}
