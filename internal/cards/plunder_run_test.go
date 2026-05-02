package cards

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
)

// Plunder Run from hand registers the on-hit-draw trigger and skips the +N{p} grant.
func TestPlunderRun_FromHandQueuesTriggerNoBonus(t *testing.T) {
	target := &sim.CardState{Card: testutils.GenericAttack(0, 4)}
	s := sim.TurnState{CardsRemaining: []*sim.CardState{target}}
	self := &sim.CardState{Card: PlunderRunRed{}}
	(PlunderRunRed{}).Play(&s, self)
	if got := s.PendingNextAttackActionHits(); got != 1 {
		t.Errorf("queued triggers = %d, want 1", got)
	}
	if target.BonusAttack != 0 {
		t.Errorf("target.BonusAttack = %d, want 0 (hand-played skips +N{p})", target.BonusAttack)
	}
}

// From arsenal: registers the trigger and adds +N{p} to the next attack action in
// CardsRemaining. Each printing carries its own N.
func TestPlunderRun_FromArsenalAddsBonusAttack(t *testing.T) {
	cases := []struct {
		c        sim.Card
		wantBoon int
	}{
		{PlunderRunRed{}, 3},
		{PlunderRunYellow{}, 2},
		{PlunderRunBlue{}, 1},
	}
	for _, tc := range cases {
		target := &sim.CardState{Card: testutils.GenericAttack(0, 4)}
		s := sim.TurnState{CardsRemaining: []*sim.CardState{target}}
		self := &sim.CardState{Card: tc.c, FromArsenal: true}
		tc.c.Play(&s, self)
		if got := s.PendingNextAttackActionHits(); got != 1 {
			t.Errorf("%s: queued triggers = %d, want 1", tc.c.Name(), got)
		}
		if target.BonusAttack != tc.wantBoon {
			t.Errorf("%s: target.BonusAttack = %d, want %d", tc.c.Name(), target.BonusAttack, tc.wantBoon)
		}
	}
}

// Multiple Plunder Runs queue independent triggers — they all fire on the same hit.
func TestPlunderRun_TriggersStack(t *testing.T) {
	s := sim.TurnState{}
	(PlunderRunRed{}).Play(&s, &sim.CardState{Card: PlunderRunRed{}})
	(PlunderRunBlue{}).Play(&s, &sim.CardState{Card: PlunderRunBlue{}})
	if got := s.PendingNextAttackActionHits(); got != 2 {
		t.Errorf("queued triggers = %d, want 2 (two independent listeners)", got)
	}
}
