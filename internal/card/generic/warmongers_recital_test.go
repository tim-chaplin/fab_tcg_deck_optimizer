package generic

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

// TestWarmongersRecital_NoAttackReturnsZero: no qualifying next attack card → +3 rider fizzles.
func TestWarmongersRecital_NoAttackReturnsZero(t *testing.T) {
	s := card.TurnState{}
	for _, c := range []card.Card{WarmongersRecitalRed{}, WarmongersRecitalYellow{}, WarmongersRecitalBlue{}} {
		if got := c.Play(&s, nil); got != 0 {
			t.Errorf("%s: Play() = %d, want 0", c.Name(), got)
		}
	}
}

// TestWarmongersRecital_NonAttackInRemainingFizzles: non-attack action fails the predicate.
func TestWarmongersRecital_NonAttackInRemainingFizzles(t *testing.T) {
	s := card.TurnState{CardsRemaining: []*card.PlayedCard{{Card: stubGenericAction()}}}
	if got := (WarmongersRecitalRed{}).Play(&s, nil); got != 0 {
		t.Errorf("Play() = %d, want 0 (non-attack skipped)", got)
	}
}

// TestWarmongersRecital_NextAttackReturnsBonus: first attack-action in CardsRemaining triggers
// the per-variant bonus (Red +3, Yellow +2, Blue +1).
func TestWarmongersRecital_NextAttackReturnsBonus(t *testing.T) {
	s := card.TurnState{CardsRemaining: []*card.PlayedCard{{Card: stubGenericAttack(0, 0)}}}
	cases := []struct {
		c    card.Card
		want int
	}{
		{WarmongersRecitalRed{}, 3},
		{WarmongersRecitalYellow{}, 2},
		{WarmongersRecitalBlue{}, 1},
	}
	for _, tc := range cases {
		if got := tc.c.Play(&s, nil); got != tc.want {
			t.Errorf("%s: Play() = %d, want %d", tc.c.Name(), got, tc.want)
		}
	}
}
