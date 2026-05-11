package sim

import (
	"testing"
)

// Tests that an OncePerTurn AttackAction trigger fires on the first call and is gated by
// FiredThisTurn on the second within the same turn.
func TestFireAttackActionAuras_FiresOnceWhenGated(t *testing.T) {
	aura := FakeRedAttack{}
	calls := 0
	state := &TurnState{Auras: []Aura{{
		Trigger: Trigger{
			TriggerType: TriggerAttackAction,
			Handler: func(s *TurnState, l Logger, _ *Trigger, _ *Aura) {
				calls++
				s.AddValue(1)
				l.AppendPreTriggerf("TestCard", 1, "test trigger fired")
			},
		},
		Self:        CardOrTokenType{Card: aura},
		Count:       3,
		OncePerTurn: true,
	}}}
	trigger := FakeRedAttack{}
	fireAttackActionAuras(state, trigger)
	if state.Value != 1 {
		t.Errorf("first fire Value = %d, want 1", state.Value)
	}
	fireAttackActionAuras(state, trigger)
	if state.Value != 1 {
		t.Errorf("second fire Value = %d, want 1 (OncePerTurn gate kept second fire from crediting)", state.Value)
	}
	if calls != 1 {
		t.Errorf("handler call count = %d, want 1 (gate prevented second call)", calls)
	}
	if len(state.Auras) != 1 || state.Auras[0].Count != 3 {
		t.Errorf("trigger state = %+v, want one entry with Count=3 (sim never mutates Count)", state.Auras)
	}
	if !state.Auras[0].FiredThisTurn {
		t.Errorf("FiredThisTurn = false, want true (single fire latched)")
	}
}

// TestFireAttackActionAuras_GraveyardsExhaustedAura: a handler that calls DestroyAura
// drops the entry from Auras and lands Self in the graveyard.
func TestFireAttackActionAuras_GraveyardsExhaustedAura(t *testing.T) {
	aura := FakeRedAttack{}
	state := &TurnState{Auras: []Aura{{
		Trigger: Trigger{
			TriggerType: TriggerAttackAction,
			Handler: func(s *TurnState, _ Logger, _ *Trigger, a *Aura) {
				s.AddValue(1)
				s.DestroyAura(a, true)
			},
		},
		Self:  CardOrTokenType{Card: aura},
		Count: 1,
	}}}
	fireAttackActionAuras(state, FakeRedAttack{})
	if len(state.Auras) != 0 {
		t.Errorf("Auras = %+v, want empty (handler called DestroyAura)", state.Auras)
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
	aura := FakeRedAttack{}
	calls := 0
	state := &TurnState{Auras: []Aura{{
		Trigger: Trigger{
			TriggerType: TriggerStartOfTurn,
			Handler:     func(*TurnState, Logger, *Trigger, *Aura) { calls++ },
		},
		Self:  CardOrTokenType{Card: aura},
		Count: 1,
	}}}
	fireAttackActionAuras(state, FakeRedAttack{})
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
