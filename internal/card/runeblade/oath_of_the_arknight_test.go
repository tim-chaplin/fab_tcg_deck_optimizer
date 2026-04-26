package runeblade

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

func TestOathOfTheArknight_NoRemainingCards(t *testing.T) {
	s := &card.TurnState{}
	(OathOfTheArknightRed{}).Play(s, &card.CardState{Card: OathOfTheArknightRed{}})
	if got := s.Value; got != 1 {
		t.Errorf("Play() = %d, want 1 (Runechant only, no attack to buff)", got)
	}
}

func TestOathOfTheArknight_RunebladeAttackInRemaining(t *testing.T) {
	// Oath always creates a Runechant (+1 damage credited to Oath itself). The +N{p} buff
	// rides on the target's BonusAttack — so Play returns just the Runechant value, and the
	// target's BonusAttack picks up +N.
	cases := []struct {
		c     card.Card
		bonus int
	}{
		{OathOfTheArknightRed{}, 3},
		{OathOfTheArknightYellow{}, 2},
		{OathOfTheArknightBlue{}, 1},
	}
	for _, tc := range cases {
		target := &card.CardState{Card: stubRunebladeAttack{}}
		s := &card.TurnState{CardsRemaining: []*card.CardState{target}}
		tc.c.Play(s, &card.CardState{Card: tc.c})
		if got := s.Value; got != 1 {
			t.Errorf("%s: Play() = %d, want 1 (Runechant only; +N rides on target's BonusAttack)", tc.c.Name(), got)
		}
		if target.BonusAttack != tc.bonus {
			t.Errorf("%s: target BonusAttack = %d, want %d", tc.c.Name(), target.BonusAttack, tc.bonus)
		}
	}
}

func TestOathOfTheArknight_WeaponCountsAsAttack(t *testing.T) {
	target := &card.CardState{Card: stubRunebladeWeapon{}}
	s := &card.TurnState{CardsRemaining: []*card.CardState{target}}
	(OathOfTheArknightRed{}).Play(s, &card.CardState{Card: OathOfTheArknightRed{}})
	if got := s.Value; got != 1 {
		t.Errorf("Play() = %d, want 1 (Runechant only; +3 rides on weapon's BonusAttack)", got)
	}
	if target.BonusAttack != 3 {
		t.Errorf("weapon BonusAttack = %d, want 3", target.BonusAttack)
	}
}

func TestOathOfTheArknight_NonRunebladeAttackDoesNotQualify(t *testing.T) {
	s := &card.TurnState{CardsRemaining: []*card.CardState{{Card: stubNonRunebladeAttack{}}}}
	(OathOfTheArknightRed{}).Play(s, &card.CardState{Card: OathOfTheArknightRed{}})
	if got := s.Value; got != 1 {
		t.Errorf("Play() = %d, want 1 (non-Runeblade attack shouldn't trigger bonus)", got)
	}
}

func TestOathOfTheArknight_RunebladeNonAttackDoesNotQualify(t *testing.T) {
	// Read the Runes is Runeblade + Action but NOT Attack or Weapon.
	s := &card.TurnState{CardsRemaining: []*card.CardState{{Card: stubNonAttack{}}}}
	(OathOfTheArknightRed{}).Play(s, &card.CardState{Card: OathOfTheArknightRed{}})
	if got := s.Value; got != 1 {
		t.Errorf("Play() = %d, want 1 (non-attack Runeblade card shouldn't trigger bonus)", got)
	}
}
