package hand

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card/fake"
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
	state := &card.TurnState{EphemeralAttackTriggers: []card.EphemeralAttackTrigger{{
		Source: fake.RedAttack{},
		Handler: func(s *card.TurnState, target *card.CardState) int {
			calls++
			return s.AddPostTriggerLogEntry("test ephemeral fired", card.DisplayName(target.Card), 1)
		},
	}}}
	target1 := &card.CardState{Card: fake.RedAttack{}}
	target2 := &card.CardState{Card: fake.YellowAttack{}}

	fireEphemeralAttackTriggers(state, target1)
	if state.Value != 1 {
		t.Errorf("first fire Value = %d, want 1", state.Value)
	}
	if len(state.EphemeralAttackTriggers) != 0 {
		t.Errorf("after first fire EphemeralAttackTriggers = %+v, want empty (trigger consumed)",
			state.EphemeralAttackTriggers)
	}
	fireEphemeralAttackTriggers(state, target2)
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
	state := &card.TurnState{EphemeralAttackTriggers: []card.EphemeralAttackTrigger{{
		Source: fake.RedAttack{},
		Matches: func(target *card.CardState) bool {
			return target.Card.Types().Has(card.TypeRuneblade)
		},
		Handler: func(*card.TurnState, *card.CardState) int {
			calls++
			return 5
		},
	}}}
	generic := &card.CardState{Card: fake.RedAttack{}} // Generic, not Runeblade

	fireEphemeralAttackTriggers(state, generic)
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
	state := &card.TurnState{EphemeralAttackTriggers: []card.EphemeralAttackTrigger{{
		Source: fake.RedAttack{},
		Handler: func(s *card.TurnState, target *card.CardState) int {
			return s.AddPostTriggerLogEntry("test ephemeral fired", card.DisplayName(target.Card), 2)
		},
	}}}
	target := &card.CardState{Card: fake.YellowAttack{}} // Generic, but Matches=nil accepts all

	fireEphemeralAttackTriggers(state, target)
	if state.Value != 2 {
		t.Errorf("Value = %d, want 2 (nil Matches accepts any target)", state.Value)
	}
	if len(state.EphemeralAttackTriggers) != 0 {
		t.Errorf("trigger should be consumed on nil-Matches fire, got %+v",
			state.EphemeralAttackTriggers)
	}
}
