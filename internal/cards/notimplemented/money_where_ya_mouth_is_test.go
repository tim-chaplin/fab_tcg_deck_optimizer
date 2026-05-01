package notimplemented

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
)

// TestMoneyWhereYaMouthIs_NoAttackReturnsZero: no qualifying next attack card → +3 fizzles.
func TestMoneyWhereYaMouthIs_NoAttackReturnsZero(t *testing.T) {
	s := sim.TurnState{}
	for _, c := range []sim.Card{MoneyWhereYaMouthIsRed{}, MoneyWhereYaMouthIsYellow{}, MoneyWhereYaMouthIsBlue{}} {
		c.Play(&s, &sim.CardState{Card: c})
		if got := s.Value; got != 0 {
			t.Errorf("%s: Play() = %d, want 0", c.Name(), got)
		}
	}
}

// TestMoneyWhereYaMouthIs_NonAttackInRemainingFizzles: non-attack action fails the predicate.
func TestMoneyWhereYaMouthIs_NonAttackInRemainingFizzles(t *testing.T) {
	s := sim.TurnState{CardsRemaining: []*sim.CardState{{Card: testutils.GenericAction()}}}
	(MoneyWhereYaMouthIsRed{}).Play(&s, &sim.CardState{Card: MoneyWhereYaMouthIsRed{}})
	if got := s.Value; got != 0 {
		t.Errorf("Play() = %d, want 0 (non-attack skipped)", got)
	}
}

// TestMoneyWhereYaMouthIs_NextAttackReturnsBonus: first attack-action triggers the per-variant
// bonus (Red +3, Yellow +2, Blue +1).
func TestMoneyWhereYaMouthIs_NextAttackReturnsBonus(t *testing.T) {
	cases := []struct {
		c    sim.Card
		want int
	}{
		{MoneyWhereYaMouthIsRed{}, 3},
		{MoneyWhereYaMouthIsYellow{}, 2},
		{MoneyWhereYaMouthIsBlue{}, 1},
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
