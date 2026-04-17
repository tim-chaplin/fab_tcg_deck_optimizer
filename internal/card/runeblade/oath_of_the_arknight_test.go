package runeblade

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

func TestOathOfTheArknight_NoRemainingCards(t *testing.T) {
	s := &card.TurnState{}
	if got := (OathOfTheArknightRed{}).Play(s); got != 1 {
		t.Errorf("Play() = %d, want 1 (Runechant only, no attack to buff)", got)
	}
}

func TestOathOfTheArknight_RunebladeAttackInRemaining(t *testing.T) {
	cases := []struct {
		c    card.Card
		want int
	}{
		{OathOfTheArknightRed{}, 1 + 3},
		{OathOfTheArknightYellow{}, 1 + 2},
		{OathOfTheArknightBlue{}, 1 + 1},
	}
	for _, tc := range cases {
		s := &card.TurnState{CardsRemaining: []*card.PlayedCard{{Card: stubRunebladeAttack{}}}}
		if got := tc.c.Play(s); got != tc.want {
			t.Errorf("%s: Play() = %d, want %d", tc.c.Name(), got, tc.want)
		}
	}
}

func TestOathOfTheArknight_WeaponCountsAsAttack(t *testing.T) {
	s := &card.TurnState{CardsRemaining: []*card.PlayedCard{{Card: stubRunebladeWeapon{}}}}
	if got := (OathOfTheArknightRed{}).Play(s); got != 4 {
		t.Errorf("Play() = %d, want 4 (1 Runechant + 3 bonus from weapon)", got)
	}
}

func TestOathOfTheArknight_NonRunebladeAttackDoesNotQualify(t *testing.T) {
	s := &card.TurnState{CardsRemaining: []*card.PlayedCard{{Card: stubNonRunebladeAttack{}}}}
	if got := (OathOfTheArknightRed{}).Play(s); got != 1 {
		t.Errorf("Play() = %d, want 1 (non-Runeblade attack shouldn't trigger bonus)", got)
	}
}

func TestOathOfTheArknight_RunebladeNonAttackDoesNotQualify(t *testing.T) {
	// Read the Runes is Runeblade + Action but NOT Attack or Weapon.
	s := &card.TurnState{CardsRemaining: []*card.PlayedCard{{Card: stubNonAttack{}}}}
	if got := (OathOfTheArknightRed{}).Play(s); got != 1 {
		t.Errorf("Play() = %d, want 1 (non-attack Runeblade card shouldn't trigger bonus)", got)
	}
}
