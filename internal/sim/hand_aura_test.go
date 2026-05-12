package sim

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/gameengine"
)

// Tests that an OncePerTurn AttackAction aura fires on the first call and is gated by
// FiredThisTurn on the second within the same turn.
func TestFireAttackActionAuras_FiresOnceWhenGated(t *testing.T) {
	aura := FakeRedAttack{}
	calls := 0
	state := gameengine.New()
	state.AppendAura(NewCardAura(
		&card.CardState{Card: aura},
		gameengine.TriggerAttackAction,
		func(g card.GameEngine, l card.Logger, _ card.Aura) {
			calls++
			g.AddValue(1)
			l.AppendPreTriggerf("TestCard", 1, "test trigger fired")
		},
		3,
		true, // oncePerTurn
	))
	trigger := FakeRedAttack{}
	state.FireAttackAction(trigger)
	if state.Value() != 1 {
		t.Errorf("first fire Value = %d, want 1", state.Value())
	}
	state.FireAttackAction(trigger)
	if state.Value() != 1 {
		t.Errorf("second fire Value = %d, want 1 (OncePerTurn gate kept second fire from crediting)", state.Value())
	}
	if calls != 1 {
		t.Errorf("handler call count = %d, want 1 (gate prevented second call)", calls)
	}
	if len(state.Auras()) != 1 || state.Auras()[0].Count() != 3 {
		t.Errorf("aura state = %+v, want one entry with Count=3", state.Auras())
	}
	if !state.Auras()[0].FiredThisTurn() {
		t.Errorf("FiredThisTurn = false, want true (single fire latched)")
	}
}

// Tests that a handler calling Destroy drops the entry from Auras and lands Self in the
// graveyard.
func TestFireAttackActionAuras_GraveyardsExhaustedAura(t *testing.T) {
	aura := FakeRedAttack{}
	state := gameengine.New()
	state.AppendAura(NewCardAura(
		&card.CardState{Card: aura},
		gameengine.TriggerAttackAction,
		func(g card.GameEngine, _ card.Logger, a card.Aura) {
			g.AddValue(1)
			a.Destroy(true)
		},
		1,
		false,
	))
	state.FireAttackAction(FakeRedAttack{})
	if len(state.Auras()) != 0 {
		t.Errorf("Auras = %+v, want empty (handler called Destroy)", state.Auras())
	}
	g := state.GraveyardRaw()
	if len(g) != 1 || g[0] != aura {
		t.Errorf("Graveyard = %v, want [aura]", g)
	}
}

// Tests that a TriggerStartOfTurn aura is left untouched by FireAttackAction.
func TestFireAttackActionAuras_PassesThroughNonAttackActionTriggers(t *testing.T) {
	aura := FakeRedAttack{}
	calls := 0
	state := gameengine.New()
	state.AppendAura(NewCardAura(
		&card.CardState{Card: aura},
		gameengine.TriggerStartOfTurn,
		func(card.GameEngine, card.Logger, card.Aura) { calls++ },
		1,
		false,
	))
	state.FireAttackAction(FakeRedAttack{})
	if state.Value() != 0 {
		t.Errorf("Value = %d, want 0 (start-of-turn aura doesn't fire on attack action)", state.Value())
	}
	if calls != 0 {
		t.Errorf("handler call count = %d, want 0", calls)
	}
	if len(state.Auras()) != 1 || state.Auras()[0].Count() != 1 {
		t.Errorf("aura should be untouched, got %+v", state.Auras())
	}
}
