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

func TestDeathlyDuet_NonAttackActionPitchedCreatesRunechants(t *testing.T) {
	// Non-attack action pitched → 2 Runechant tokens enter play. Play() returns just the base
	// power (the Runechant damage lands when some attack consumes the tokens — here it'd be
	// Deathly Duet's own attack, done downstream by playSequence).
	s := card.TurnState{Pitched: []card.Card{stubNonAttack{}}}
	if got := (DeathlyDuetRed{}).Play(&s); got != 4 {
		t.Errorf("Deathly Duet Red with non-attack pitched: Play() = %d, want 4", got)
	}
	if s.Runechants != 2 {
		t.Errorf("Runechants = %d, want 2", s.Runechants)
	}
	if !s.AuraCreated {
		t.Errorf("AuraCreated should be set when Runechants are created")
	}
}

func TestDeathlyDuet_BothBranchesFire(t *testing.T) {
	// Both an attack AND a non-attack action in Pitched → both riders fire: +2 power bonus baked
	// into Play's return, plus 2 Runechants on state (consumed later when an attack resolves).
	s := card.TurnState{Pitched: []card.Card{stubRunebladeAttack{}, stubNonAttack{}}}
	if got := (DeathlyDuetRed{}).Play(&s); got != 6 {
		t.Errorf("Deathly Duet Red with both pitched: Play() = %d, want 6", got)
	}
	if s.Runechants != 2 {
		t.Errorf("Runechants = %d, want 2", s.Runechants)
	}
}
