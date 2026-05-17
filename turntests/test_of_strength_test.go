package turntests

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card/cards"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/gameengine"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
)

// Tests that Test of Strength creates a Gold token when the deck-top wins the clash
// (top power ≥ 6 vs opponent's modelled 5).
func TestTestOfStrength_WinCreatesGold(t *testing.T) {
	for _, power := range []int{6, 7} {
		ge := &gameengine.GameEngine{GameState: gameengine.GameStateBuilder().SetCards([]card.Card{testutils.GenericAttack(0, power)}).Build()}
		ge.ResolveChainStep(ge.Logger(), &card.CardState{Card: cards.TestOfStrengthRed{}})
		if ge.GoldCount() != 1 {
			t.Errorf("top power %d: Gold = %d, want 1 (clash win)", power, ge.GoldCount())
		}
	}
}

// Tests that a tied clash (top power == 5) creates no Gold token.
func TestTestOfStrength_TieNoGold(t *testing.T) {
	ge := &gameengine.GameEngine{GameState: gameengine.GameStateBuilder().SetCards([]card.Card{testutils.GenericAttack(0, 5)}).Build()}
	ge.ResolveChainStep(ge.Logger(), &card.CardState{Card: cards.TestOfStrengthRed{}})
	if ge.GoldCount() != 0 {
		t.Errorf("top power 5: Gold = %d, want 0 (tie)", ge.GoldCount())
	}
}

// Tests that a lost clash (top power ≤ 4) creates no Gold token and subtracts one Value
// to reflect the opponent's Gold token.
func TestTestOfStrength_LossNoGoldAndDocksValue(t *testing.T) {
	for _, power := range []int{0, 1, 2, 3, 4} {
		ge := &gameengine.GameEngine{GameState: gameengine.GameStateBuilder().SetCards([]card.Card{testutils.GenericAttack(0, power)}).Build()}
		valueBefore := ge.Value()
		ge.ResolveChainStep(ge.Logger(), &card.CardState{Card: cards.TestOfStrengthRed{}})
		if ge.GoldCount() != 0 {
			t.Errorf("top power %d: Gold = %d, want 0 (clash loss)", power, ge.GoldCount())
		}
		if ge.Value()-valueBefore != -1 {
			t.Errorf("top power %d: Value delta = %d, want -1 (Clash loss costs 1)",
				power, ge.Value()-valueBefore)
		}
	}
}

// Tests that an empty deck makes the clash a no-op (no Gold created).
func TestTestOfStrength_EmptyDeckNoGold(t *testing.T) {
	ge := gameengine.New()
	ge.ResolveChainStep(ge.Logger(), &card.CardState{Card: cards.TestOfStrengthRed{}})
	if ge.GoldCount() != 0 {
		t.Errorf("empty deck: Gold = %d, want 0 (clash fails)", ge.GoldCount())
	}
}
