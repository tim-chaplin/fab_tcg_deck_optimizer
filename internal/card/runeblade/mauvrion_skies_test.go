package runeblade

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

// Reuses stubRunebladeAttack / stubRunebladeWeapon / stubAttackWithPower from stubs_test.go —
// all tests live in the same package.

// playMauvrion exercises Mauvrion Skies's Play and returns the combined damage: Play's own
// return plus whatever damage the registered ephemeral trigger credits when fired against
// target (nil target means "no attack action follows", so the trigger just fizzles).
func playMauvrion(c card.Card, s *card.TurnState, target *card.CardState) int {
	dmg := c.Play(s, &card.CardState{Card: c})
	if target == nil {
		return dmg
	}
	for _, t := range s.EphemeralAttackTriggers {
		if t.Matches != nil && !t.Matches(target) {
			continue
		}
		dmg += t.Handler(s, target)
	}
	s.EphemeralAttackTriggers = nil
	return dmg
}

func TestMauvrionSkies_NoNextAttackReturnsZero(t *testing.T) {
	// No qualifying next attack in the chain → go-again grant fizzles, the ephemeral
	// trigger stays registered but never fires (fizzles at end of turn without a matching
	// attack).
	s := card.TurnState{} // no CardsRemaining
	if got := (MauvrionSkiesRed{}).Play(&s, &card.CardState{}); got != 0 {
		t.Fatalf("Play returned %d, want 0 (rune creation is deferred to trigger fire)", got)
	}
	if s.AuraCreated {
		t.Fatalf("AuraCreated should stay false until the trigger actually fires and creates runes")
	}
	if len(s.EphemeralAttackTriggers) != 1 {
		t.Fatalf("EphemeralAttackTriggers len = %d, want 1 (trigger always registered)", len(s.EphemeralAttackTriggers))
	}
}

func TestMauvrionSkies_WeaponNextDoesNotQualify(t *testing.T) {
	// A Runeblade weapon swing later in the turn is not an attack action card — the
	// go-again grant skips it (weapons already get go-again implicitly via the weapon swing
	// slot), and the ephemeral trigger's Matches predicate rejects it too.
	target := &card.CardState{Card: stubRunebladeWeapon{}}
	s := card.TurnState{CardsRemaining: []*card.CardState{target}}
	if got := (MauvrionSkiesRed{}).Play(&s, &card.CardState{}); got != 0 {
		t.Fatalf("Play returned %d, want 0", got)
	}
	if target.GrantedGoAgain {
		t.Fatalf("grant should stay false when target is a weapon")
	}
	// The trigger must still be registered — its Matches predicate, not Play, is what
	// filters weapon-only targets, so the trigger just never finds a match.
	if len(s.EphemeralAttackTriggers) != 1 {
		t.Fatalf("EphemeralAttackTriggers len = %d, want 1", len(s.EphemeralAttackTriggers))
	}
	if s.EphemeralAttackTriggers[0].Matches(target) {
		t.Fatalf("trigger's Matches(%s) = true, want false (weapon target shouldn't fire the rider)",
			target.Card.Name())
	}
}

func TestMauvrionSkies_NonRunebladeAttackDoesNotQualify(t *testing.T) {
	// A Generic attack action card later in the turn is not a Runeblade attack — the
	// go-again look-ahead skips it, and the ephemeral trigger's Matches predicate rejects
	// it at fire time so the Runechant rider doesn't land.
	target := &card.CardState{Card: stubNonRunebladeAttack{}}
	s := card.TurnState{CardsRemaining: []*card.CardState{target}}
	if got := playMauvrion(MauvrionSkiesRed{}, &s, target); got != 0 {
		t.Fatalf("combined Play+trigger damage = %d, want 0", got)
	}
	if target.GrantedGoAgain {
		t.Fatalf("grant should stay false when target isn't a Runeblade card")
	}
	if s.Runechants != 0 || s.AuraCreated {
		t.Fatalf("state leaked runes/aura: Runechants=%d AuraCreated=%v", s.Runechants, s.AuraCreated)
	}
}

func TestMauvrionSkies_LikelyHitTargetGrantsGoAgainAndRunechants(t *testing.T) {
	// Target's printed power (4) is in the likely-to-hit set. Play grants go-again
	// immediately via the look-ahead; the ephemeral trigger fires on the target's
	// resolution and creates each variant's printed number of Runechant tokens. Tokens are
	// credited +1 each at creation; the test collapses Play + trigger fire into a single
	// damage sum via playMauvrion.
	cases := []struct {
		c card.Card
		n int
	}{
		{MauvrionSkiesRed{}, 3},
		{MauvrionSkiesYellow{}, 2},
		{MauvrionSkiesBlue{}, 1},
	}
	for _, tc := range cases {
		target := &card.CardState{Card: stubAttackWithPower{power: 4}}
		s := card.TurnState{CardsRemaining: []*card.CardState{target}}
		if got := playMauvrion(tc.c, &s, target); got != tc.n {
			t.Errorf("%s: combined damage = %d, want %d", tc.c.Name(), got, tc.n)
		}
		if s.Runechants != tc.n {
			t.Errorf("%s: Runechants = %d, want %d", tc.c.Name(), s.Runechants, tc.n)
		}
		if !target.GrantedGoAgain {
			t.Errorf("%s: target GrantedGoAgain should be set (look-ahead runs during Play)", tc.c.Name())
		}
		if !s.AuraCreated {
			t.Errorf("%s: AuraCreated should be set once the trigger fires", tc.c.Name())
		}
	}
}

func TestMauvrionSkies_BlockableTargetGrantsGoAgainButNoRunechants(t *testing.T) {
	// Target's printed power (3) falls in the blockable range → the trigger's Handler runs
	// LikelyToHit and drops the rider (no Runechants created). The go-again grant still
	// lands because it isn't gated on the attack hitting.
	target := &card.CardState{Card: stubAttackWithPower{power: 3}}
	s := card.TurnState{CardsRemaining: []*card.CardState{target}}
	if got := playMauvrion(MauvrionSkiesRed{}, &s, target); got != 0 {
		t.Errorf("combined damage = %d, want 0 (blockable target drops the Runechants)", got)
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
