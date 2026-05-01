package notimplemented

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
)

// TestCaptainsCall_NoAttackReturnsZero: no qualifying next attack card → +2 rider fizzles.
func TestCaptainsCall_NoAttackReturnsZero(t *testing.T) {
	s := sim.TurnState{}
	for _, c := range []sim.Card{CaptainsCallRed{}, CaptainsCallYellow{}, CaptainsCallBlue{}} {
		c.Play(&s, &sim.CardState{Card: c})
		if got := s.Value; got != 0 {
			t.Errorf("%s: Play() = %d, want 0", c.Name(), got)
		}
	}
}

// TestCaptainsCall_HighCostFilteredOut: a cost-3 attack is past Red's cost<=2 gate.
func TestCaptainsCall_HighCostFilteredOut(t *testing.T) {
	s := sim.TurnState{CardsRemaining: []*sim.CardState{{Card: testutils.GenericAttack(3, 0)}}}
	(CaptainsCallRed{}).Play(&s, &sim.CardState{Card: CaptainsCallRed{}})
	if got := s.Value; got != 0 {
		t.Errorf("Play() = %d, want 0 (cost 3 > 2)", got)
	}
}

// TestCaptainsCall_CostThresholdPerVariant: the +2 bonus is flat across variants, but each
// variant has its own cost threshold: Red cost<=2, Yellow cost<=1, Blue cost==0. A cost-2
// attack only triggers Red; a cost-1 attack triggers Red and Yellow; a cost-0 attack
// triggers all three. Each variant is exercised against a fresh target (rather than sharing
// one across the three calls) so a successful grant doesn't accumulate into the next
// variant's BonusAttack.
func TestCaptainsCall_CostThresholdPerVariant(t *testing.T) {
	cases := []struct {
		name   string
		cost   int
		red    int
		yellow int
		blue   int
	}{
		{"cost 2", 2, 2, 0, 0},
		{"cost 1", 1, 2, 2, 0},
		{"cost 0", 0, 2, 2, 2},
	}
	check := func(t *testing.T, label string, c sim.Card, cost, want int) {
		t.Helper()
		target := &sim.CardState{Card: testutils.GenericAttack(cost, 0)}
		s := sim.TurnState{CardsRemaining: []*sim.CardState{target}}
		c.Play(&s, &sim.CardState{Card: c})
		if got := s.Value; got != 0 {
			t.Errorf("%s: Play() = %d, want 0 (granter returns 0; +N rides on target's BonusAttack)", label, got)
		}
		if target.BonusAttack != want {
			t.Errorf("%s: target BonusAttack = %d, want %d", label, target.BonusAttack, want)
		}
	}
	for _, tc := range cases {
		check(t, tc.name+" Red", CaptainsCallRed{}, tc.cost, tc.red)
		check(t, tc.name+" Yellow", CaptainsCallYellow{}, tc.cost, tc.yellow)
		check(t, tc.name+" Blue", CaptainsCallBlue{}, tc.cost, tc.blue)
	}
}
