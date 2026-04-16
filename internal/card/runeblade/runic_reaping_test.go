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
	// Next attack exists, but nothing attack-typed was pitched → N Runechant tokens created.
	// Play returns N (each token credited +1 at creation); the pitched-attack +1 rider doesn't
	// fire. state.Runechants tracks the tokens for downstream consume.
	cases := []struct {
		c card.Card
		n int
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
		if got := tc.c.Play(&s); got != tc.n {
			t.Errorf("%s: Play() = %d, want %d", tc.c.Name(), got, tc.n)
		}
		if s.Runechants != tc.n {
			t.Errorf("%s: Runechants = %d, want %d", tc.c.Name(), s.Runechants, tc.n)
		}
		if !s.AuraCreated {
			t.Errorf("%s: AuraCreated should be set when bonus fires", tc.c.Name())
		}
	}
}

func TestRunicReaping_NextAttackWithPitchedAttack(t *testing.T) {
	// Next attack exists AND an attack card was pitched → Play returns N (token credits) plus 1
	// (the pitched-attack rider). state.Runechants holds only the N tokens — the rider damage is
	// direct, not a runechant.
	cases := []struct {
		c card.Card
		n int
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
		if got := tc.c.Play(&s); got != tc.n+1 {
			t.Errorf("%s: Play() = %d, want %d (N tokens + 1 pitched-attack bonus)", tc.c.Name(), got, tc.n+1)
		}
		if s.Runechants != tc.n {
			t.Errorf("%s: Runechants = %d, want %d", tc.c.Name(), s.Runechants, tc.n)
		}
	}
}
