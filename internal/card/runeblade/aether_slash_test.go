package runeblade

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

func TestAetherSlash_BaseDamage(t *testing.T) {
	// Nothing pitched → just printed power + printed 1 arcane.
	cases := []struct {
		c    card.Card
		want int
	}{
		{AetherSlashRed{}, 5},
		{AetherSlashYellow{}, 4},
		{AetherSlashBlue{}, 3},
	}
	for _, tc := range cases {
		var s card.TurnState
		if got := tc.c.Play(&s); got != tc.want {
			t.Errorf("%s: Play() = %d, want %d", tc.c.Name(), got, tc.want)
		}
	}
}

func TestAetherSlash_NonAttackActionPitchedAddsArcane(t *testing.T) {
	// A non-attack action in Pitched triggers the +1 arcane rider.
	cases := []struct {
		c    card.Card
		want int
	}{
		{AetherSlashRed{}, 6},
		{AetherSlashYellow{}, 5},
		{AetherSlashBlue{}, 4},
	}
	for _, tc := range cases {
		s := card.TurnState{Pitched: []card.Card{stubNonAttack{}}}
		if got := tc.c.Play(&s); got != tc.want {
			t.Errorf("%s: Play() = %d, want %d", tc.c.Name(), got, tc.want)
		}
	}
}

func TestAetherSlash_AttackPitchedDoesNotTrigger(t *testing.T) {
	// Pitching an attack card does NOT satisfy the "non-attack action pitched" rider.
	s := card.TurnState{Pitched: []card.Card{stubRunebladeAttack{}}}
	if got := (AetherSlashRed{}).Play(&s); got != 5 {
		t.Errorf("Aether Slash Red: Play() = %d, want 5 (no rider)", got)
	}
}
