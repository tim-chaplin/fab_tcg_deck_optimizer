package notimplemented

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
)

// TestSapwoodElixir_NoAttackReturnsZero: no qualifying next attack card → +3 rider fizzles.
func TestSapwoodElixir_NoAttackReturnsZero(t *testing.T) {
	s := sim.TurnState{}
	(SapwoodElixirRed{}).Play(&s, &sim.CardState{Card: SapwoodElixirRed{}})
	if got := s.Value; got != 0 {
		t.Errorf("Play() = %d, want 0", got)
	}
}

// TestSapwoodElixir_NonAttackInRemainingFizzles: non-attack action fails the predicate.
func TestSapwoodElixir_NonAttackInRemainingFizzles(t *testing.T) {
	s := sim.TurnState{CardsRemaining: []*sim.CardState{{Card: testutils.GenericAction()}}}
	(SapwoodElixirRed{}).Play(&s, &sim.CardState{Card: SapwoodElixirRed{}})
	if got := s.Value; got != 0 {
		t.Errorf("Play() = %d, want 0 (non-attack skipped)", got)
	}
}

// TestSapwoodElixir_NextAttackGrantsBonusAttack: first attack-action picks up +3 on its
// BonusAttack. Granter returns 0; the +3 attributes to the target.
func TestSapwoodElixir_NextAttackGrantsBonusAttack(t *testing.T) {
	target := &sim.CardState{Card: testutils.GenericAttack(0, 0)}
	s := sim.TurnState{CardsRemaining: []*sim.CardState{target}}
	(SapwoodElixirRed{}).Play(&s, &sim.CardState{Card: SapwoodElixirRed{}})
	if got := s.Value; got != 0 {
		t.Errorf("Play() = %d, want 0 (granter returns 0; +N rides on target's BonusAttack)", got)
	}
	if target.BonusAttack != 3 {
		t.Errorf("target BonusAttack = %d, want 3", target.BonusAttack)
	}
}
