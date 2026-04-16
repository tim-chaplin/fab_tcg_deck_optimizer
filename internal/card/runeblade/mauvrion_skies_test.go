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
	// tokens on state (consumed when a later attack fires), flips GrantedGoAgain on the target
	// PlayedCard, and sets AuraCreated. Play itself returns 0 damage — the token damage is
	// accounted for downstream by the attack pipeline.
	cases := []struct {
		c             card.Card
		wantRunechant int
	}{
		{MauvrionSkiesRed{}, 3},
		{MauvrionSkiesYellow{}, 2},
		{MauvrionSkiesBlue{}, 1},
	}
	for _, tc := range cases {
		target := &card.PlayedCard{Card: stubRunebladeAttack{}}
		s := card.TurnState{CardsRemaining: []*card.PlayedCard{target}}
		if got := tc.c.Play(&s); got != 0 {
			t.Errorf("%s: Play() = %d, want 0", tc.c.Name(), got)
		}
		if s.Runechants != tc.wantRunechant {
			t.Errorf("%s: Runechants = %d, want %d", tc.c.Name(), s.Runechants, tc.wantRunechant)
		}
		if !target.GrantedGoAgain {
			t.Errorf("%s: target GrantedGoAgain should be set", tc.c.Name())
		}
		if !s.AuraCreated {
			t.Errorf("%s: AuraCreated should be set", tc.c.Name())
		}
	}
}
