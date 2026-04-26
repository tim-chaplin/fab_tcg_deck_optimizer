package generic

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

// TestTestOfStrength_WinCreditsGoldToken: a top-of-deck attack of 6 or 7 wins the clash; Play
// returns +card.GoldTokenValue.
func TestTestOfStrength_WinCreditsGoldToken(t *testing.T) {
	for _, power := range []int{6, 7} {
		s := &card.TurnState{Deck: []card.Card{stubGenericAttack(0, power)}}
		(TestOfStrengthRed{}).Play(s, &card.CardState{Card: TestOfStrengthRed{}})
		if got := s.Value; got != card.GoldTokenValue{
			t.Errorf("top power %d: Play() = %d, want %d", power, got, card.GoldTokenValue)
		}
	}
}

// TestTestOfStrength_TieCreditsZero: a top-of-deck attack of exactly 5 ties the clash; nobody
// gets the Gold token, so Play returns 0.
func TestTestOfStrength_TieCreditsZero(t *testing.T) {
	s := &card.TurnState{Deck: []card.Card{stubGenericAttack(0, 5)}}
	(TestOfStrengthRed{}).Play(s, &card.CardState{Card: TestOfStrengthRed{}})
	if got := s.Value; got != 0{
		t.Errorf("top power 5: Play() = %d, want 0 (tie)", got)
	}
}

// TestTestOfStrength_LossSubtractsGoldToken: a top-of-deck attack of 4 or below loses the clash;
// the opponent creates the Gold token, so Play returns -card.GoldTokenValue.
func TestTestOfStrength_LossSubtractsGoldToken(t *testing.T) {
	for _, power := range []int{0, 1, 2, 3, 4} {
		s := &card.TurnState{Deck: []card.Card{stubGenericAttack(0, power)}}
		(TestOfStrengthRed{}).Play(s, &card.CardState{Card: TestOfStrengthRed{}})
		if got := s.Value; got != -card.GoldTokenValue{
			t.Errorf("top power %d: Play() = %d, want %d", power, got, -card.GoldTokenValue)
		}
	}
}

// TestTestOfStrength_EmptyDeckReturnsZero: with no card to reveal the clash effect fails (per
// comprehensive rules 8.5.45d); Play returns 0.
func TestTestOfStrength_EmptyDeckReturnsZero(t *testing.T) {
	s := &card.TurnState{}
	(TestOfStrengthRed{}).Play(s, &card.CardState{Card: TestOfStrengthRed{}})
	if got := s.Value; got != 0{
		t.Errorf("empty deck: Play() = %d, want 0 (clash fails)", got)
	}
}
