package generic

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

// TestMinnowism_NoAttackReturnsZero: no qualifying next attack card → +3 rider fizzles.
func TestMinnowism_NoAttackReturnsZero(t *testing.T) {
	s := card.TurnState{}
	for _, c := range []card.Card{MinnowismRed{}, MinnowismYellow{}, MinnowismBlue{}} {
		c.Play(&s, &card.CardState{Card: c})
		if got := s.Value; got != 0{
			t.Errorf("%s: Play() = %d, want 0", c.Name(), got)
		}
	}
}

// TestMinnowism_HighPowerFilteredOut: a power-4 attack is seen but the power<=3 filter rejects it,
// so the rider fizzles without falling through to a later match.
func TestMinnowism_HighPowerFilteredOut(t *testing.T) {
	s := card.TurnState{CardsRemaining: []*card.CardState{{Card: stubGenericAttack(0, 4)}}}
	(MinnowismRed{}).Play(&s, &card.CardState{Card: MinnowismRed{}})
	if got := s.Value; got != 0{
		t.Errorf("Play() = %d, want 0 (power 4 > 3)", got)
	}
}

// TestMinnowism_LowPowerReturnsBonus: first power<=3 attack triggers the per-variant bonus
// (Red +3, Yellow +2, Blue +1).
func TestMinnowism_LowPowerReturnsBonus(t *testing.T) {
	cases := []struct {
		c    card.Card
		want int
	}{
		{MinnowismRed{}, 3},
		{MinnowismYellow{}, 2},
		{MinnowismBlue{}, 1},
	}
	for _, tc := range cases {
		target := &card.CardState{Card: stubGenericAttack(0, 3)}
		s := card.TurnState{CardsRemaining: []*card.CardState{target}}
		tc.c.Play(&s, &card.CardState{Card: tc.c})
		if got := s.Value; got != 0{
			t.Errorf("%s: Play() = %d, want 0 (granter returns 0; +N rides on target's BonusAttack)", tc.c.Name(), got)
		}
		if target.BonusAttack != tc.want {
			t.Errorf("%s: target BonusAttack = %d, want %d", tc.c.Name(), target.BonusAttack, tc.want)
		}
	}
}
