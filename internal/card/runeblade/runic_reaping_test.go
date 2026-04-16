package runeblade

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

func TestRunicReaping_NoNextAttackReturnsZero(t *testing.T) {
	// No attack action following → no bonus at all, and AuraCreated must remain false.
	s := card.TurnState{Pitched: []card.Card{stubRunebladeAttack{}}}
	if got := (RunicReapingRed{}).Play(&s); got != 0 {
		t.Fatalf("want 0 when no next attack, got %d", got)
	}
	if s.AuraCreated {
		t.Fatalf("AuraCreated should stay false when no bonus fires")
	}
}

func TestRunicReaping_WeaponNextDoesNotQualify(t *testing.T) {
	// A Runeblade weapon swing later in the turn is not an attack action card, so the rider doesn't
	// trigger.
	s := card.TurnState{CardsRemaining: []*card.PlayedCard{{Card: stubRunebladeWeapon{}}}}
	if got := (RunicReapingRed{}).Play(&s); got != 0 {
		t.Fatalf("want 0 with weapon-only next, got %d", got)
	}
}

func TestRunicReaping_NextAttackNoPitchedAttack(t *testing.T) {
	// Next attack exists, but nothing attack-typed was pitched → N Runechant tokens created on
	// state (consumed downstream by the attack pipeline). Play itself returns 0 damage; the +1
	// pitched-attack rider also doesn't fire.
	cases := []struct {
		c             card.Card
		wantRunechant int
	}{
		{RunicReapingRed{}, 3},
		{RunicReapingYellow{}, 2},
		{RunicReapingBlue{}, 1},
	}
	for _, tc := range cases {
		s := card.TurnState{
			CardsRemaining: []*card.PlayedCard{{Card: stubRunebladeAttack{}}},
			Pitched:        []card.Card{stubNonAttack{}},
		}
		if got := tc.c.Play(&s); got != 0 {
			t.Errorf("%s: Play() = %d, want 0", tc.c.Name(), got)
		}
		if s.Runechants != tc.wantRunechant {
			t.Errorf("%s: Runechants = %d, want %d", tc.c.Name(), s.Runechants, tc.wantRunechant)
		}
		if !s.AuraCreated {
			t.Errorf("%s: AuraCreated should be set when bonus fires", tc.c.Name())
		}
	}
}

func TestRunicReaping_NextAttackWithPitchedAttack(t *testing.T) {
	// Next attack exists AND an attack card was pitched → Play returns the +1{p} rider damage;
	// the Runechant count is unaffected by that rider (only by the create-runechants branch).
	cases := []struct {
		c             card.Card
		wantRunechant int
	}{
		{RunicReapingRed{}, 3},
		{RunicReapingYellow{}, 2},
		{RunicReapingBlue{}, 1},
	}
	for _, tc := range cases {
		s := card.TurnState{
			CardsRemaining: []*card.PlayedCard{{Card: stubRunebladeAttack{}}},
			Pitched:        []card.Card{stubRunebladeAttack{}},
		}
		if got := tc.c.Play(&s); got != 1 {
			t.Errorf("%s: Play() = %d, want 1 (pitched-attack bonus only)", tc.c.Name(), got)
		}
		if s.Runechants != tc.wantRunechant {
			t.Errorf("%s: Runechants = %d, want %d", tc.c.Name(), s.Runechants, tc.wantRunechant)
		}
	}
}
