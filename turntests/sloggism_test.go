package turntests

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card/cards"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/gameengine"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
)

// TestSloggism_NoAttackReturnsZero: no qualifying next attack card → +6 rider fizzles.
func TestSloggism_NoAttackReturnsZero(t *testing.T) {
	ge := gameengine.New()
	for _, c := range []card.Card{cards.SloggismRed{}, cards.SloggismYellow{}, cards.SloggismBlue{}} {
		ge.ResolveAttackStep(ge.Logger(), &card.CardState{Card: c})
		if got := ge.Value(); got != 0 {
			t.Errorf("%s: Play() = %d, want 0", c.Name(), got)
		}
	}
}

// TestSloggism_LowCostFilteredOut: a cost-1 attack is seen but the cost>=2 filter rejects it.
func TestSloggism_LowCostFilteredOut(t *testing.T) {
	ge := &gameengine.GameEngine{GameState: gameengine.GameStateBuilder().SetCardsRemaining([]*card.CardState{{Card: testutils.FakeRedAttack().WithCost(1)}}).Build()}
	ge.ResolveAttackStep(ge.Logger(), &card.CardState{Card: cards.SloggismRed{}})
	if got := ge.Value(); got != 0 {
		t.Errorf("Play() = %d, want 0 (cost 1 < 2)", got)
	}
}

// TestSloggism_HighCostReturnsBonus: first cost>=2 attack triggers the per-variant bonus
// (Red +6, Yellow +5, Blue +4).
func TestSloggism_HighCostReturnsBonus(t *testing.T) {
	cases := []struct {
		c    card.Card
		want int
	}{
		{cards.SloggismRed{}, 6},
		{cards.SloggismYellow{}, 5},
		{cards.SloggismBlue{}, 4},
	}
	for _, tc := range cases {
		target := &card.CardState{Card: testutils.FakeRedAttack().WithCost(2)}
		ge := &gameengine.GameEngine{GameState: gameengine.GameStateBuilder().SetCardsRemaining([]*card.CardState{target}).Build()}
		ge.ResolveAttackStep(ge.Logger(), &card.CardState{Card: tc.c})
		if got := ge.Value(); got != 0 {
			t.Errorf("%s: Play() = %d, want 0 (granter returns 0; +N rides on target'ge BonusAttack)", tc.c.Name(), got)
		}
		if target.BonusAttack != tc.want {
			t.Errorf("%s: target BonusAttack = %d, want %d", tc.c.Name(), target.BonusAttack, tc.want)
		}
	}
}
