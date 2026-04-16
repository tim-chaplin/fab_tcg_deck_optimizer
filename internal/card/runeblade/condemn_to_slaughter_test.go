package runeblade

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

func TestCondemnToSlaughter_NoNextAttackReturnsZero(t *testing.T) {
	// No Runeblade attack follows → rider doesn't fire, Play returns 0.
	var s card.TurnState
	if got := (CondemnToSlaughterRed{}).Play(&s); got != 0 {
		t.Errorf("Play() = %d, want 0 when CardsRemaining is empty", got)
	}
}

func TestCondemnToSlaughter_NextAttackActionTriggers(t *testing.T) {
	// A Runeblade attack action card in CardsRemaining triggers the +N{p} rider.
	cases := []struct {
		c card.Card
		n int
	}{
		{CondemnToSlaughterRed{}, 3},
		{CondemnToSlaughterYellow{}, 2},
		{CondemnToSlaughterBlue{}, 1},
	}
	for _, tc := range cases {
		s := card.TurnState{CardsRemaining: []*card.PlayedCard{{Card: stubRunebladeAttack{}}}}
		if got := tc.c.Play(&s); got != tc.n {
			t.Errorf("%s: Play() = %d, want %d", tc.c.Name(), got, tc.n)
		}
	}
}

func TestCondemnToSlaughter_WeaponCountsAsNextAttack(t *testing.T) {
	// Unlike Runic Reaping, Condemn's rider accepts weapon swings as the "next attack."
	s := card.TurnState{CardsRemaining: []*card.PlayedCard{{Card: stubRunebladeWeapon{}}}}
	if got := (CondemnToSlaughterRed{}).Play(&s); got != 3 {
		t.Errorf("Play() = %d, want 3 (weapon should qualify)", got)
	}
}

func TestCondemnToSlaughter_NonRunebladeAttackDoesNotQualify(t *testing.T) {
	// A Generic attack-action card later in the chain doesn't satisfy the Runeblade-only rider.
	s := card.TurnState{CardsRemaining: []*card.PlayedCard{{Card: stubNonRunebladeAttack{}}}}
	if got := (CondemnToSlaughterRed{}).Play(&s); got != 0 {
		t.Errorf("Play() = %d, want 0 (non-Runeblade attack shouldn't qualify)", got)
	}
}
