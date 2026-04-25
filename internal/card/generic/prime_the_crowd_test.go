package generic

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

// TestPrimeTheCrowd_NoAttackReturnsZero: no qualifying next attack card → +4 rider fizzles.
func TestPrimeTheCrowd_NoAttackReturnsZero(t *testing.T) {
	s := card.TurnState{}
	for _, c := range []card.Card{PrimeTheCrowdRed{}, PrimeTheCrowdYellow{}, PrimeTheCrowdBlue{}} {
		if got := c.Play(&s, &card.CardState{}); got != 0 {
			t.Errorf("%s: Play() = %d, want 0", c.Name(), got)
		}
	}
}

// TestPrimeTheCrowd_NonAttackInRemainingFizzles: non-attack action fails the predicate.
func TestPrimeTheCrowd_NonAttackInRemainingFizzles(t *testing.T) {
	s := card.TurnState{CardsRemaining: []*card.CardState{{Card: stubGenericAction()}}}
	if got := (PrimeTheCrowdRed{}).Play(&s, &card.CardState{}); got != 0 {
		t.Errorf("Play() = %d, want 0 (non-attack skipped)", got)
	}
}

// TestPrimeTheCrowd_NextAttackReturnsBonus: first attack-action triggers the per-variant bonus
// (Red +4, Yellow +3, Blue +2).
func TestPrimeTheCrowd_NextAttackReturnsBonus(t *testing.T) {
	cases := []struct {
		c    card.Card
		want int
	}{
		{PrimeTheCrowdRed{}, 4},
		{PrimeTheCrowdYellow{}, 3},
		{PrimeTheCrowdBlue{}, 2},
	}
	for _, tc := range cases {
		target := &card.CardState{Card: stubGenericAttack(0, 0)}
		s := card.TurnState{CardsRemaining: []*card.CardState{target}}
		if got := tc.c.Play(&s, &card.CardState{}); got != 0 {
			t.Errorf("%s: Play() = %d, want 0 (granter returns 0; +N rides on target's BonusAttack)", tc.c.Name(), got)
		}
		if target.BonusAttack != tc.want {
			t.Errorf("%s: target BonusAttack = %d, want %d", tc.c.Name(), target.BonusAttack, tc.want)
		}
	}
}
