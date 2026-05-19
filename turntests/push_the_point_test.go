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

// Tests that Push the Point gains +2{p} when the last attack on this chain hit.
func TestPushThePoint_LastAttackHitGrantsBonus(t *testing.T) {
	for _, c := range []card.Card{cards.PushThePointRed{}, cards.PushThePointYellow{}, cards.PushThePointBlue{}} {
		ge := &gameengine.GameEngine{GameState: gameengine.GameStateBuilder().Build()}
		ge.SetLastAttackHit(true)
		pc := &card.CardState{Card: c}
		ge.ResolveChainStep(ge.Logger(), pc)
		want := c.Attack() + 2
		if ge.Value() != want {
			t.Errorf("%s: Value = %d, want %d (printed + 2{p} bonus)", c.Name(), ge.Value(), want)
		}
	}
}

// Tests that Push the Point stays at printed power when no prior attack hit.
func TestPushThePoint_NoPriorHitNoBonus(t *testing.T) {
	for _, c := range []card.Card{cards.PushThePointRed{}, cards.PushThePointYellow{}, cards.PushThePointBlue{}} {
		ge := gameengine.New()
		pc := &card.CardState{Card: c}
		ge.ResolveChainStep(ge.Logger(), pc)
		if ge.Value() != c.Attack() {
			t.Errorf("%s: Value = %d, want %d (printed power)", c.Name(), ge.Value(), c.Attack())
		}
	}
}

// Tests that finalizeActiveAttack records the hit and Push the Point reads it on the
// next chain step.
func TestPushThePoint_ChainRunnerWiresPriorHit(t *testing.T) {
	d := deck.New(testutils.Hero{Intel: 4}, nil, nil)
	summary := sim.EvalOneTurnForTesting(
		d,
		gameengine.GameStateBuilder().SetIncomingDamage(0).Build(),
		[]card.Card{testutils.BlueAttack{}, cards.PushThePointRed{}, testutils.BluePitch{}},
	)
	if summary.Value != 7 {
		t.Errorf("Value = %d, want 7 (BlueAttack 1 + PushThePoint 4 + LastAttackHit bonus 2)", summary.Value)
	}
}
