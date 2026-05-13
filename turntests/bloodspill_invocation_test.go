package turntests

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/cards"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/gameengine"
)

// TestBloodspillInvocation_BlockCoversIncomingReturnsN: the aura survives the opponent's turn
// when we block all incoming damage, and pays N on a future-turn pop.
func TestBloodspillInvocation_BlockCoversIncomingReturnsN(t *testing.T) {
	cases := []struct {
		c card.Card
		n int
	}{
		{cards.BloodspillInvocationRed{}, 3},
		{cards.BloodspillInvocationYellow{}, 2},
		{cards.BloodspillInvocationBlue{}, 1},
	}
	for _, tc := range cases {
		s := &gameengine.GameEngine{GameState: gameengine.GameStateBuilder().
			SetIncomingDamage(3).
			SetBlockTotal(3).
			Build()}
		s.ResolveChainStep(s.Logger(), &card.CardState{Card: tc.c})
		if got := s.Value(); got != tc.n {
			t.Errorf("%s: Play() = %d, want %d (block == incoming)", tc.c.Name(), got, tc.n)
		}
	}
}

// TestBloodspillInvocation_BlockShortReturnsZero: if we take damage and have no same-turn
// attack action likely to hit, Bloodspill dies without creating Runechants.
func TestBloodspillInvocation_BlockShortReturnsZero(t *testing.T) {
	cases := []card.Card{
		cards.BloodspillInvocationRed{},
		cards.BloodspillInvocationYellow{},
		cards.BloodspillInvocationBlue{},
	}
	for _, c := range cases {
		s := &gameengine.GameEngine{GameState: gameengine.GameStateBuilder().
			SetIncomingDamage(3).
			SetBlockTotal(2).
			Build()}
		s.ResolveChainStep(s.Logger(), &card.CardState{Card: c})
		if got := s.Value(); got != 0 {
			t.Errorf("%s: Play() = %d, want 0 (block < incoming, no same-turn pop)", c.Name(), got)
		}
	}
}

// TestBloodspillInvocation_SameTurnPopBySalientAttackAction: a later attack action with a
// likely-to-hit power pops Bloodspill this turn for its full N — even if we're taking damage.
func TestBloodspillInvocation_SameTurnPopBySalientAttackAction(t *testing.T) {
	s := &gameengine.GameEngine{GameState: gameengine.GameStateBuilder().
		SetIncomingDamage(3).
		SetBlockTotal(0).
		SetCardsRemaining([]*card.CardState{{Card: testutils.AttackWithPower{Power: 4}}}).
		Build()}
	s.ResolveChainStep(s.Logger(), &card.CardState{Card: cards.BloodspillInvocationRed{}})
	if got := s.Value(); got != 3 {
		t.Errorf("Play() = %d, want 3 (Attack=4 attack action pops Bloodspill same turn)", got)
	}
}

// TestBloodspillInvocation_WeaponDoesNotPop: Bloodspill's rider is gated on attack ACTION cards
// hitting — a weapon swing that hits doesn't trigger its destruction, even with a Runechant
// firing alongside.
func TestBloodspillInvocation_WeaponDoesNotPop(t *testing.T) {
	s := &gameengine.GameEngine{GameState: gameengine.GameStateBuilder().
		SetIncomingDamage(3).
		SetBlockTotal(0).
		SetCardsRemaining([]*card.CardState{{Card: testutils.RunebladeWeapon{}}}).
		Build()}
	s.CreateAura(sim.NewRunechantAura(1))
	s.ResolveChainStep(s.Logger(), &card.CardState{Card: cards.BloodspillInvocationRed{}})
	if got := s.Value(); got != 0 {
		t.Errorf("Play() = %d, want 0 (weapon hits don't trigger Bloodspill; under-block collapses value)", got)
	}
}

// TestBloodspillInvocation_SameTurnPopByRunechant: a blockable attack action still pops
// Bloodspill when a lone Runechant fires alongside — the 1 arcane is likely to slip through.
func TestBloodspillInvocation_SameTurnPopByRunechant(t *testing.T) {
	s := &gameengine.GameEngine{GameState: gameengine.GameStateBuilder().
		SetIncomingDamage(3).
		SetBlockTotal(0).
		SetCardsRemaining([]*card.CardState{{Card: testutils.AttackWithPower{Power: 6}}}).
		Build()}
	s.CreateAura(sim.NewRunechantAura(1))
	s.ResolveChainStep(s.Logger(), &card.CardState{Card: cards.BloodspillInvocationRed{}})
	if got := s.Value(); got != 3 {
		t.Errorf("Play() = %d, want 3 (Attack=6 blockable, 1 Runechant likely to hit)", got)
	}
}
