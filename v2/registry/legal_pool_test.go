package registry_test

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/weapons"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/registry"
)

// Tests that LegalCards excludes every card carrying the NotImplemented marker. The deck
// builder relies on this filter rather than re-checking each printing per call to Random.
func TestLegalCards_SkipsNotImplemented(t *testing.T) {
	for _, c := range (registry.Registry{}).LegalCards() {
		if _, ok := c.(registry.NotImplemented); ok {
			t.Errorf("LegalCards included NotImplemented card %s", c.DisplayName())
		}
	}
}

// Tests that LegalCards excludes every card carrying the Unplayable marker.
func TestLegalCards_SkipsUnplayable(t *testing.T) {
	for _, c := range (registry.Registry{}).LegalCards() {
		if _, ok := c.(registry.Unplayable); ok {
			t.Errorf("LegalCards included Unplayable card %s", c.DisplayName())
		}
	}
}

// Tests that a known-tagged NotImplemented card is absent from LegalCards. Self-retires
// (skip) when the chosen card loses the tag.
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

// Tests that a known-tagged Unplayable card is absent from LegalCards. Self-retires (skip)
// when the chosen card loses the tag.
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

// Tests that LegalWeapons excludes every weapon carrying the NotImplemented marker.
func TestLegalWeapons_SkipsNotImplemented(t *testing.T) {
	for _, w := range (registry.Registry{}).LegalWeapons() {
		if _, ok := w.(registry.NotImplemented); ok {
			t.Errorf("LegalWeapons included NotImplemented weapon %s", w.Name())
		}
	}
}

// Tests that a known-tagged NotImplemented weapon is absent from LegalWeapons. Self-retires
// (skip) when the chosen weapon loses the tag.
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
