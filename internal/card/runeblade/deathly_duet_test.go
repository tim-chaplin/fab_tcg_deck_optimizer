package runeblade

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

func TestDeathlyDuet_BaseDamage(t *testing.T) {
	// Nothing pitched → just printed power.
	cases := []struct {
		c    card.Card
		want int
	}{
		{DeathlyDuetRed{}, 4},
		{DeathlyDuetYellow{}, 3},
		{DeathlyDuetBlue{}, 2},
	}
	for _, tc := range cases {
		var s card.TurnState
		if got := tc.c.Play(&s); got != tc.want {
			t.Errorf("%s: Play() = %d, want %d", tc.c.Name(), got, tc.want)
		}
	}
}

func TestDeathlyDuet_AttackPitchedAddsPower(t *testing.T) {
	// Attack pitched → +2{p}.
	s := card.TurnState{Pitched: []card.Card{stubRunebladeAttack{}}}
	if got := (DeathlyDuetRed{}).Play(&s); got != 6 {
		t.Errorf("Deathly Duet Red with attack pitched: Play() = %d, want 6", got)
	}
}

func TestDeathlyDuet_NonAttackActionPitchedWithFollowupAddsRunechants(t *testing.T) {
	// Non-attack action pitched + a following attack → +2 damage (two Runechants) and
	// AuraCreated set for following aura-conditional cards.
	s := card.TurnState{
		Pitched:        []card.Card{stubNonAttack{}},
		CardsRemaining: []*card.PlayedCard{{Card: stubRunebladeAttack{}}},
	}
	if got := (DeathlyDuetRed{}).Play(&s); got != 6 {
		t.Errorf("Deathly Duet Red with non-attack pitched + follow-up: Play() = %d, want 6", got)
	}
	if !s.AuraCreated {
		t.Errorf("AuraCreated should be set when Runechants are created")
	}
}

func TestDeathlyDuet_NonAttackActionPitchedNoFollowupFizzles(t *testing.T) {
	// Non-attack action pitched but nothing follows → Runechants fizzle, only base damage.
	s := card.TurnState{Pitched: []card.Card{stubNonAttack{}}}
	if got := (DeathlyDuetRed{}).Play(&s); got != 4 {
		t.Errorf("Deathly Duet Red, no follow-up: Play() = %d, want 4", got)
	}
	if s.AuraCreated {
		t.Errorf("AuraCreated should NOT be set when Runechants have no follow-up")
	}
}

func TestDeathlyDuet_WeaponCountsAsFollowingAttack(t *testing.T) {
	// A weapon after Deathly Duet counts as a following attack — Runechants land on the swing.
	s := card.TurnState{
		Pitched:        []card.Card{stubNonAttack{}},
		CardsRemaining: []*card.PlayedCard{{Card: stubRunebladeWeapon{}}},
	}
	if got := (DeathlyDuetRed{}).Play(&s); got != 6 {
		t.Errorf("Deathly Duet Red with weapon follow-up: Play() = %d, want 6", got)
	}
}

func TestDeathlyDuet_BothBranchesFire(t *testing.T) {
	// Both an attack AND a non-attack action in Pitched → both riders fire (+2{p} AND +2 from
	// Runechants when there's a follow-up).
	s := card.TurnState{
		Pitched: []card.Card{stubRunebladeAttack{}, stubNonAttack{}},
		CardsRemaining: []*card.PlayedCard{{Card: stubRunebladeAttack{}}},
	}
	if got := (DeathlyDuetRed{}).Play(&s); got != 8 {
		t.Errorf("Deathly Duet Red with both pitched + follow-up: Play() = %d, want 8", got)
	}
}
