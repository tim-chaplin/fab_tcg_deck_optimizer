package cards

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
)

// Tests that Test of Strength creates a Gold token when the deck-top wins the clash
// (top power ≥ 6 vs opponent's modelled 5).
func TestTestOfStrength_WinCreatesGold(t *testing.T) {
	for _, power := range []int{6, 7} {
		s := sim.NewTurnStateFromCards([]card.Card{testutils.GenericAttack(0, power)}, nil)
		sim.ResolveChainStep(s, s.Logger(), &card.CardState{Card: TestOfStrengthRed{}})
		if s.Gold() != 1 {
			t.Errorf("top power %d: Gold = %d, want 1 (clash win)", power, s.Gold())
		}
	}
}

// Tests that a tied clash (top power == 5) creates no Gold token.
func TestTestOfStrength_TieNoGold(t *testing.T) {
	s := sim.NewTurnStateFromCards([]card.Card{testutils.GenericAttack(0, 5)}, nil)
	sim.ResolveChainStep(s, s.Logger(), &card.CardState{Card: TestOfStrengthRed{}})
	if s.Gold() != 0 {
		t.Errorf("top power 5: Gold = %d, want 0 (tie)", s.Gold())
	}
}

// Tests that a lost clash (top power ≤ 4) creates no Gold token and subtracts one Value
// to reflect the opponent's Gold token.
func TestTestOfStrength_LossNoGoldAndDocksValue(t *testing.T) {
	for _, power := range []int{0, 1, 2, 3, 4} {
		s := sim.NewTurnStateFromCards([]card.Card{testutils.GenericAttack(0, power)}, nil)
		valueBefore := s.Value()
		sim.ResolveChainStep(s, s.Logger(), &card.CardState{Card: TestOfStrengthRed{}})
		if s.Gold() != 0 {
			t.Errorf("top power %d: Gold = %d, want 0 (clash loss)", power, s.Gold())
		}
		if s.Value()-valueBefore != -1 {
			t.Errorf("top power %d: Value delta = %d, want -1 (Clash loss costs 1)",
				power, s.Value()-valueBefore)
		}
	}
}

// Tests that an empty deck makes the clash a no-op (no Gold created).
func TestTestOfStrength_EmptyDeckNoGold(t *testing.T) {
	s := sim.NewTurnStateFromCards(nil, nil)
	sim.ResolveChainStep(s, s.Logger(), &card.CardState{Card: TestOfStrengthRed{}})
	if s.Gold() != 0 {
		t.Errorf("empty deck: Gold = %d, want 0 (clash fails)", s.Gold())
	}
}
