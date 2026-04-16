package runeblade

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

// Reuses stubRunebladeAttack / stubRunebladeWeapon from runic_reaping_test.go — both tests live in
// the same package.

func TestMauvrionSkies_NoNextAttackReturnsZero(t *testing.T) {
	// No qualifying next attack → both the runechant rider and the go-again grant fizzle.
	s := card.TurnState{} // no CardsRemaining
	if got := (MauvrionSkiesRed{}).Play(&s); got != 0 {
		t.Fatalf("want 0 when no next attack, got %d", got)
	}
	if s.AuraCreated {
		t.Fatalf("AuraCreated should stay false when no bonus fires")
	}
}

func TestMauvrionSkies_WeaponNextDoesNotQualify(t *testing.T) {
	// A Runeblade weapon swing later in the turn is not an attack action card — rider fizzles.
	target := &card.PlayedCard{Card: stubRunebladeWeapon{}}
	s := card.TurnState{CardsRemaining: []*card.PlayedCard{target}}
	if got := (MauvrionSkiesRed{}).Play(&s); got != 0 {
		t.Fatalf("want 0 with weapon-only next, got %d", got)
	}
	if target.GrantedGoAgain {
		t.Fatalf("grant should stay false when target isn't an attack action card")
	}
}

func TestMauvrionSkies_NextAttackGrantsGoAgainAndRunechants(t *testing.T) {
	// Qualifying next attack exists → each variant creates its printed number of Runechant
	// tokens. Play returns N (tokens credited +1 each at creation); state.Runechants holds the
	// tokens for downstream consume; target.GrantedGoAgain is flipped; AuraCreated is set.
	cases := []struct {
		c card.Card
		n int
	}{
		{MauvrionSkiesRed{}, 3},
		{MauvrionSkiesYellow{}, 2},
		{MauvrionSkiesBlue{}, 1},
	}
	for _, tc := range cases {
		target := &card.PlayedCard{Card: stubRunebladeAttack{}}
		s := card.TurnState{CardsRemaining: []*card.PlayedCard{target}}
		if got := tc.c.Play(&s); got != tc.n {
			t.Errorf("%s: Play() = %d, want %d", tc.c.Name(), got, tc.n)
		}
		if s.Runechants != tc.n {
			t.Errorf("%s: Runechants = %d, want %d", tc.c.Name(), s.Runechants, tc.n)
		}
		if !target.GrantedGoAgain {
			t.Errorf("%s: target GrantedGoAgain should be set", tc.c.Name())
		}
		if !s.AuraCreated {
			t.Errorf("%s: AuraCreated should be set", tc.c.Name())
		}
	}
}
