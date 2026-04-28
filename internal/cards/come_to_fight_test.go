package cards

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
)

// TestComeToFight_NoAttackReturnsZero: no qualifying next attack card → +3 rider fizzles.
func TestComeToFight_NoAttackReturnsZero(t *testing.T) {
	s := sim.TurnState{}
	for _, c := range []sim.Card{ComeToFightRed{}, ComeToFightYellow{}, ComeToFightBlue{}} {
		c.Play(&s, &sim.CardState{Card: c})
		if got := s.Value; got != 0 {
			t.Errorf("%s: Play() = %d, want 0", c.Name(), got)
		}
	}
}

// TestComeToFight_NonAttackInRemainingFizzles: only an action (no attack) in CardsRemaining — the
// attack-action predicate rejects it.
func TestComeToFight_NonAttackInRemainingFizzles(t *testing.T) {
	s := sim.TurnState{CardsRemaining: []*sim.CardState{{Card: testutils.GenericAction()}}}
	(ComeToFightRed{}).Play(&s, &sim.CardState{Card: ComeToFightRed{}})
	if got := s.Value; got != 0 {
		t.Errorf("Play() = %d, want 0 (non-attack skipped)", got)
	}
}

// TestComeToFight_NextAttackReturnsBonus: first attack-action in CardsRemaining triggers the
// per-variant bonus (Red +3, Yellow +2, Blue +1).
func TestComeToFight_NextAttackReturnsBonus(t *testing.T) {
	cases := []struct {
		c    sim.Card
		want int
	}{
		{ComeToFightRed{}, 3},
		{ComeToFightYellow{}, 2},
		{ComeToFightBlue{}, 1},
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
