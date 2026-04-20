package runeblade

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

// Reuses stubRunebladeAttack / stubRunebladeWeapon / stubAttackWithPower from stubs_test.go —
// all tests live in the same package.

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
	// A Runeblade weapon swing later in the turn is not an attack action card — rider fizzles,
	// and the go-again grant is skipped too (weapons already get go-again implicitly via the
	// weapon swing slot; Mauvrion's grant targets attack action CARDS only).
	target := &card.PlayedCard{Card: stubRunebladeWeapon{}}
	s := card.TurnState{CardsRemaining: []*card.PlayedCard{target}}
	if got := (MauvrionSkiesRed{}).Play(&s); got != 0 {
		t.Fatalf("want 0 with weapon-only next, got %d", got)
	}
	if target.GrantedGoAgain {
		t.Fatalf("grant should stay false when target is a weapon")
	}
}

func TestMauvrionSkies_NonRunebladeAttackDoesNotQualify(t *testing.T) {
	// A Generic attack action card later in the turn is not a Runeblade attack — Mauvrion's
	// rider is gated on the "next Runeblade attack action card", so the grant must skip it
	// (and no Runechants fire).
	target := &card.PlayedCard{Card: stubNonRunebladeAttack{}}
	s := card.TurnState{CardsRemaining: []*card.PlayedCard{target}}
	if got := (MauvrionSkiesRed{}).Play(&s); got != 0 {
		t.Fatalf("want 0 with non-Runeblade attack next, got %d", got)
	}
	if target.GrantedGoAgain {
		t.Fatalf("grant should stay false when target isn't a Runeblade card")
	}
}

func TestMauvrionSkies_LikelyHitTargetGrantsGoAgainAndRunechants(t *testing.T) {
	// Target's printed power (4) is in the likely-to-hit set → each variant creates its
	// printed number of Runechant tokens. Play returns N (tokens credited +1 each at
	// creation); state.Runechants holds the tokens for downstream consume; target.
	// GrantedGoAgain is flipped; AuraCreated is set.
	cases := []struct {
		c card.Card
		n int
	}{
		{MauvrionSkiesRed{}, 3},
		{MauvrionSkiesYellow{}, 2},
		{MauvrionSkiesBlue{}, 1},
	}
	for _, tc := range cases {
		target := &card.PlayedCard{Card: stubAttackWithPower{power: 4}}
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

func TestMauvrionSkies_BlockableTargetGrantsGoAgainButNoRunechants(t *testing.T) {
	// Target's printed power (3) falls in the blockable range → the "if this hits" clause
	// doesn't fire, so no Runechants are created. The go-again grant still lands because it
	// isn't gated on the attack hitting.
	target := &card.PlayedCard{Card: stubAttackWithPower{power: 3}}
	s := card.TurnState{CardsRemaining: []*card.PlayedCard{target}}
	if got := (MauvrionSkiesRed{}).Play(&s); got != 0 {
		t.Errorf("Play() = %d, want 0 (blockable target drops the Runechants)", got)
	}
	if s.Runechants != 0 {
		t.Errorf("Runechants = %d, want 0", s.Runechants)
	}
	if !target.GrantedGoAgain {
		t.Error("target GrantedGoAgain should still be set — go-again isn't gated on hitting")
	}
	if s.AuraCreated {
		t.Error("AuraCreated should stay false when no Runechant fires")
	}
}
