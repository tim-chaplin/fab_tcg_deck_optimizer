package turntests

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card/cards"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/gameengine"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
)

// TestNimblism_NoAttackReturnsZero: no qualifying next attack card → +3 rider fizzles.
func TestNimblism_NoAttackReturnsZero(t *testing.T) {
	ge := gameengine.New()
	for _, c := range []card.Card{cards.NimblismRed{}, cards.NimblismYellow{}, cards.NimblismBlue{}} {
		ge.ResolveChainStep(ge.Logger(), &card.CardState{Card: c})
		if got := ge.Value(); got != 0 {
			t.Errorf("%s: Play() = %d, want 0", c.Name(), got)
		}
	}
}

// TestNimblism_HighCostFilteredOut: a cost-2 attack is seen but the cost<=1 filter rejects it.
func TestNimblism_HighCostFilteredOut(t *testing.T) {
	ge := &gameengine.GameEngine{GameState: gameengine.GameStateBuilder().SetCardsRemaining([]*card.CardState{{Card: testutils.FakeRedAttack()}}).Build()}
	ge.ResolveChainStep(ge.Logger(), &card.CardState{Card: cards.NimblismRed{}})
	if got := ge.Value(); got != 0 {
		t.Errorf("Play() = %d, want 0 (cost 2 > 1)", got)
	}
}

// TestNimblism_LowCostReturnsBonus: first cost<=1 attack triggers the per-variant bonus
// (Red +3, Yellow +2, Blue +1).
func TestNimblism_LowCostReturnsBonus(t *testing.T) {
	cases := []struct {
		c    card.Card
		want int
	}{
		{cards.NimblismRed{}, 3},
		{cards.NimblismYellow{}, 2},
		{cards.NimblismBlue{}, 1},
	}
	for _, tc := range cases {
		target := &card.CardState{Card: testutils.FakeRedAttack()}
		ge := &gameengine.GameEngine{GameState: gameengine.GameStateBuilder().SetCardsRemaining([]*card.CardState{target}).Build()}
		ge.ResolveChainStep(ge.Logger(), &card.CardState{Card: tc.c})
		if got := ge.Value(); got != 0 {
			t.Errorf("%s: Play() = %d, want 0 (granter returns 0; +N rides on target'ge BonusAttack)", tc.c.Name(), got)
		}
		if target.BonusAttack != tc.want {
			t.Errorf("%s: target BonusAttack = %d, want %d", tc.c.Name(), target.BonusAttack, tc.want)
		}
	}
}
