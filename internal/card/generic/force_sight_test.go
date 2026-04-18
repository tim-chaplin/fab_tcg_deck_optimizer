package generic

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

// TestForceSight_NoAttackReturnsZero: no qualifying next attack card → +3 rider fizzles.
func TestForceSight_NoAttackReturnsZero(t *testing.T) {
	s := card.TurnState{}
	for _, c := range []card.Card{ForceSightRed{}, ForceSightYellow{}, ForceSightBlue{}} {
		if got := c.Play(&s); got != 0 {
			t.Errorf("%s: Play() = %d, want 0", c.Name(), got)
		}
	}
}

// TestForceSight_NonAttackInRemainingFizzles: non-attack action fails the predicate.
func TestForceSight_NonAttackInRemainingFizzles(t *testing.T) {
	s := card.TurnState{CardsRemaining: []*card.PlayedCard{{Card: stubGenericAction()}}}
	if got := (ForceSightRed{}).Play(&s); got != 0 {
		t.Errorf("Play() = %d, want 0 (non-attack skipped)", got)
	}
}

// TestForceSight_NextAttackReturnsBonus: first attack-action triggers the per-variant bonus
// (Red +3, Yellow +2, Blue +1).
func TestForceSight_NextAttackReturnsBonus(t *testing.T) {
	s := card.TurnState{CardsRemaining: []*card.PlayedCard{{Card: stubGenericAttack(0, 0)}}}
	cases := []struct {
		c    card.Card
		want int
	}{
		{ForceSightRed{}, 3},
		{ForceSightYellow{}, 2},
		{ForceSightBlue{}, 1},
	}
	for _, tc := range cases {
		if got := tc.c.Play(&s); got != tc.want {
			t.Errorf("%s: Play() = %d, want %d", tc.c.Name(), got, tc.want)
		}
	}
}
