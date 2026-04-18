package generic

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

// TestWarmongersRecital_NoAttackReturnsZero: no qualifying next attack card → +3 rider fizzles.
func TestWarmongersRecital_NoAttackReturnsZero(t *testing.T) {
	s := card.TurnState{}
	for _, c := range []card.Card{WarmongersRecitalRed{}, WarmongersRecitalYellow{}, WarmongersRecitalBlue{}} {
		if got := c.Play(&s); got != 0 {
			t.Errorf("%s: Play() = %d, want 0", c.Name(), got)
		}
	}
}

// TestWarmongersRecital_NonAttackInRemainingFizzles: non-attack action fails the predicate.
func TestWarmongersRecital_NonAttackInRemainingFizzles(t *testing.T) {
	s := card.TurnState{CardsRemaining: []*card.PlayedCard{{Card: stubGenericAction()}}}
	if got := (WarmongersRecitalRed{}).Play(&s); got != 0 {
		t.Errorf("Play() = %d, want 0 (non-attack skipped)", got)
	}
}

// TestWarmongersRecital_NextAttackReturnsBonus: first attack-action in CardsRemaining triggers +3.
func TestWarmongersRecital_NextAttackReturnsBonus(t *testing.T) {
	s := card.TurnState{CardsRemaining: []*card.PlayedCard{{Card: stubGenericAttack(0, 0)}}}
	for _, c := range []card.Card{WarmongersRecitalRed{}, WarmongersRecitalYellow{}, WarmongersRecitalBlue{}} {
		if got := c.Play(&s); got != 3 {
			t.Errorf("%s: Play() = %d, want 3", c.Name(), got)
		}
	}
}
