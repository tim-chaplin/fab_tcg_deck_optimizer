package cards

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

func TestCondemnToSlaughter_NoNextAttackReturnsZero(t *testing.T) {
	// No Runeblade attack follows → rider doesn't fire, Play returns 0.
	var s card.TurnState
	(CondemnToSlaughterRed{}).Play(&s, &card.CardState{Card: CondemnToSlaughterRed{}})
	if got := s.Value; got != 0 {
		t.Errorf("Play() = %d, want 0 when CardsRemaining is empty", got)
	}
}

func TestCondemnToSlaughter_NextAttackActionTriggers(t *testing.T) {
	// A Runeblade attack action card in CardsRemaining picks up +N{p} on its BonusAttack;
	// Play returns 0 (the +N attributes to the buffed attack, not Condemn).
	cases := []struct {
		c card.Card
		n int
	}{
		{CondemnToSlaughterRed{}, 3},
		{CondemnToSlaughterYellow{}, 2},
		{CondemnToSlaughterBlue{}, 1},
	}
	for _, tc := range cases {
		target := &card.CardState{Card: stubRunebladeAttack{}}
		s := card.TurnState{CardsRemaining: []*card.CardState{target}}
		tc.c.Play(&s, &card.CardState{Card: tc.c})
		if got := s.Value; got != 0 {
			t.Errorf("%s: Play() = %d, want 0 (granter returns 0; +N rides on target's BonusAttack)", tc.c.Name(), got)
		}
		if target.BonusAttack != tc.n {
			t.Errorf("%s: target BonusAttack = %d, want %d", tc.c.Name(), target.BonusAttack, tc.n)
		}
	}
}

func TestCondemnToSlaughter_WeaponCountsAsNextAttack(t *testing.T) {
	// Unlike Runic Reaping, Condemn's rider accepts weapon swings as the "next attack."
	target := &card.CardState{Card: stubRunebladeWeapon{}}
	s := card.TurnState{CardsRemaining: []*card.CardState{target}}
	(CondemnToSlaughterRed{}).Play(&s, &card.CardState{Card: CondemnToSlaughterRed{}})
	if got := s.Value; got != 0 {
		t.Errorf("Play() = %d, want 0 (granter returns 0; +N rides on target's BonusAttack)", got)
	}
	if target.BonusAttack != 3 {
		t.Errorf("target BonusAttack = %d, want 3 (weapon should qualify)", target.BonusAttack)
	}
}

func TestCondemnToSlaughter_NonRunebladeAttackDoesNotQualify(t *testing.T) {
	// A Generic attack-action card later in the chain doesn't satisfy the Runeblade-only rider.
	s := card.TurnState{CardsRemaining: []*card.CardState{{Card: stubNonRunebladeAttack{}}}}
	(CondemnToSlaughterRed{}).Play(&s, &card.CardState{Card: CondemnToSlaughterRed{}})
	if got := s.Value; got != 0 {
		t.Errorf("Play() = %d, want 0 (non-Runeblade attack shouldn't qualify)", got)
	}
}
