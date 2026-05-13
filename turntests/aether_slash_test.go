package turntests

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/cards"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/gameengine"
)

func TestAetherSlash_BaseDamage(t *testing.T) {
	// Nothing attributed to this card → just printed power. The CSV "Arcane: 1" is the text
	// rider's damage (not a separate baseline), so with the non-attack-action condition unmet
	// the card deals no arcane.
	cases := []struct {
		c    card.Card
		want int
	}{
		{cards.AetherSlashRed{}, 4},
		{cards.AetherSlashYellow{}, 3},
		{cards.AetherSlashBlue{}, 2},
	}
	for _, tc := range cases {
		s := &gameengine.GameEngine{GameState: gameengine.NewState()}
		s.ResolveChainStep(s.Logger(), &card.CardState{Card: tc.c})
		if got := s.Value(); got != tc.want {
			t.Errorf("%s: Play() = %d, want %d", tc.c.Name(), got, tc.want)
		}
	}
}

func TestAetherSlash_NonAttackActionAttributedFiresRider(t *testing.T) {
	// A non-attack action attributed to this card via PitchedToPlay fires the +1 arcane rider.
	cases := []struct {
		c    card.Card
		want int
	}{
		{cards.AetherSlashRed{}, 5},
		{cards.AetherSlashYellow{}, 4},
		{cards.AetherSlashBlue{}, 3},
	}
	for _, tc := range cases {
		s := &gameengine.GameEngine{GameState: gameengine.NewState()}
		self := &card.CardState{Card: tc.c, PitchedToPlay: []card.Card{testutils.NonAttack{}}}
		s.ResolveChainStep(s.Logger(), self)
		if got := s.Value(); got != tc.want {
			t.Errorf("%s: Play() = %d, want %d", tc.c.Name(), got, tc.want)
		}
	}
}

func TestAetherSlash_AttackAttributedDoesNotFireRider(t *testing.T) {
	// Pitch attribution containing only an attack-typed card does NOT satisfy the rider —
	// even if a non-attack action is present in the broader pitch bag (s.Pitched()), only the
	// cards funded specifically to play this Aether Slash (PitchedToPlay) count.
	self := &card.CardState{
		Card:          cards.AetherSlashRed{},
		PitchedToPlay: []card.Card{testutils.RunebladeAttack{}},
	}
	s := gameengine.NewFromSpec(gameengine.Spec{Pitched: []card.Card{testutils.RunebladeAttack{}, testutils.NonAttack{}}})
	s.ResolveChainStep(s.Logger(), self)
	if got := s.Value(); got != 4 {
		t.Errorf("Aether Slash Red: Play() = %d, want 4 (attack attributed; rider gated to PitchedToPlay)", got)
	}
}

func TestAetherSlash_FlagsArcaneDamageDealtOnlyWhenTriggered(t *testing.T) {
	// The ArcaneDamageDealt flag should only be set when the rider actually fires — otherwise
	// same-turn triggers like Meat and Greet's go-again would spuriously enable themselves.
	s := &gameengine.GameEngine{GameState: gameengine.NewState()}
	s.ResolveChainStep(s.Logger(), &card.CardState{Card: cards.AetherSlashRed{}})
	if s.ArcaneDamageDealt() {
		t.Error("ArcaneDamageDealt = true with no qualifying pitch attribution; want false")
	}
	s = &gameengine.GameEngine{GameState: gameengine.NewState()}
	self := &card.CardState{Card: cards.AetherSlashRed{}, PitchedToPlay: []card.Card{testutils.NonAttack{}}}
	s.ResolveChainStep(s.Logger(), self)
	if !s.ArcaneDamageDealt() {
		t.Error("ArcaneDamageDealt = false with non-attack action attributed; want true")
	}
}
