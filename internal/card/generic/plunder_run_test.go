package generic

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

// TestPlunderRun_NoAttackReturnsZero: no qualifying next attack card → +3 rider fizzles.
func TestPlunderRun_NoAttackReturnsZero(t *testing.T) {
	s := card.TurnState{}
	for _, c := range []card.Card{PlunderRunRed{}, PlunderRunYellow{}, PlunderRunBlue{}} {
		if got := c.Play(&s); got != 0 {
			t.Errorf("%s: Play() = %d, want 0", c.Name(), got)
		}
	}
}

// TestPlunderRun_NonAttackInRemainingFizzles: non-attack action fails the predicate.
func TestPlunderRun_NonAttackInRemainingFizzles(t *testing.T) {
	s := card.TurnState{CardsRemaining: []*card.PlayedCard{{Card: stubGenericAction()}}}
	if got := (PlunderRunRed{}).Play(&s); got != 0 {
		t.Errorf("Play() = %d, want 0 (non-attack skipped)", got)
	}
}

// TestPlunderRun_NextAttackReturnsBonus: arsenal-played copy with a queued attack action
// triggers the per-variant bonus (Red +3, Yellow +2, Blue +1).
func TestPlunderRun_NextAttackReturnsBonus(t *testing.T) {
	cases := []struct {
		c    card.Card
		want int
	}{
		{PlunderRunRed{}, 3},
		{PlunderRunYellow{}, 2},
		{PlunderRunBlue{}, 1},
	}
	for _, tc := range cases {
		s := card.TurnState{
			SelfFromArsenal: true,
			CardsRemaining:  []*card.PlayedCard{{Card: stubGenericAttack(0, 0)}},
		}
		if got := tc.c.Play(&s); got != tc.want {
			t.Errorf("%s: Play() = %d, want %d", tc.c.Name(), got, tc.want)
		}
	}
}

// TestPlunderRun_HandPlayedFizzles: hand-played copy fails the from-arsenal gate even when a
// queued attack action would otherwise satisfy the rider.
func TestPlunderRun_HandPlayedFizzles(t *testing.T) {
	s := card.TurnState{
		CardsRemaining: []*card.PlayedCard{{Card: stubGenericAttack(0, 0)}},
	}
	if got := (PlunderRunRed{}).Play(&s); got != 0 {
		t.Errorf("Play() = %d, want 0 (hand-played, not from arsenal)", got)
	}
}
