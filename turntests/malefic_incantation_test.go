package turntests

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card/cards"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/deck"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/gameengine"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/triggertype"
)

// Tests that Play credits 0 immediately and registers an OncePerTurn CardOrAbility aura
// with the variant's Count.
func TestMaleficIncantation_PlayRegistersCardOrAbilityTrigger(t *testing.T) {
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
		ge.ResolveAttackStep(ge.Logger(), &card.CardState{Card: tc.c})
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
		if tr.TriggerType() != triggertype.CardOrAbility {
			t.Errorf("%s: trigger Type = %d, want CardOrAbility", tc.c.Name(), tr.TriggerType())
		}
		if !tr.OncePerTurn() {
			t.Errorf("%s: OncePerTurn = false, want true", tc.c.Name())
		}
		if tr.Count() != tc.n {
			t.Errorf("%s: Count = %d, want %d (one per verse counter)", tc.c.Name(), tr.Count(), tc.n)
		}
	}
}

// Tests that a carried Malefic Incantation fires its handler when an attack action card
// resolves on the next turn, producing exactly one Runechant. The OncePerTurn gate keeps
// it to a single fire per turn even if the aura has verse counters left.
func TestMaleficIncantation_OncePerTurnFiresOnAttackActionCreatesOneRune(t *testing.T) {
	for _, c := range []card.Card{cards.MaleficIncantationRed{}, cards.MaleficIncantationYellow{}, cards.MaleficIncantationBlue{}} {
		prior := gameengine.GameStateBuilder().CreateAuraFromCard(c).Build()
		d := deck.New(testutils.Hero{Intel: 4}, nil, nil)
		hand := []card.Card{testutils.FakeRedAttack().WithTypes(card.TypeRuneblade)}

		summary := sim.EvalOneTurnForTesting(d, prior, hand)

		if got := summary.State.RunechantCount(); got != 1 {
			t.Errorf("%s: Runechants = %d, want 1 (one fire per turn even with multiple verses)", c.Name(), got)
		}
	}
}
