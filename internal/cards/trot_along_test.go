package cards

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
)

// TestTrotAlong_NoAttackReturnsZero covers the miss branch: no qualifying next attack → grant
// fizzles.
func TestTrotAlong_NoAttackReturnsZero(t *testing.T) {
	s := sim.TurnState{}
	(TrotAlongBlue{}).Play(&s, &sim.CardState{Card: TrotAlongBlue{}})
	if got := s.Value; got != 0 {
		t.Errorf("Play() = %d, want 0", got)
	}
}

// TestTrotAlong_HighPowerAttackDoesNotFire exercises the power<=3 filter: a power-4 attack in
// CardsRemaining is seen but doesn't pass the predicate, so the grant doesn't fire.
func TestTrotAlong_HighPowerAttackDoesNotFire(t *testing.T) {
	target := &sim.CardState{Card: testutils.GenericAttack(0, 4)}
	s := sim.TurnState{CardsRemaining: []*sim.CardState{target}}
	(TrotAlongBlue{}).Play(&s, &sim.CardState{Card: TrotAlongBlue{}})
	if got := s.Value; got != 0 {
		t.Errorf("Play() = %d, want 0 (power 4 > 3)", got)
	}
	if target.GrantedGoAgain {
		t.Error("target GrantedGoAgain = true, want false (power 4 > 3)")
	}
}

// TestTrotAlong_LowPowerAttackGrantsGoAgain exercises the hit branch: a power-3 attack qualifies.
func TestTrotAlong_LowPowerAttackGrantsGoAgain(t *testing.T) {
	target := &sim.CardState{Card: testutils.GenericAttack(0, 3)}
	s := sim.TurnState{CardsRemaining: []*sim.CardState{target}}
	(TrotAlongBlue{}).Play(&s, &sim.CardState{Card: TrotAlongBlue{}})
	if got := s.Value; got != 0 {
		t.Errorf("Play() = %d, want 0 (Trot Along grants go again, not damage)", got)
	}
	if !target.GrantedGoAgain {
		t.Error("target GrantedGoAgain = false, want true (power 3 ≤ 3)")
	}
}
