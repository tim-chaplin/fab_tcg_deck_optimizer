package sim_test

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	. "github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
)

// TestFireEphemeralAttackTriggers_SingleFireDropsFromList: ephemeral triggers fire at most
// once — the first matching attack consumes the trigger and removes it from the slice. A
// later attack can't re-fire it even when its Matches predicate would have accepted. Pins
// the core invariant that distinguishes EphemeralAttackTrigger from AuraTrigger's counted
// fires: no OncePerTurn/Count bookkeeping needed because consumption IS the signal.
// Ephemeral attack triggers are post-triggers, so handlers credit Value through
// AddPostTriggerLogEntry.
func TestFireEphemeralAttackTriggers_SingleFireDropsFromList(t *testing.T) {
	calls := 0
	state := &TurnState{EphemeralAttackTriggers: []EphemeralAttackTrigger{{
		Source: testutils.RedAttack{},
		Handler: func(s *TurnState, target *CardState) int {
			calls++
			return s.AddPostTriggerLogEntry("test ephemeral fired", DisplayName(target.Card), 1)
		},
	}}}
	target1 := &CardState{Card: testutils.RedAttack{}}
	target2 := &CardState{Card: testutils.YellowAttack{}}

	FireEphemeralAttackTriggers(state, target1)
	if state.Value != 1 {
		t.Errorf("first fire Value = %d, want 1", state.Value)
	}
	if len(state.EphemeralAttackTriggers) != 0 {
		t.Errorf("after first fire EphemeralAttackTriggers = %+v, want empty (trigger consumed)",
			state.EphemeralAttackTriggers)
	}
	FireEphemeralAttackTriggers(state, target2)
	if state.Value != 1 {
		t.Errorf("second fire Value = %d, want 1 (trigger already consumed; no second credit)", state.Value)
	}
	if calls != 1 {
		t.Errorf("handler call count = %d, want 1 (a second attack must not re-fire a consumed trigger)",
			calls)
	}
}

// TestFireEphemeralAttackTriggers_NonMatchingTargetLeavesTriggerInPlace: a target that
// fails the Matches predicate doesn't consume the trigger; it stays in the slice waiting
// for a later attack action that does match. Pins the "waits for the right attack"
// semantics that let Mauvrion Skies skip past a Generic attack and land on the next
// Runeblade attack action.
func TestFireEphemeralAttackTriggers_NonMatchingTargetLeavesTriggerInPlace(t *testing.T) {
	calls := 0
	state := &TurnState{EphemeralAttackTriggers: []EphemeralAttackTrigger{{
		Source: testutils.RedAttack{},
		Matches: func(target *CardState) bool {
			return target.Card.Types().Has(card.TypeRuneblade)
		},
		Handler: func(*TurnState, *CardState) int {
			calls++
			return 5
		},
	}}}
	generic := &CardState{Card: testutils.RedAttack{}} // Generic, not Runeblade

	FireEphemeralAttackTriggers(state, generic)
	if state.Value != 0 {
		t.Errorf("non-matching fire Value = %d, want 0", state.Value)
	}
	if calls != 0 {
		t.Errorf("handler ran %d times on a rejected target, want 0", calls)
	}
	if len(state.EphemeralAttackTriggers) != 1 {
		t.Errorf("trigger should still be registered, got %+v", state.EphemeralAttackTriggers)
	}
}

// TestFireEphemeralAttackTriggers_NilMatchesAcceptsAnyTarget: Matches=nil is the "match any
// attack action" shortcut so simple one-shot triggers don't have to spell out an
// accept-everything predicate.
func TestFireEphemeralAttackTriggers_NilMatchesAcceptsAnyTarget(t *testing.T) {
	state := &TurnState{EphemeralAttackTriggers: []EphemeralAttackTrigger{{
		Source: testutils.RedAttack{},
		Handler: func(s *TurnState, target *CardState) int {
			return s.AddPostTriggerLogEntry("test ephemeral fired", DisplayName(target.Card), 2)
		},
	}}}
	target := &CardState{Card: testutils.YellowAttack{}} // Generic, but Matches=nil accepts all

	FireEphemeralAttackTriggers(state, target)
	if state.Value != 2 {
		t.Errorf("Value = %d, want 2 (nil Matches accepts any target)", state.Value)
	}
	if len(state.EphemeralAttackTriggers) != 0 {
		t.Errorf("trigger should be consumed on nil-Matches fire, got %+v",
			state.EphemeralAttackTriggers)
	}
}
