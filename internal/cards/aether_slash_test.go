package cards

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
)

func TestAetherSlash_BaseDamage(t *testing.T) {
	// Nothing attributed to this card → just printed power. The CSV "Arcane: 1" is the text
	// rider's damage (not a separate baseline), so with the non-attack-action condition unmet
	// the card deals no arcane.
	cases := []struct {
		c    sim.Card
		want int
	}{
		{AetherSlashRed{}, 4},
		{AetherSlashYellow{}, 3},
		{AetherSlashBlue{}, 2},
	}
	for _, tc := range cases {
		var s sim.TurnState
		tc.c.Play(&s, &sim.CardState{Card: tc.c})
		if got := s.Value; got != tc.want {
			t.Errorf("%s: Play() = %d, want %d", tc.c.Name(), got, tc.want)
		}
	}
}

func TestAetherSlash_NonAttackActionAttributedFiresRider(t *testing.T) {
	// A non-attack action attributed to this card via PitchedToPlay fires the +1 arcane rider.
	cases := []struct {
		c    sim.Card
		want int
	}{
		{AetherSlashRed{}, 5},
		{AetherSlashYellow{}, 4},
		{AetherSlashBlue{}, 3},
	}
	for _, tc := range cases {
		var s sim.TurnState
		self := &sim.CardState{Card: tc.c, PitchedToPlay: []sim.Card{testutils.NonAttack{}}}
		tc.c.Play(&s, self)
		if got := s.Value; got != tc.want {
			t.Errorf("%s: Play() = %d, want %d", tc.c.Name(), got, tc.want)
		}
	}
}

func TestAetherSlash_AttackAttributedDoesNotFireRider(t *testing.T) {
	// Pitch attribution containing only an attack-typed card does NOT satisfy the rider —
	// even if a non-attack action is present in the broader pitch bag (s.Pitched), only the
	// cards funded specifically to play this Aether Slash (PitchedToPlay) count.
	self := &sim.CardState{
		Card:          AetherSlashRed{},
		PitchedToPlay: []sim.Card{testutils.RunebladeAttack{}},
	}
	s := sim.TurnState{Pitched: []sim.Card{testutils.RunebladeAttack{}, testutils.NonAttack{}}}
	(AetherSlashRed{}).Play(&s, self)
	if got := s.Value; got != 4 {
		t.Errorf("Aether Slash Red: Play() = %d, want 4 (attack attributed; rider gated to PitchedToPlay)", got)
	}
}

func TestAetherSlash_FlagsArcaneDamageDealtOnlyWhenTriggered(t *testing.T) {
	// The ArcaneDamageDealt flag should only be set when the rider actually fires — otherwise
	// same-turn triggers like Meat and Greet's go-again would spuriously enable themselves.
	var s sim.TurnState
	(AetherSlashRed{}).Play(&s, &sim.CardState{Card: AetherSlashRed{}})
	if s.ArcaneDamageDealt {
		t.Error("ArcaneDamageDealt = true with no qualifying pitch attribution; want false")
	}
	s = sim.TurnState{}
	self := &sim.CardState{Card: AetherSlashRed{}, PitchedToPlay: []sim.Card{testutils.NonAttack{}}}
	(AetherSlashRed{}).Play(&s, self)
	if !s.ArcaneDamageDealt {
		t.Error("ArcaneDamageDealt = false with non-attack action attributed; want true")
	}
}
