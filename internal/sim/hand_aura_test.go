package sim

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/aura"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/gameengine"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/triggertype"
)

// Tests that an OncePerTurn CardOrAbility aura fires on the first call and is gated by
// FiredThisTurn on the second within the same turn.
func TestFireCardOrAbilityAuras_FiresOnceWhenGated(t *testing.T) {
	src := FakeRedAttack{}
	calls := 0
	ge := gameengine.New()
	ge.AppendAura(aura.NewFromCard(
		src,
		triggertype.CardOrAbility,
		func(ge card.GameEngine, l card.Logger, _ card.Aura) {
			calls++
			ge.AddValue(1)
			l.AppendPreTriggerf("TestCard", 1, "test trigger fired")
		},
		3,
		true, // oncePerTurn
		nil,
	))
	trigger := FakeRedAttack{}
	ge.FireTriggers(triggertype.CardOrAbility, trigger)
	if ge.Value() != 1 {
		t.Errorf("first fire Value = %d, want 1", ge.Value())
	}
	ge.FireTriggers(triggertype.CardOrAbility, trigger)
	if ge.Value() != 1 {
		t.Errorf("second fire Value = %d, want 1 (OncePerTurn gate kept second fire from crediting)", ge.Value())
	}
	if calls != 1 {
		t.Errorf("handler call count = %d, want 1 (gate prevented second call)", calls)
	}
	if len(ge.Auras()) != 1 || ge.Auras()[0].Count() != 3 {
		t.Errorf("aura ge = %+v, want one entry with Count=3", ge.Auras())
	}
	if !ge.Auras()[0].FiredThisTurn() {
		t.Errorf("FiredThisTurn = false, want true (single fire latched)")
	}
}

// Tests that a handler calling Destroy drops the entry from Auras and lands Self in the
// graveyard.
func TestFireCardOrAbilityAuras_GraveyardsExhaustedAura(t *testing.T) {
	src := FakeRedAttack{}
	ge := gameengine.New()
	ge.AppendAura(aura.NewFromCard(
		src,
		triggertype.CardOrAbility,
		func(ge card.GameEngine, _ card.Logger, a card.Aura) {
			ge.AddValue(1)
			a.Destroy(true)
		},
		1,
		false,
		nil,
	))
	ge.FireTriggers(triggertype.CardOrAbility, FakeRedAttack{})
	if len(ge.Auras()) != 0 {
		t.Errorf("Auras = %+v, want empty (handler called Destroy)", ge.Auras())
	}
	g := ge.Graveyard()
	if len(g) != 1 || g[0] != src {
		t.Errorf("Graveyard = %v, want [src]", g)
	}
}

// Tests that a StartOfTurn aura is left untouched by a CardOrAbility fire.
func TestFireCardOrAbilityAuras_PassesThroughNonCardOrAbilityTriggers(t *testing.T) {
	src := FakeRedAttack{}
	calls := 0
	ge := gameengine.New()
	ge.AppendAura(aura.NewFromCard(
		src,
		triggertype.StartOfTurn,
		func(card.GameEngine, card.Logger, card.Aura) { calls++ },
		1,
		false,
		nil,
	))
	ge.FireTriggers(triggertype.CardOrAbility, FakeRedAttack{})
	if ge.Value() != 0 {
		t.Errorf("Value = %d, want 0 (start-of-turn aura doesn't fire on a card played)", ge.Value())
	}
	if calls != 0 {
		t.Errorf("handler call count = %d, want 0", calls)
	}
	if len(ge.Auras()) != 1 || ge.Auras()[0].Count() != 1 {
		t.Errorf("aura should be untouched, got %+v", ge.Auras())
	}
}
