package registry

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/weapon/weapons"
)

// Tests that LegalCards excludes every card carrying the NotImplemented marker. The deck
// builder relies on this filter rather than re-checking each printing per call to Random.
func TestLegalCards_SkipsNotImplemented(t *testing.T) {
	for _, c := range (Registry{}).LegalCards() {
		if _, ok := c.(NotImplemented); ok {
			t.Errorf("LegalCards included NotImplemented card %s", c.DisplayName())
		}
	}
}

// Tests that LegalCards excludes every card carrying the Unplayable marker.
func TestLegalCards_SkipsUnplayable(t *testing.T) {
	for _, c := range (Registry{}).LegalCards() {
		if _, ok := c.(Unplayable); ok {
			t.Errorf("LegalCards included Unplayable card %s", c.DisplayName())
		}
	}
}

// Tests that a known-tagged NotImplemented card is absent from LegalCards. Self-retires
// (skip) when the chosen card loses the tag.
func TestLegalCards_ExcludesTaggedNotImplementedByID(t *testing.T) {
	if _, ok := GetCard(ids.StrikeGoldRed).(NotImplemented); !ok {
		t.Skip("Strike Gold [R] is no longer NotImplemented — pick another tagged card or drop this test")
	}
	for _, c := range (Registry{}).LegalCards() {
		if c.ID() == ids.StrikeGoldRed {
			t.Fatalf("LegalCards included Strike Gold [R] despite its NotImplemented tag")
		}
	}
}

// Tests that a known-tagged Unplayable card is absent from LegalCards. Self-retires (skip)
// when the chosen card loses the tag.
func TestLegalCards_ExcludesTaggedUnplayableByID(t *testing.T) {
	if _, ok := GetCard(ids.PotionOfSeeingBlue).(Unplayable); !ok {
		t.Skip("Potion of Seeing [B] is no longer Unplayable — pick another tagged card or drop this test")
	}
	for _, c := range (Registry{}).LegalCards() {
		if c.ID() == ids.PotionOfSeeingBlue {
			t.Fatalf("LegalCards included Potion of Seeing [B] despite its Unplayable tag")
		}
	}
}

// Tests that weapons never leak into LegalCards. Weapons are now full card.Cards sharing
// the CardID pool with deck cards, but they're equipment, not deck cards — they must not be
// reachable by the deck builder / mutations as main-deck entries. Weapon-slot mutations use
// LegalWeapons instead (see TestLegalWeapons_IncludesImplementedWeapon).
func TestLegalCards_ExcludesWeapons(t *testing.T) {
	inPool := make(map[ids.CardID]bool)
	for _, c := range (Registry{}).LegalCards() {
		inPool[c.ID()] = true
	}
	for _, w := range AllWeapons {
		wc, ok := w.(interface{ ID() ids.CardID })
		if !ok {
			t.Fatalf("weapon %s exposes no CardID — the platonic weapon card should be a card.Card", w.Name())
		}
		if inPool[wc.ID()] {
			t.Errorf("LegalCards included weapon %s (CardID %d); weapons must not enter the main-deck pool", w.Name(), wc.ID())
		}
	}
}

// Tests that an implemented, format-legal weapon still appears in LegalWeapons, so weapon-
// slot mutations keep considering it even though weapons are excluded from the card pool.
func TestLegalWeapons_IncludesImplementedWeapon(t *testing.T) {
	legal := make(map[string]bool)
	for _, w := range (Registry{}).LegalWeapons() {
		legal[w.Name()] = true
	}
	// Nebula Blade carries no exclusion marker, so it must be a weapon-slot candidate.
	if !legal["Nebula Blade"] {
		t.Errorf("LegalWeapons missing Nebula Blade; weapon-slot mutations would never consider it")
	}
}

// Tests that LegalWeapons excludes every weapon carrying the NotImplemented marker.
func TestLegalWeapons_SkipsNotImplemented(t *testing.T) {
	for _, w := range (Registry{}).LegalWeapons() {
		if _, ok := w.(NotImplemented); ok {
			t.Errorf("LegalWeapons included NotImplemented weapon %s", w.Name())
		}
	}
}

// Tests that a known-tagged NotImplemented weapon is absent from LegalWeapons. Self-retires
// when the weapon loses the tag.
func TestLegalWeapons_ExcludesTaggedWeaponByID(t *testing.T) {
	tagged := weapons.AnnalsOfSutcliffe{}
	if _, ok := any(tagged).(NotImplemented); !ok {
		t.Skip("Annals of Sutcliffe is no longer NotImplemented — pick another tagged weapon or drop this test")
	}
	for _, w := range (Registry{}).LegalWeapons() {
		if w == tagged {
			t.Fatalf("LegalWeapons included Annals of Sutcliffe despite its NotImplemented tag")
		}
	}
}
