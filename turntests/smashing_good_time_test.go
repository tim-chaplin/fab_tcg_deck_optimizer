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

// Tests that an arsenal-played Smashing Good Time grants +N{p} to the next attack action card
// in hand — Red +3, Yellow +2, Blue +1. The buff lands on the target's BonusAttack so its
// EffectiveAttack folds it into the LikelyToHit check.
func TestSmashingGoodTime_ArsenalPlayBuffsNextAttackAction(t *testing.T) {
	cases := []struct {
		name      string
		c         card.Card
		wantValue int
	}{
		{"Red", cards.SmashingGoodTimeRed{}, 6},       // 3 base + 3 buff
		{"Yellow", cards.SmashingGoodTimeYellow{}, 5}, // 3 base + 2 buff
		{"Blue", cards.SmashingGoodTimeBlue{}, 4},     // 3 base + 1 buff
	}
	for _, tc := range cases {
		handAttack := testutils.FakeRedAttack().WithPower(3)
		d := deck.New(testutils.Hero{Intel: 4}, nil, nil)
		state := gameengine.GameStateBuilder().SetArsenal(tc.c).Build()
		summary := sim.EvalOneTurnForTesting(d, state, []card.Card{handAttack})
		if summary.Value != tc.wantValue {
			t.Errorf("%s arsenal: Value = %d, want %d (3 base + variant +N{p})", tc.name, summary.Value, tc.wantValue)
		}
	}
}

// Tests that a hand-played Smashing Good Time grants no buff — the +N{p} rider is gated on
// self.FromArsenal. The follow-up attack resolves at its base power only.
func TestSmashingGoodTime_HandPlayNoBuff(t *testing.T) {
	for _, c := range []card.Card{cards.SmashingGoodTimeRed{}, cards.SmashingGoodTimeYellow{}, cards.SmashingGoodTimeBlue{}} {
		handAttack := testutils.FakeRedAttack().WithPower(3)
		d := deck.New(testutils.Hero{Intel: 4}, nil, nil)
		summary := sim.EvalOneTurnForTesting(d, nil, []card.Card{c, handAttack})
		if summary.Value != 3 {
			t.Errorf("%s hand-play: Value = %d, want 3 (no arsenal → no +N{p} buff; only base attack)", c.Name(), summary.Value)
		}
	}
}
