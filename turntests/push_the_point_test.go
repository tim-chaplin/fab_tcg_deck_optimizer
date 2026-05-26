package turntests

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card/cards"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/deck"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/gameengine"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
)

// Tests that Push the Point gains +2{p} when the last attack on this attack turn hit. The
// preceding BlueAttack lands first, flipping LastAttackHit, so Push the Point reads it on
// the next attack step.
func TestPushThePoint_LastAttackHitGrantsBonus(t *testing.T) {
	for _, c := range []card.Card{cards.PushThePointRed{}, cards.PushThePointYellow{}, cards.PushThePointBlue{}} {
		d := deck.New(testutils.Hero{Intel: 4}, nil, nil)
		hand := []card.Card{testutils.FakeBlueAttack().
			WithPower(1).
			WithGoAgain(), c, testutils.FakeBlueResource(), testutils.FakeBlueResource()}
		summary := sim.EvalOneTurnForTesting(d, gameengine.GameStateBuilder().SetIncomingDamage(0).Build(), hand)
		want := testutils.FakeBlueAttack().
			WithPower(1).
			WithGoAgain().Attack() + c.Attack() + 2
		if summary.Value != want {
			t.Errorf("%s: Value = %d, want %d (BlueAttack + printed + 2{p} bonus)", c.Name(), summary.Value, want)
		}
	}
}

// Tests that Push the Point stays at printed power when no prior attack hit.
func TestPushThePoint_NoPriorHitNoBonus(t *testing.T) {
	for _, c := range []card.Card{cards.PushThePointRed{}, cards.PushThePointYellow{}, cards.PushThePointBlue{}} {
		d := deck.New(testutils.Hero{Intel: 4}, nil, nil)
		hand := []card.Card{c, testutils.FakeBlueResource()}
		summary := sim.EvalOneTurnForTesting(d, gameengine.GameStateBuilder().SetIncomingDamage(0).Build(), hand)
		if summary.Value != c.Attack() {
			t.Errorf("%s: Value = %d, want %d (printed power)", c.Name(), summary.Value, c.Attack())
		}
	}
}
