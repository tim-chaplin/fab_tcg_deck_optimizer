package turntests

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/cards"

	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/gameengine"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/triggertype"
)

// Tests that Play credits 0 immediately and registers an OncePerTurn TriggerAttackAction
// aura with the variant's Count.
func TestMaleficIncantation_PlayRegistersAttackActionTrigger(t *testing.T) {
	cases := []struct {
		c card.Card
		n int
	}{
		{cards.MaleficIncantationRed{}, 3},
		{cards.MaleficIncantationYellow{}, 2},
		{cards.MaleficIncantationBlue{}, 1},
	}
	for _, tc := range cases {
		ge := gameengine.New()
		ge.ResolveChainStep(ge.Logger(), &card.CardState{Card: tc.c})
		if got := ge.Value(); got != 0 {
			t.Errorf("%s: Play() = %d, want 0 (rune comes from trigger, not Play)", tc.c.Name(), got)
		}
		if !ge.AuraCreated() {
			t.Errorf("%s: AuraCreated should be set", tc.c.Name())
		}
		if ge.RunechantCount() != 0 {
			t.Errorf("%s: Runechants = %d, want 0 (trigger not yet fired)", tc.c.Name(), ge.RunechantCount())
		}
		if len(ge.Auras()) != 1 {
			t.Fatalf("%s: Auras len = %d, want 1", tc.c.Name(), len(ge.Auras()))
		}
		tr := ge.Auras()[0]
		if tr.TriggerType() != triggertype.AttackAction {
			t.Errorf("%s: trigger Type = %d, want TriggerAttackAction", tc.c.Name(), tr.TriggerType())
		}
		if !tr.OncePerTurn() {
			t.Errorf("%s: OncePerTurn = false, want true", tc.c.Name())
		}
		if tr.Count() != tc.n {
			t.Errorf("%s: Count = %d, want %d (one per verse counter)", tc.c.Name(), tr.Count(), tc.n)
		}
	}
}

// Tests that one handler invocation creates one Runechant and credits 1 damage.
func TestMaleficIncantation_HandlerCreatesOneRunechantPerFire(t *testing.T) {
	for _, c := range []card.Card{cards.MaleficIncantationRed{}, cards.MaleficIncantationYellow{}, cards.MaleficIncantationBlue{}} {
		ge := gameengine.New()
		ge.ResolveChainStep(ge.Logger(), &card.CardState{Card: c})
		chain := gameengine.New()
		chain.SetTriggeringCard(c)
		chain.CreateAura(ge.Auras()[0])
		chain.FireAttackAction(c)
		if chain.Value() != 1 {
			t.Errorf("%s: handler Value = %d, want 1", c.Name(), chain.Value())
		}
		if chain.RunechantCount() != 1 {
			t.Errorf("%s: Runechants = %d, want 1 (handler creates one live rune)", c.Name(), chain.RunechantCount())
		}
	}
}
