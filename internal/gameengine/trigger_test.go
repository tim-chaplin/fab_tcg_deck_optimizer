package gameengine

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/aura"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/trigger"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/triggertype"
)

// Tests that an end-of-turn Trigger fires once and is removed.
func TestFireEndOfTurn_FiresOnceAndRemoves(t *testing.T) {
	ge := New()
	calls := 0
	ge.CreateTrigger(trigger.NewFromCard(
		stubCard{name: "src"},
		triggertype.EndOfTurn,
		func(_ card.GameEngine, _ card.Logger, _ card.Trigger) { calls++ },
		nil,
	))
	ge.FireTriggers(triggertype.EndOfTurn, nil)
	if calls != 1 {
		t.Fatalf("handler calls = %d, want 1", calls)
	}
	if len(ge.Triggers()) != 0 {
		t.Fatalf("triggers after fire = %d, want 0", len(ge.Triggers()))
	}
}

// Tests that a non-matching TriggerType stays queued when end-of-turn fires.
func TestFireEndOfTurn_LeavesNonMatchingType(t *testing.T) {
	ge := New()
	calls := 0
	ge.CreateTrigger(trigger.NewFromCard(
		stubCard{name: "src"},
		triggertype.Attack,
		func(_ card.GameEngine, _ card.Logger, _ card.Trigger) { calls++ },
		nil,
	))
	ge.FireTriggers(triggertype.EndOfTurn, nil)
	if calls != 0 {
		t.Fatalf("handler calls = %d, want 0 (TriggerAttack should not fire from end-of-turn walk)", calls)
	}
	if len(ge.Triggers()) != 1 {
		t.Fatalf("triggers after fire = %d, want 1 (non-matching trigger preserved)", len(ge.Triggers()))
	}
}

// Tests that a handler appending a new trigger during fire queues it for a future fire
// walk rather than firing it on the current pass.
func TestFireEndOfTurn_HandlerAddTriggerSafeReentry(t *testing.T) {
	ge := New()
	calls := 0
	ge.CreateTrigger(trigger.NewFromCard(
		stubCard{name: "src"},
		triggertype.EndOfTurn,
		func(engine card.GameEngine, _ card.Logger, _ card.Trigger) {
			calls++
			ts := engine.(*GameEngine)
			ts.CreateTrigger(trigger.NewFromCard(
				stubCard{name: "added"},
				triggertype.EndOfTurn,
				func(_ card.GameEngine, _ card.Logger, _ card.Trigger) { calls++ },
				nil,
			))
		},
		nil,
	))
	ge.FireTriggers(triggertype.EndOfTurn, nil)
	if calls != 1 {
		t.Fatalf("handler calls during first walk = %d, want 1 (handler-added trigger should not fire on the same pass)", calls)
	}
	if len(ge.Triggers()) != 1 {
		t.Fatalf("triggers after fire = %d, want 1 (handler-added trigger preserved)", len(ge.Triggers()))
	}
	ge.FireTriggers(triggertype.EndOfTurn, nil)
	if calls != 2 {
		t.Fatalf("handler calls after second walk = %d, want 2 (queued trigger fires on next pass)", calls)
	}
}

// Tests that ResetEphemeralState re-arms the OncePerTurn gate, so an aura that fired last
// turn can fire again on the next.
func TestResetEphemeralState_RearmsOncePerTurnAuras(t *testing.T) {
	ge := New()
	a := aura.NewFromCard(
		stubCard{name: "src"},
		triggertype.AttackAction,
		func(card.GameEngine, card.Logger, card.Aura) {},
		1,
		true,
	)
	a.SetFiredThisTurn(true)
	ge.CreateAura(a)
	ge.ResetEphemeralState()
	if ge.Auras()[0].FiredThisTurn() {
		t.Errorf("FiredThisTurn = true after ResetEphemeralState, want false (turn-boundary re-arm)")
	}
}
