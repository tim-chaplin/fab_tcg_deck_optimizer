package turntests

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/cards"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/gameengine"
)

// TestMinnowism_NoAttackReturnsZero: no qualifying next attack card → +3 rider fizzles.
func TestMinnowism_NoAttackReturnsZero(t *testing.T) {
	s := gameengine.New()
	for _, c := range []card.Card{cards.MinnowismRed{}, cards.MinnowismYellow{}, cards.MinnowismBlue{}} {
		s.ResolveChainStep(s.Logger(), &card.CardState{Card: c})
		if got := s.Value(); got != 0 {
			t.Errorf("%s: Play() = %d, want 0", c.Name(), got)
		}
	}
}

// TestMinnowism_HighPowerFilteredOut: a power-4 attack is seen but the power<=3 filter rejects it,
// so the rider fizzles without falling through to a later match.
func TestMinnowism_HighPowerFilteredOut(t *testing.T) {
	s := &gameengine.GameEngine{GameState: gameengine.GameStateBuilder().SetCardsRemaining([]*card.CardState{{Card: testutils.GenericAttack(0, 4)}}).Build()}
	s.ResolveChainStep(s.Logger(), &card.CardState{Card: cards.MinnowismRed{}})
	if got := s.Value(); got != 0 {
		t.Errorf("Play() = %d, want 0 (power 4 > 3)", got)
	}
}

// TestMinnowism_LowPowerReturnsBonus: first power<=3 attack triggers the per-variant bonus
// (Red +3, Yellow +2, Blue +1).
func TestMinnowism_LowPowerReturnsBonus(t *testing.T) {
	cases := []struct {
		c    card.Card
		want int
	}{
		{cards.MinnowismRed{}, 3},
		{cards.MinnowismYellow{}, 2},
		{cards.MinnowismBlue{}, 1},
	}
	for _, tc := range cases {
		target := &card.CardState{Card: testutils.GenericAttack(0, 3)}
		s := &gameengine.GameEngine{GameState: gameengine.GameStateBuilder().SetCardsRemaining([]*card.CardState{target}).Build()}
		s.ResolveChainStep(s.Logger(), &card.CardState{Card: tc.c})
		if got := s.Value(); got != 0 {
			t.Errorf("%s: Play() = %d, want 0 (granter returns 0; +N rides on target's BonusAttack)", tc.c.Name(), got)
		}
		if target.BonusAttack != tc.want {
			t.Errorf("%s: target BonusAttack = %d, want %d", tc.c.Name(), target.BonusAttack, tc.want)
		}
	}
}
