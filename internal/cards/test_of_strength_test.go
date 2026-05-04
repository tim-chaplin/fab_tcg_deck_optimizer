package cards

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
)

// Tests that Test of Strength creates a Gold token when the deck-top wins the clash
// (top power ≥ 6 vs opponent's modelled 5).
func TestTestOfStrength_WinCreatesGold(t *testing.T) {
	for _, power := range []int{6, 7} {
		s := sim.NewTurnState([]sim.Card{testutils.GenericAttack(0, power)}, nil)
		(TestOfStrengthRed{}).Play(s, &sim.CardState{Card: TestOfStrengthRed{}})
		if s.Gold() != 1 {
			t.Errorf("top power %d: Gold = %d, want 1 (clash win)", power, s.Gold())
		}
	}
}

// Tests that a tied clash (top power == 5) creates no Gold token.
func TestTestOfStrength_TieNoGold(t *testing.T) {
	s := sim.NewTurnState([]sim.Card{testutils.GenericAttack(0, 5)}, nil)
	(TestOfStrengthRed{}).Play(s, &sim.CardState{Card: TestOfStrengthRed{}})
	if s.Gold() != 0 {
		t.Errorf("top power 5: Gold = %d, want 0 (tie)", s.Gold())
	}
}

// Tests that a lost clash (top power ≤ 4) creates no Gold token (opponent's gold isn't
// modelled).
func TestTestOfStrength_LossNoGold(t *testing.T) {
	for _, power := range []int{0, 1, 2, 3, 4} {
		s := sim.NewTurnState([]sim.Card{testutils.GenericAttack(0, power)}, nil)
		(TestOfStrengthRed{}).Play(s, &sim.CardState{Card: TestOfStrengthRed{}})
		if s.Gold() != 0 {
			t.Errorf("top power %d: Gold = %d, want 0 (clash loss)", power, s.Gold())
		}
	}
}

// Tests that an empty deck makes the clash a no-op (no Gold created).
func TestTestOfStrength_EmptyDeckNoGold(t *testing.T) {
	s := sim.NewTurnState(nil, nil)
	(TestOfStrengthRed{}).Play(s, &sim.CardState{Card: TestOfStrengthRed{}})
	if s.Gold() != 0 {
		t.Errorf("empty deck: Gold = %d, want 0 (clash fails)", s.Gold())
	}
}
