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
	// Next attack exists, but nothing attack-typed was pitched → just N runechants. Each variant
	// contributes its printed count.
	cases := []struct {
		c    card.Card
		want int
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
		if got := tc.c.Play(&s); got != tc.want {
			t.Errorf("%s: Play() = %d, want %d", tc.c.Name(), got, tc.want)
		}
		if !s.AuraCreated {
			t.Errorf("%s: AuraCreated should be set when bonus fires", tc.c.Name())
		}
	}
}

func TestRunicReaping_NextAttackWithPitchedAttack(t *testing.T) {
	// Next attack exists AND an attack card was pitched → N+1 (the +1{p} rider stacks on the
	// runechant count).
	cases := []struct {
		c    card.Card
		want int
	}{
		{RunicReapingRed{}, 4},
		{RunicReapingYellow{}, 3},
		{RunicReapingBlue{}, 2},
	}
	for _, tc := range cases {
		s := card.TurnState{
			CardsRemaining: []*card.PlayedCard{{Card: stubRunebladeAttack{}}},
			Pitched:        []card.Card{stubRunebladeAttack{}},
		}
		if got := tc.c.Play(&s); got != tc.want {
			t.Errorf("%s: Play() = %d, want %d", tc.c.Name(), got, tc.want)
		}
	}
}
