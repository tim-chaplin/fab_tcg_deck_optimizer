package sim

import (
	"testing"
)

// Tests that an end-of-turn Trigger fires once and is removed.
func TestFireEndOfTurn_FiresOnceAndRemoves(t *testing.T) {
	s := NewTurnState(nil, nil)
	calls := 0
	s.AddTrigger(Trigger{
		TriggerType: TriggerEndOfTurn,
		Handler:     func(_ *TurnState, _ Logger, _ *Trigger, _ *Aura) { calls++ },
	})
	FireEndOfTurn(s)
	if calls != 1 {
		t.Fatalf("handler calls = %d, want 1", calls)
	}
	if len(s.Triggers) != 0 {
		t.Fatalf("triggers after fire = %d, want 0", len(s.Triggers))
	}
}

// Tests that a non-matching TriggerType stays queued when end-of-turn fires.
func TestFireEndOfTurn_LeavesNonMatchingType(t *testing.T) {
	s := NewTurnState(nil, nil)
	calls := 0
	s.AddTrigger(Trigger{
		TriggerType: TriggerAttack,
		Handler:     func(_ *TurnState, _ Logger, _ *Trigger, _ *Aura) { calls++ },
	})
	FireEndOfTurn(s)
	if calls != 0 {
		t.Fatalf("handler calls = %d, want 0 (TriggerAttack should not fire from end-of-turn walk)", calls)
	}
	if len(s.Triggers) != 1 {
		t.Fatalf("triggers after fire = %d, want 1 (non-matching trigger preserved)", len(s.Triggers))
	}
}

// Tests that a handler calling AddTrigger during fire queues the new entry for a future
// fire walk rather than firing it on the current pass.
func TestFireEndOfTurn_HandlerAddTriggerSafeReentry(t *testing.T) {
	s := NewTurnState(nil, nil)
	calls := 0
	s.AddTrigger(Trigger{
		TriggerType: TriggerEndOfTurn,
		Handler: func(s *TurnState, _ Logger, _ *Trigger, _ *Aura) {
			calls++
			s.AddTrigger(Trigger{
				TriggerType: TriggerEndOfTurn,
				Handler:     func(_ *TurnState, _ Logger, _ *Trigger, _ *Aura) { calls++ },
			})
		},
	})
	FireEndOfTurn(s)
	if calls != 1 {
		t.Fatalf("handler calls during first walk = %d, want 1 (handler-added trigger should not fire on the same pass)", calls)
	}
	if len(s.Triggers) != 1 {
		t.Fatalf("triggers after fire = %d, want 1 (handler-added trigger preserved)", len(s.Triggers))
	}
	FireEndOfTurn(s)
	if calls != 2 {
		t.Fatalf("handler calls after second walk = %d, want 2 (queued trigger fires on next pass)", calls)
	}
}
