package sim

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/gameengine"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/triggertype"
)

// Tests that an end-of-turn Trigger fires once and is removed.
func TestFireEndOfTurn_FiresOnceAndRemoves(t *testing.T) {
	s := gameengine.NewFromState(nil)
	calls := 0
	s.CreateTrigger(NewCardTrigger(
		&card.CardState{Card: FakeRedAttack{}},
		triggertype.EndOfTurn,
		func(_ card.GameEngine, _ card.Logger, _ card.Trigger) { calls++ },
		nil,
	))
	s.FireEndOfTurn()
	if calls != 1 {
		t.Fatalf("handler calls = %d, want 1", calls)
	}
	if len(s.Triggers()) != 0 {
		t.Fatalf("triggers after fire = %d, want 0", len(s.Triggers()))
	}
}

// Tests that a non-matching TriggerType stays queued when end-of-turn fires.
func TestFireEndOfTurn_LeavesNonMatchingType(t *testing.T) {
	s := gameengine.NewFromState(nil)
	calls := 0
	s.CreateTrigger(NewCardTrigger(
		&card.CardState{Card: FakeRedAttack{}},
		triggertype.Attack,
		func(_ card.GameEngine, _ card.Logger, _ card.Trigger) { calls++ },
		nil,
	))
	s.FireEndOfTurn()
	if calls != 0 {
		t.Fatalf("handler calls = %d, want 0 (TriggerAttack should not fire from end-of-turn walk)", calls)
	}
	if len(s.Triggers()) != 1 {
		t.Fatalf("triggers after fire = %d, want 1 (non-matching trigger preserved)", len(s.Triggers()))
	}
}

// Tests that a handler appending a new trigger during fire queues it for a future fire
// walk rather than firing it on the current pass.
func TestFireEndOfTurn_HandlerAddTriggerSafeReentry(t *testing.T) {
	s := gameengine.NewFromState(nil)
	calls := 0
	s.CreateTrigger(NewCardTrigger(
		&card.CardState{Card: FakeRedAttack{}},
		triggertype.EndOfTurn,
		func(g card.GameEngine, _ card.Logger, _ card.Trigger) {
			calls++
			ts := g.(*gameengine.GameEngine)
			ts.CreateTrigger(NewCardTrigger(
				&card.CardState{Card: FakeRedAttack{}},
				triggertype.EndOfTurn,
				func(_ card.GameEngine, _ card.Logger, _ card.Trigger) { calls++ },
				nil,
			))
		},
		nil,
	))
	s.FireEndOfTurn()
	if calls != 1 {
		t.Fatalf("handler calls during first walk = %d, want 1 (handler-added trigger should not fire on the same pass)", calls)
	}
	if len(s.Triggers()) != 1 {
		t.Fatalf("triggers after fire = %d, want 1 (handler-added trigger preserved)", len(s.Triggers()))
	}
	s.FireEndOfTurn()
	if calls != 2 {
		t.Fatalf("handler calls after second walk = %d, want 2 (queued trigger fires on next pass)", calls)
	}
}
