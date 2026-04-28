package cards

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
)

func TestOathOfTheArknight_NoRemainingCards(t *testing.T) {
	s := &sim.TurnState{}
	(OathOfTheArknightRed{}).Play(s, &sim.CardState{Card: OathOfTheArknightRed{}})
	if got := s.Value; got != 1 {
		t.Errorf("Play() = %d, want 1 (Runechant only, no attack to buff)", got)
	}
}

func TestOathOfTheArknight_RunebladeAttackInRemaining(t *testing.T) {
	// Oath always creates a Runechant (+1 damage credited to Oath itself). The +N{p} buff
	// rides on the target's BonusAttack — so Play returns just the Runechant value, and the
	// target's BonusAttack picks up +N.
	cases := []struct {
		c     sim.Card
		bonus int
	}{
		{OathOfTheArknightRed{}, 3},
		{OathOfTheArknightYellow{}, 2},
		{OathOfTheArknightBlue{}, 1},
	}
	for _, tc := range cases {
		target := &sim.CardState{Card: testutils.RunebladeAttack{}}
		s := &sim.TurnState{CardsRemaining: []*sim.CardState{target}}
		tc.c.Play(s, &sim.CardState{Card: tc.c})
		if got := s.Value; got != 1 {
			t.Errorf("%s: Play() = %d, want 1 (Runechant only; +N rides on target's BonusAttack)", tc.c.Name(), got)
		}
		if target.BonusAttack != tc.bonus {
			t.Errorf("%s: target BonusAttack = %d, want %d", tc.c.Name(), target.BonusAttack, tc.bonus)
		}
	}
}

func TestOathOfTheArknight_WeaponCountsAsAttack(t *testing.T) {
	target := &sim.CardState{Card: testutils.RunebladeWeapon{}}
	s := &sim.TurnState{CardsRemaining: []*sim.CardState{target}}
	(OathOfTheArknightRed{}).Play(s, &sim.CardState{Card: OathOfTheArknightRed{}})
	if got := s.Value; got != 1 {
		t.Errorf("Play() = %d, want 1 (Runechant only; +3 rides on weapon's BonusAttack)", got)
	}
	if target.BonusAttack != 3 {
		t.Errorf("weapon BonusAttack = %d, want 3", target.BonusAttack)
	}
}

func TestOathOfTheArknight_NonRunebladeAttackDoesNotQualify(t *testing.T) {
	s := &sim.TurnState{CardsRemaining: []*sim.CardState{{Card: testutils.NonRunebladeAttack{}}}}
	(OathOfTheArknightRed{}).Play(s, &sim.CardState{Card: OathOfTheArknightRed{}})
	if got := s.Value; got != 1 {
		t.Errorf("Play() = %d, want 1 (non-Runeblade attack shouldn't trigger bonus)", got)
	}
}

func TestOathOfTheArknight_RunebladeNonAttackDoesNotQualify(t *testing.T) {
	// Read the Runes is Runeblade + Action but NOT Attack or Weapon.
	s := &sim.TurnState{CardsRemaining: []*sim.CardState{{Card: testutils.NonAttack{}}}}
	(OathOfTheArknightRed{}).Play(s, &sim.CardState{Card: OathOfTheArknightRed{}})
	if got := s.Value; got != 1 {
		t.Errorf("Play() = %d, want 1 (non-attack Runeblade card shouldn't trigger bonus)", got)
	}
}
