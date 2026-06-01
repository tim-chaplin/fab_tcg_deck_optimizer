package registry

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
)

// fakeHero is the minimal card.Hero heroCanPlay reads: a class and a type line (which carries
// the hero's talents). Local so the test can vary the class / talents freely, unlike the
// fixed-class testutils.Hero.
type fakeHero struct {
	class   card.CardType
	typeSet card.TypeSet
}

func (h fakeHero) Class() card.CardType { return h.class }
func (h fakeHero) Types() card.TypeSet  { return h.typeSet }

// universalFake is a card.Card carrying the Universal marker, for the class-fold branch.
type universalFake struct{ testutils.Fake }

func (universalFake) Universal() {}

func TestHeroCanPlay(t *testing.T) {
	// Viserai: Runeblade class, no talent.
	viserai := fakeHero{class: card.TypeRuneblade, typeSet: card.NewTypeSet(card.TypeRuneblade)}
	// Aurora-like: Runeblade class, Lightning talent.
	aurora := fakeHero{
		class:   card.TypeRuneblade,
		typeSet: card.NewTypeSet(card.TypeRuneblade, card.TypeLightning),
	}

	// Cards built off a neutral base shape; only the class / talent bits added via WithTypes
	// matter to heroCanPlay.
	generic := testutils.FakeRedAction().WithTypes(card.TypeGeneric)
	runeblade := testutils.FakeRedAction().WithTypes(card.TypeRuneblade)
	genericLightning := testutils.FakeRedAction().WithTypes(card.TypeGeneric, card.TypeLightning)
	runebladeLightning := testutils.FakeRedAction().WithTypes(card.TypeRuneblade, card.TypeLightning)
	classlessLightning := testutils.FakeRedAction().WithTypes(card.TypeLightning)

	cases := []struct {
		name string
		hero fakeHero
		card card.Card
		want bool
	}{
		{"viserai plays generic", viserai, generic, true},
		{"viserai plays runeblade", viserai, runeblade, true},
		{"viserai rejects generic lightning", viserai, genericLightning, false},
		{"viserai rejects runeblade lightning", viserai, runebladeLightning, false},
		{"viserai rejects classless lightning", viserai, classlessLightning, false},
		{"aurora plays generic", aurora, generic, true},
		{"aurora plays runeblade", aurora, runeblade, true},
		{"aurora plays generic lightning", aurora, genericLightning, true},
		{"aurora plays runeblade lightning", aurora, runebladeLightning, true},
		{"aurora plays classless lightning", aurora, classlessLightning, true},
		// A Generic card is legal for a hero of any class (Generic is always allowed).
		{"runeblade card illegal for off-class hero", fakeHero{class: card.TypeGeneric}, runeblade, false},
		// Universal cards adopt the hero's class, so they're always class-legal.
		{"universal card legal regardless of class", viserai, universalFake{}, true},
	}
	for _, c := range cases {
		if got := heroCanPlay(c.hero, c.card); got != c.want {
			t.Errorf("%s: heroCanPlay = %v, want %v", c.name, got, c.want)
		}
	}
}
