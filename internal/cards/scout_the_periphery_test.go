package cards

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
)

// TestScoutThePeriphery_NoAttackReturnsZero: no qualifying next attack card → +3 rider fizzles.
func TestScoutThePeriphery_NoAttackReturnsZero(t *testing.T) {
	s := sim.TurnState{}
	for _, c := range []sim.Card{ScoutThePeripheryRed{}, ScoutThePeripheryYellow{}, ScoutThePeripheryBlue{}} {
		c.Play(&s, &sim.CardState{Card: c})
		if got := s.Value; got != 0 {
			t.Errorf("%s: Play() = %d, want 0", c.Name(), got)
		}
	}
}

// TestScoutThePeriphery_NonAttackInRemainingFizzles: non-attack action (even from arsenal)
// fails the predicate — only attack actions count as the rider's target.
func TestScoutThePeriphery_NonAttackInRemainingFizzles(t *testing.T) {
	s := sim.TurnState{CardsRemaining: []*sim.CardState{{Card: testutils.GenericAction(), FromArsenal: true}}}
	(ScoutThePeripheryRed{}).Play(&s, &sim.CardState{Card: ScoutThePeripheryRed{}})
	if got := s.Value; got != 0 {
		t.Errorf("Play() = %d, want 0 (non-attack skipped)", got)
	}
}

// TestScoutThePeriphery_HandPlayedAttackFizzles: queued attack action that wasn't played from
// arsenal fails the rider's "next attack action card you play from arsenal" target gate.
func TestScoutThePeriphery_HandPlayedAttackFizzles(t *testing.T) {
	s := sim.TurnState{CardsRemaining: []*sim.CardState{{Card: testutils.GenericAttack(0, 0)}}}
	(ScoutThePeripheryRed{}).Play(&s, &sim.CardState{Card: ScoutThePeripheryRed{}})
	if got := s.Value; got != 0 {
		t.Errorf("Play() = %d, want 0 (target attack not from arsenal)", got)
	}
}

// TestScoutThePeriphery_NextArsenalAttackReturnsBonus: when the queued attack action is
// itself played from arsenal the per-variant bonus lands on that card's BonusAttack so
// the buffed attack's EffectiveAttack picks up the +N (Red +3, Yellow +2, Blue +1). The
// granter itself credits 0 — the bonus rides on the target.
func TestScoutThePeriphery_NextArsenalAttackReturnsBonus(t *testing.T) {
	cases := []struct {
		c    sim.Card
		want int
	}{
		{ScoutThePeripheryRed{}, 3},
		{ScoutThePeripheryYellow{}, 2},
		{ScoutThePeripheryBlue{}, 1},
	}
	for _, tc := range cases {
		target := &sim.CardState{Card: testutils.GenericAttack(0, 0), FromArsenal: true}
		s := sim.TurnState{CardsRemaining: []*sim.CardState{target}}
		tc.c.Play(&s, &sim.CardState{Card: tc.c})
		if got := s.Value; got != 0 {
			t.Errorf("%s: granter credits %d, want 0 (bonus rides on target)", tc.c.Name(), got)
		}
		if got := target.BonusAttack; got != tc.want {
			t.Errorf("%s: target.BonusAttack = %d, want %d", tc.c.Name(), got, tc.want)
		}
	}
}
