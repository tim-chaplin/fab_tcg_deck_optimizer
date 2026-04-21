package runeblade

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

// TestBloodspillInvocation_BlockCoversIncomingReturnsN: the aura survives the opponent's turn
// when we block all incoming damage, and pays N on a future-turn pop.
func TestBloodspillInvocation_BlockCoversIncomingReturnsN(t *testing.T) {
	cases := []struct {
		c card.Card
		n int
	}{
		{BloodspillInvocationRed{}, 3},
		{BloodspillInvocationYellow{}, 2},
		{BloodspillInvocationBlue{}, 1},
	}
	for _, tc := range cases {
		s := card.TurnState{IncomingDamage: 3, BlockTotal: 3}
		if got := tc.c.Play(&s, nil); got != tc.n {
			t.Errorf("%s: Play() = %d, want %d (block == incoming)", tc.c.Name(), got, tc.n)
		}
	}
}

// TestBloodspillInvocation_BlockShortReturnsZero: if we take damage and have no same-turn
// attack action likely to hit, Bloodspill dies without creating Runechants.
func TestBloodspillInvocation_BlockShortReturnsZero(t *testing.T) {
	cases := []card.Card{
		BloodspillInvocationRed{},
		BloodspillInvocationYellow{},
		BloodspillInvocationBlue{},
	}
	for _, c := range cases {
		s := card.TurnState{IncomingDamage: 3, BlockTotal: 2}
		if got := c.Play(&s, nil); got != 0 {
			t.Errorf("%s: Play() = %d, want 0 (block < incoming, no same-turn pop)", c.Name(), got)
		}
	}
}

// TestBloodspillInvocation_SameTurnPopBySalientAttackAction: a later attack action with a
// likely-to-hit power pops Bloodspill this turn for its full N — even if we're taking damage.
func TestBloodspillInvocation_SameTurnPopBySalientAttackAction(t *testing.T) {
	s := card.TurnState{
		IncomingDamage: 3,
		BlockTotal:     0,
		CardsRemaining: []*card.PlayedCard{{Card: stubAttackWithPower{power: 4}}},
	}
	if got := (BloodspillInvocationRed{}).Play(&s, nil); got != 3 {
		t.Errorf("Play() = %d, want 3 (Attack=4 attack action pops Bloodspill same turn)", got)
	}
}

// TestBloodspillInvocation_WeaponDoesNotPop: Bloodspill's rider is gated on attack ACTION cards
// hitting — a weapon swing that hits doesn't trigger its destruction, even with a Runechant
// firing alongside.
func TestBloodspillInvocation_WeaponDoesNotPop(t *testing.T) {
	s := card.TurnState{
		IncomingDamage: 3,
		BlockTotal:     0,
		Runechants:     1,
		CardsRemaining: []*card.PlayedCard{{Card: stubRunebladeWeapon{}}},
	}
	if got := (BloodspillInvocationRed{}).Play(&s, nil); got != 0 {
		t.Errorf("Play() = %d, want 0 (weapon hits don't trigger Bloodspill; under-block collapses value)", got)
	}
}

// TestBloodspillInvocation_SameTurnPopByRunechant: a blockable attack action still pops
// Bloodspill when a lone Runechant fires alongside — the 1 arcane is likely to slip through.
func TestBloodspillInvocation_SameTurnPopByRunechant(t *testing.T) {
	s := card.TurnState{
		IncomingDamage: 3,
		BlockTotal:     0,
		Runechants:     1,
		CardsRemaining: []*card.PlayedCard{{Card: stubAttackWithPower{power: 6}}},
	}
	if got := (BloodspillInvocationRed{}).Play(&s, nil); got != 3 {
		t.Errorf("Play() = %d, want 3 (Attack=6 blockable, 1 Runechant likely to hit)", got)
	}
}
