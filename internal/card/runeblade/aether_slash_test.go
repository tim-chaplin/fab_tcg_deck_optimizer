package runeblade

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

func TestAetherSlash_BaseDamage(t *testing.T) {
	// Nothing pitched → just printed power. The CSV "Arcane: 1" is the text rider's damage (not
	// a separate baseline), so with the non-attack-action condition unmet the card deals no
	// arcane.
	cases := []struct {
		c    card.Card
		want int
	}{
		{AetherSlashRed{}, 4},
		{AetherSlashYellow{}, 3},
		{AetherSlashBlue{}, 2},
	}
	for _, tc := range cases {
		var s card.TurnState
		if got := tc.c.Play(&s, nil); got != tc.want {
			t.Errorf("%s: Play() = %d, want %d", tc.c.Name(), got, tc.want)
		}
	}
}

func TestAetherSlash_NonAttackActionPitchedAddsArcane(t *testing.T) {
	// A non-attack action in Pitched fires the text rider for +1 arcane.
	cases := []struct {
		c    card.Card
		want int
	}{
		{AetherSlashRed{}, 5},
		{AetherSlashYellow{}, 4},
		{AetherSlashBlue{}, 3},
	}
	for _, tc := range cases {
		s := card.TurnState{Pitched: []card.Card{stubNonAttack{}}}
		if got := tc.c.Play(&s, nil); got != tc.want {
			t.Errorf("%s: Play() = %d, want %d", tc.c.Name(), got, tc.want)
		}
	}
}

func TestAetherSlash_AttackPitchedDoesNotTrigger(t *testing.T) {
	// Pitching an attack card does NOT satisfy the "non-attack action pitched" rider.
	s := card.TurnState{Pitched: []card.Card{stubRunebladeAttack{}}}
	if got := (AetherSlashRed{}).Play(&s, nil); got != 4 {
		t.Errorf("Aether Slash Red: Play() = %d, want 4 (no rider)", got)
	}
}

func TestAetherSlash_FlagsArcaneDamageDealtOnlyWhenTriggered(t *testing.T) {
	// The ArcaneDamageDealt flag should only be set when the rider actually fires — otherwise
	// same-turn triggers like Meat and Greet's go-again would spuriously enable themselves.
	var s card.TurnState
	(AetherSlashRed{}).Play(&s, nil)
	if s.ArcaneDamageDealt {
		t.Error("ArcaneDamageDealt = true with no qualifying pitch; want false")
	}
	s = card.TurnState{Pitched: []card.Card{stubNonAttack{}}}
	(AetherSlashRed{}).Play(&s, nil)
	if !s.ArcaneDamageDealt {
		t.Error("ArcaneDamageDealt = false with non-attack action pitched; want true")
	}
}
