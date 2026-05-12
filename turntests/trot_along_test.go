package turntests

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/cards"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/gameengine"
)

// TestTrotAlong_NoAttackReturnsZero covers the miss branch: no qualifying next attack → grant
// fizzles.
func TestTrotAlong_NoAttackReturnsZero(t *testing.T) {
	s := gameengine.New()
	s.ResolveChainStep(s.Logger(), &card.CardState{Card: cards.TrotAlongBlue{}})
	if got := s.Value(); got != 0 {
		t.Errorf("Play() = %d, want 0", got)
	}
}

// TestTrotAlong_HighPowerAttackDoesNotFire exercises the power<=3 filter: a power-4 attack in
// CardsRemaining is seen but doesn't pass the predicate, so the grant doesn't fire.
func TestTrotAlong_HighPowerAttackDoesNotFire(t *testing.T) {
	target := &card.CardState{Card: testutils.GenericAttack(0, 4)}
	s := gameengine.NewFromSpec(gameengine.Spec{CardsRemaining: []*card.CardState{target}})
	s.ResolveChainStep(s.Logger(), &card.CardState{Card: cards.TrotAlongBlue{}})
	if got := s.Value(); got != 0 {
		t.Errorf("Play() = %d, want 0 (power 4 > 3)", got)
	}
	if target.GrantedGoAgain {
		t.Error("target GrantedGoAgain = true, want false (power 4 > 3)")
	}
}

// TestTrotAlong_LowPowerAttackGrantsGoAgain exercises the hit branch: a power-3 attack qualifies.
func TestTrotAlong_LowPowerAttackGrantsGoAgain(t *testing.T) {
	target := &card.CardState{Card: testutils.GenericAttack(0, 3)}
	s := gameengine.NewFromSpec(gameengine.Spec{CardsRemaining: []*card.CardState{target}})
	s.ResolveChainStep(s.Logger(), &card.CardState{Card: cards.TrotAlongBlue{}})
	if got := s.Value(); got != 0 {
		t.Errorf("Play() = %d, want 0 (Trot Along grants go again, not damage)", got)
	}
	if !target.GrantedGoAgain {
		t.Error("target GrantedGoAgain = false, want true (power 3 ≤ 3)")
	}
}

// TestTrotAlong_GrantsGoAgainToWeaponSwing pins the "your next attack" wording: a weapon
// swing (TypeWeapon, no TypeAction) with base power ≤ 3 qualifies. RunebladeWeapon's
// Attack() is 0 so the power gate trivially passes.
func TestTrotAlong_GrantsGoAgainToWeaponSwing(t *testing.T) {
	target := &card.CardState{Card: testutils.RunebladeWeapon{}}
	s := gameengine.NewFromSpec(gameengine.Spec{CardsRemaining: []*card.CardState{target}})
	s.ResolveChainStep(s.Logger(), &card.CardState{Card: cards.TrotAlongBlue{}})
	if !target.GrantedGoAgain {
		t.Error("weapon swing should get go again ('your next attack' has no 'action card' qualifier)")
	}
}
