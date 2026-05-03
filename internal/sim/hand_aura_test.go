package sim_test

import (
	. "github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
)

// TestFireAttackActionAuras_FiresOnceWhenGated: a single OncePerTurn AttackAction
// trigger fires on the first call and is gated on the second within the same turn — its
// Count ticks only once, FiredThisTurn latches. Aura attack-action triggers are
// pre-triggers, so handlers credit Value through AddPreTriggerLogEntry.
func TestFireAttackActionAuras_FiresOnceWhenGated(t *testing.T) {
	aura := testutils.RedAttack{}
	calls := 0
	state := &TurnState{Auras: []Aura{{
		Self:        aura,
		TriggerType: TriggerAttackAction,
		Count:       3,
		OncePerTurn: true,
		Handler: func(s *TurnState, _ *Aura) int {
			calls++
			s.AddValue(1)
			s.LogPreTriggerf("TestCard", 1, "test trigger fired")
			return 1
		},
	}}}
	trigger := testutils.RedAttack{}
	FireAttackActionAuras(state, trigger)
	if state.Value != 1 {
		t.Errorf("first fire Value = %d, want 1", state.Value)
	}
	FireAttackActionAuras(state, trigger)
	if state.Value != 1 {
		t.Errorf("second fire Value = %d, want 1 (OncePerTurn gate kept second fire from crediting)", state.Value)
	}
	if calls != 1 {
		t.Errorf("handler call count = %d, want 1 (gate prevented second call)", calls)
	}
	if len(state.Auras) != 1 || state.Auras[0].Count != 2 {
		t.Errorf("trigger state = %+v, want one entry with Count=2", state.Auras)
	}
	if !state.Auras[0].FiredThisTurn {
		t.Errorf("FiredThisTurn = false, want true (single fire latched)")
	}
}

// TestFireAttackActionAuras_GraveyardsExhaustedAura: a Count=1 trigger fires once, hits
// Count=0, and the sim drops it from Auras and graveyards Self.
func TestFireAttackActionAuras_GraveyardsExhaustedAura(t *testing.T) {
	aura := testutils.RedAttack{}
	state := &TurnState{Auras: []Aura{{
		Self:        aura,
		TriggerType: TriggerAttackAction,
		Count:       1,
		Handler:     func(*TurnState, *Aura) int { return 1 },
	}}}
	FireAttackActionAuras(state, testutils.RedAttack{})
	if len(state.Auras) != 0 {
		t.Errorf("Auras = %+v, want empty (Count hit zero)", state.Auras)
	}
	g := state.Graveyard()
	if len(g) != 1 || g[0] != aura {
		t.Errorf("Graveyard = %v, want [aura]", g)
	}
}

// TestFireAttackActionAuras_PassesThroughNonAttackActionTriggers: a TriggerStartOfTurn
// trigger is left untouched by FireAttackActionAuras — only AttackAction-typed entries
// fire here.
func TestFireAttackActionAuras_PassesThroughNonAttackActionTriggers(t *testing.T) {
	aura := testutils.RedAttack{}
	calls := 0
	state := &TurnState{Auras: []Aura{{
		Self:        aura,
		TriggerType: TriggerStartOfTurn,
		Count:       1,
		Handler:     func(*TurnState, *Aura) int { calls++; return 5 },
	}}}
	FireAttackActionAuras(state, testutils.RedAttack{})
	if state.Value != 0 {
		t.Errorf("Value = %d, want 0 (start-of-turn trigger doesn't fire on attack action)", state.Value)
	}
	if calls != 0 {
		t.Errorf("handler call count = %d, want 0", calls)
	}
	if len(state.Auras) != 1 || state.Auras[0].Count != 1 {
		t.Errorf("trigger should be untouched, got %+v", state.Auras)
	}
}
