package generic

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

// TestWarmongersRecital_NoAttackReturnsZero: no qualifying next attack card → +3 rider fizzles.
func TestWarmongersRecital_NoAttackReturnsZero(t *testing.T) {
	s := card.TurnState{}
	for _, c := range []card.Card{WarmongersRecitalRed{}, WarmongersRecitalYellow{}, WarmongersRecitalBlue{}} {
		c.Play(&s, &card.CardState{Card: c})
		if got := s.Value; got != 0 {
			t.Errorf("%s: Play() = %d, want 0", c.Name(), got)
		}
	}
}

// TestWarmongersRecital_NonAttackInRemainingFizzles: non-attack action fails the predicate.
func TestWarmongersRecital_NonAttackInRemainingFizzles(t *testing.T) {
	s := card.TurnState{CardsRemaining: []*card.CardState{{Card: stubGenericAction()}}}
	(WarmongersRecitalRed{}).Play(&s, &card.CardState{Card: WarmongersRecitalRed{}})
	if got := s.Value; got != 0 {
		t.Errorf("Play() = %d, want 0 (non-attack skipped)", got)
	}
}

// TestWarmongersRecital_NextAttackReturnsBonus: first attack-action in CardsRemaining triggers
// the per-variant bonus (Red +3, Yellow +2, Blue +1).
func TestWarmongersRecital_NextAttackReturnsBonus(t *testing.T) {
	cases := []struct {
		c    card.Card
		want int
	}{
		{WarmongersRecitalRed{}, 3},
		{WarmongersRecitalYellow{}, 2},
		{WarmongersRecitalBlue{}, 1},
	}
	for _, tc := range cases {
		target := &card.CardState{Card: stubGenericAttack(0, 0)}
		s := card.TurnState{CardsRemaining: []*card.CardState{target}}
		tc.c.Play(&s, &card.CardState{Card: tc.c})
		if got := s.Value; got != 0 {
			t.Errorf("%s: Play() = %d, want 0 (granter returns 0; +N rides on target's BonusAttack)", tc.c.Name(), got)
		}
		if target.BonusAttack != tc.want {
			t.Errorf("%s: target BonusAttack = %d, want %d", tc.c.Name(), target.BonusAttack, tc.want)
		}
	}
}
