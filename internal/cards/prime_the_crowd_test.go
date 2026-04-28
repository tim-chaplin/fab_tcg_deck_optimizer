package cards

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
)

// TestPrimeTheCrowd_NoAttackReturnsZero: no qualifying next attack card → +4 rider fizzles.
func TestPrimeTheCrowd_NoAttackReturnsZero(t *testing.T) {
	s := sim.TurnState{}
	for _, c := range []sim.Card{PrimeTheCrowdRed{}, PrimeTheCrowdYellow{}, PrimeTheCrowdBlue{}} {
		c.Play(&s, &sim.CardState{Card: c})
		if got := s.Value; got != 0 {
			t.Errorf("%s: Play() = %d, want 0", c.Name(), got)
		}
	}
}

// TestPrimeTheCrowd_NonAttackInRemainingFizzles: non-attack action fails the predicate.
func TestPrimeTheCrowd_NonAttackInRemainingFizzles(t *testing.T) {
	s := sim.TurnState{CardsRemaining: []*sim.CardState{{Card: testutils.GenericAction()}}}
	(PrimeTheCrowdRed{}).Play(&s, &sim.CardState{Card: PrimeTheCrowdRed{}})
	if got := s.Value; got != 0 {
		t.Errorf("Play() = %d, want 0 (non-attack skipped)", got)
	}
}

// TestPrimeTheCrowd_NextAttackReturnsBonus: first attack-action triggers the per-variant bonus
// (Red +4, Yellow +3, Blue +2).
func TestPrimeTheCrowd_NextAttackReturnsBonus(t *testing.T) {
	cases := []struct {
		c    sim.Card
		want int
	}{
		{PrimeTheCrowdRed{}, 4},
		{PrimeTheCrowdYellow{}, 3},
		{PrimeTheCrowdBlue{}, 2},
	}
	for _, tc := range cases {
		target := &sim.CardState{Card: testutils.GenericAttack(0, 0)}
		s := sim.TurnState{CardsRemaining: []*sim.CardState{target}}
		tc.c.Play(&s, &sim.CardState{Card: tc.c})
		if got := s.Value; got != 0 {
			t.Errorf("%s: Play() = %d, want 0 (granter returns 0; +N rides on target's BonusAttack)", tc.c.Name(), got)
		}
		if target.BonusAttack != tc.want {
			t.Errorf("%s: target BonusAttack = %d, want %d", tc.c.Name(), target.BonusAttack, tc.want)
		}
	}
}
