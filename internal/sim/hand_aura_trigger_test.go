package sim_test

import (
	. "github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
)

// TestFireAttackActionTriggers_FiresOnceWhenGated: a single OncePerTurn AttackAction
// trigger fires on the first call and is gated on the second within the same turn — its
// Count ticks only once, FiredThisTurn latches. Aura attack-action triggers are
// pre-triggers, so handlers credit Value through AddPreTriggerLogEntry.
func TestFireAttackActionTriggers_FiresOnceWhenGated(t *testing.T) {
	aura := testutils.RedAttack{}
	calls := 0
	state := &TurnState{AuraTriggers: []AuraTrigger{{
		Self:        aura,
		Type:        TriggerAttackAction,
		Count:       3,
		OncePerTurn: true,
		Handler: func(s *TurnState, _ *AuraTrigger) int {
			calls++
			return s.AddPreTriggerLogEntry("test trigger fired", "TestCard", 1)
		},
	}}}
	trigger := testutils.RedAttack{}
	FireAttackActionTriggers(state, trigger)
	if state.Value != 1 {
		t.Errorf("first fire Value = %d, want 1", state.Value)
	}
	FireAttackActionTriggers(state, trigger)
	if state.Value != 1 {
		t.Errorf("second fire Value = %d, want 1 (OncePerTurn gate kept second fire from crediting)", state.Value)
	}
	if calls != 1 {
		t.Errorf("handler call count = %d, want 1 (gate prevented second call)", calls)
	}
	if len(state.AuraTriggers) != 1 || state.AuraTriggers[0].Count != 2 {
		t.Errorf("trigger state = %+v, want one entry with Count=2", state.AuraTriggers)
	}
	if !state.AuraTriggers[0].FiredThisTurn {
		t.Errorf("FiredThisTurn = false, want true (single fire latched)")
	}
}

// TestFireAttackActionTriggers_GraveyardsExhaustedAura: a Count=1 trigger fires once, hits
// Count=0, and the sim drops it from AuraTriggers and graveyards Self.
func TestFireAttackActionTriggers_GraveyardsExhaustedAura(t *testing.T) {
	aura := testutils.RedAttack{}
	state := &TurnState{AuraTriggers: []AuraTrigger{{
		Self:    aura,
		Type:    TriggerAttackAction,
		Count:   1,
		Handler: func(*TurnState, *AuraTrigger) int { return 1 },
	}}}
	FireAttackActionTriggers(state, testutils.RedAttack{})
	if len(state.AuraTriggers) != 0 {
		t.Errorf("AuraTriggers = %+v, want empty (Count hit zero)", state.AuraTriggers)
	}
	g := state.Graveyard()
	if len(g) != 1 || g[0] != aura {
		t.Errorf("Graveyard = %v, want [aura]", g)
	}
}

// TestFireAttackActionTriggers_PassesThroughNonAttackActionTriggers: a TriggerStartOfTurn
// trigger is left untouched by FireAttackActionTriggers — only AttackAction-typed entries
// fire here.
func TestFireAttackActionTriggers_PassesThroughNonAttackActionTriggers(t *testing.T) {
	aura := testutils.RedAttack{}
	calls := 0
	state := &TurnState{AuraTriggers: []AuraTrigger{{
		Self:    aura,
		Type:    TriggerStartOfTurn,
		Count:   1,
		Handler: func(*TurnState, *AuraTrigger) int { calls++; return 5 },
	}}}
	FireAttackActionTriggers(state, testutils.RedAttack{})
	if state.Value != 0 {
		t.Errorf("Value = %d, want 0 (start-of-turn trigger doesn't fire on attack action)", state.Value)
	}
	if calls != 0 {
		t.Errorf("handler call count = %d, want 0", calls)
	}
	if len(state.AuraTriggers) != 1 || state.AuraTriggers[0].Count != 1 {
		t.Errorf("trigger should be untouched, got %+v", state.AuraTriggers)
	}
}
