package registry_test

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/weapons"
)

// TestLegalCards_SkipsNotImplemented pins that the registry's deck-construction pool excludes
// every card carrying the NotImplemented marker — the deck builder relies on this filter
// rather than re-checking each printing per call to Random / SanitizeNotImplemented.
func TestLegalCards_SkipsNotImplemented(t *testing.T) {
	for _, c := range (registry.Registry{}).LegalCards() {
		if _, ok := c.(registry.NotImplemented); ok {
			t.Errorf("LegalCards included NotImplemented card %s", c.DisplayName())
		}
	}
}

// TestLegalCards_SkipsUnplayable pins the Unplayable counterpart: any card carrying the
// Unplayable marker is filtered out of the pool, so a deck.Random pick can never land on one.
func TestLegalCards_SkipsUnplayable(t *testing.T) {
	for _, c := range (registry.Registry{}).LegalCards() {
		if _, ok := c.(registry.Unplayable); ok {
			t.Errorf("LegalCards included Unplayable card %s", c.DisplayName())
		}
	}
}

// TestLegalCards_ExcludesTaggedNotImplementedByID gives the marker-property test teeth: a
// known-tagged NotImplemented card must be absent from LegalCards. Self-retires when the
// card loses the tag.
func TestLegalCards_ExcludesTaggedNotImplementedByID(t *testing.T) {
	if _, ok := registry.GetCard(ids.StrikeGoldRed).(registry.NotImplemented); !ok {
		t.Skip("Strike Gold [R] is no longer NotImplemented — pick another tagged card or drop this test")
	}
	for _, c := range (registry.Registry{}).LegalCards() {
		if c.ID() == ids.StrikeGoldRed {
			t.Fatalf("LegalCards included Strike Gold [R] despite its NotImplemented tag")
		}
	}
}

// TestLegalCards_ExcludesTaggedUnplayableByID gives the Unplayable-marker test teeth: a
// known-tagged Unplayable card must be absent from LegalCards. Self-retires when the card
// loses the tag.
func TestLegalCards_ExcludesTaggedUnplayableByID(t *testing.T) {
	if _, ok := registry.GetCard(ids.PotionOfSeeingBlue).(registry.Unplayable); !ok {
		t.Skip("Potion of Seeing [B] is no longer Unplayable — pick another tagged card or drop this test")
	}
	for _, c := range (registry.Registry{}).LegalCards() {
		if c.ID() == ids.PotionOfSeeingBlue {
			t.Fatalf("LegalCards included Potion of Seeing [B] despite its Unplayable tag")
		}
	}
}

// TestLegalWeapons_SkipsNotImplemented pins the weapon-side counterpart to
// TestLegalCards_SkipsNotImplemented: every weapon LegalWeapons surfaces must not implement
// NotImplemented, regardless of which weapons currently carry the tag.
func TestLegalWeapons_SkipsNotImplemented(t *testing.T) {
	for _, w := range (registry.Registry{}).LegalWeapons() {
		if _, ok := w.(registry.NotImplemented); ok {
			t.Errorf("LegalWeapons included NotImplemented weapon %s", w.Name())
		}
	}
}

// TestLegalWeapons_ExcludesTaggedWeaponByID pins that a known-tagged NotImplemented weapon is
// absent from LegalWeapons. Self-retires when the weapon loses the tag.
func TestLegalWeapons_ExcludesTaggedWeaponByID(t *testing.T) {
	tagged := weapons.AnnalsOfSutcliffe{}
	if _, ok := any(tagged).(registry.NotImplemented); !ok {
		t.Skip("Annals of Sutcliffe is no longer NotImplemented — pick another tagged weapon or drop this test")
	}
	for _, w := range (registry.Registry{}).LegalWeapons() {
		if w == tagged {
			t.Fatalf("LegalWeapons included Annals of Sutcliffe despite its NotImplemented tag")
		}
	}
}
