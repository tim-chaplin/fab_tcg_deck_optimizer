package turntests

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/v2/card/cards"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/gameengine"
)

func TestAetherSlash_BaseDamage(t *testing.T) {
	// Nothing attributed to this card → just printed power. The CSV "Arcane: 1" is the text
	// rider'ge damage (not a separate baseline), so with the non-attack-action condition unmet
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
		ge := gameengine.New()
		ge.ResolveChainStep(ge.Logger(), &card.CardState{Card: tc.c})
		if got := ge.Value(); got != tc.want {
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
		ge := gameengine.New()
		self := &card.CardState{Card: tc.c, PitchedToPlay: []card.Card{testutils.NonAttack{}}}
		ge.ResolveChainStep(ge.Logger(), self)
		if got := ge.Value(); got != tc.want {
			t.Errorf("%s: Play() = %d, want %d", tc.c.Name(), got, tc.want)
		}
	}
}

func TestAetherSlash_AttackAttributedDoesNotFireRider(t *testing.T) {
	// Pitch attribution containing only an attack-typed card does NOT satisfy the rider —
	// even if a non-attack action is present in the broader pitch bag (ge.Pitched()), only the
	// cards funded specifically to play this Aether Slash (PitchedToPlay) count.
	self := &card.CardState{
		Card:          cards.AetherSlashRed{},
		PitchedToPlay: []card.Card{testutils.RunebladeAttack{}},
	}
	ge := &gameengine.GameEngine{GameState: gameengine.GameStateBuilder().SetPitched([]card.Card{testutils.RunebladeAttack{}, testutils.NonAttack{}}).Build()}
	ge.ResolveChainStep(ge.Logger(), self)
	if got := ge.Value(); got != 4 {
		t.Errorf("Aether Slash Red: Play() = %d, want 4 (attack attributed; rider gated to PitchedToPlay)", got)
	}
}

func TestAetherSlash_FlagsArcaneDamageDealtOnlyWhenTriggered(t *testing.T) {
	// The ArcaneDamageDealt flag should only be set when the rider actually fires — otherwise
	// same-turn triggers like Meat and Greet'ge go-again would spuriously enable themselves.
	ge := gameengine.New()
	ge.ResolveChainStep(ge.Logger(), &card.CardState{Card: cards.AetherSlashRed{}})
	if ge.ArcaneDamageDealt() {
		t.Error("ArcaneDamageDealt = true with no qualifying pitch attribution; want false")
	}
	ge = gameengine.New()
	self := &card.CardState{Card: cards.AetherSlashRed{}, PitchedToPlay: []card.Card{testutils.NonAttack{}}}
	ge.ResolveChainStep(ge.Logger(), self)
	if !ge.ArcaneDamageDealt() {
		t.Error("ArcaneDamageDealt = false with non-attack action attributed; want true")
	}
}
