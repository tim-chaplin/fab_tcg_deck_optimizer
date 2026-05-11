package notimplemented

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
)

// TestRestvineElixir_NoAttackReturnsZero: no qualifying next attack card → +3 rider fizzles.
func TestRestvineElixir_NoAttackReturnsZero(t *testing.T) {
	s := sim.TurnState{}
	sim.ResolveChainStep(&s, s.Logger(), &card.CardState{Card: RestvineElixirRed{}})
	if got := s.Value(); got != 0 {
		t.Errorf("Play() = %d, want 0", got)
	}
}

// TestRestvineElixir_NonAttackInRemainingFizzles: non-attack action fails the predicate.
func TestRestvineElixir_NonAttackInRemainingFizzles(t *testing.T) {
	s := sim.NewTurnStateFromSpec(sim.TurnStateSpec{CardsRemaining: []*card.CardState{{Card: testutils.GenericAction()}}})
	sim.ResolveChainStep(&s, s.Logger(), &card.CardState{Card: RestvineElixirRed{}})
	if got := s.Value(); got != 0 {
		t.Errorf("Play() = %d, want 0 (non-attack skipped)", got)
	}
}

// TestRestvineElixir_NextAttackGrantsBonusAttack: first attack-action picks up +3 on its
// BonusAttack. Granter returns 0; the +3 attributes to the target.
func TestRestvineElixir_NextAttackGrantsBonusAttack(t *testing.T) {
	target := &card.CardState{Card: testutils.GenericAttack(0, 0)}
	s := sim.NewTurnStateFromSpec(sim.TurnStateSpec{CardsRemaining: []*card.CardState{target}})
	sim.ResolveChainStep(&s, s.Logger(), &card.CardState{Card: RestvineElixirRed{}})
	if got := s.Value(); got != 0 {
		t.Errorf("Play() = %d, want 0 (granter returns 0; +N rides on target's BonusAttack)", got)
	}
	if target.BonusAttack != 3 {
		t.Errorf("target BonusAttack = %d, want 3", target.BonusAttack)
	}
}
