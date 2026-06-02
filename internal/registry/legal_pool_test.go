package registry

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/deck"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/format"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/hero/heroes"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/weapon/weapons"
)

// Tests that LegalCardsFor's hero filter drops from Viserai's pool exactly the cards Viserai
// can't play — the Lightning-talent cards (Aurora's, like Amulet of Lightning). Every other
// implemented card is Runeblade / Generic with no talent, so it stays.
func TestLegalCardsFor_Viserai_ExcludesLightningCards(t *testing.T) {
	inHeroPool := make(map[ids.CardID]bool)
	for _, c := range (Registry{}).LegalCardsFor(format.SilverAge, heroes.Viserai) {
		inHeroPool[c.ID()] = true
	}
	excluded := 0
	for _, c := range legalCardsForFormat(format.SilverAge) {
		// Viserai is Runeblade with no talent, so a Lightning-talent card is the one thing it
		// can't play; everything else (Runeblade / Generic, no talent) it can.
		viseraiCanPlay := !c.Types(nil).Has(card.TypeLightning)
		if inHeroPool[c.ID()] != viseraiCanPlay {
			t.Errorf("%s in Viserai pool = %v, want %v", c.DisplayName(), inHeroPool[c.ID()], viseraiCanPlay)
		}
		if !viseraiCanPlay {
			excluded++
		}
	}
	if excluded == 0 {
		t.Error("no Lightning cards excluded — the test no longer exercises the hero filter")
	}
}

// Tests that LegalCards excludes every card carrying the NotImplemented marker. The deck
// builder relies on this filter rather than re-checking each printing per call to Random.
func TestLegalCards_SkipsNotImplemented(t *testing.T) {
	for _, c := range legalCardsForFormat(format.SilverAge) {
		if _, ok := c.(NotImplemented); ok {
			t.Errorf("legal-card pool included NotImplemented card %s", c.DisplayName())
		}
	}
}

// Tests that LegalCards excludes every card carrying the Unplayable marker.
func TestLegalCards_SkipsUnplayable(t *testing.T) {
	for _, c := range legalCardsForFormat(format.SilverAge) {
		if _, ok := c.(Unplayable); ok {
			t.Errorf("legal-card pool included Unplayable card %s", c.DisplayName())
		}
	}
}

// Tests that a known-tagged NotImplemented card is absent from LegalCards. Self-retires
// (skip) when the chosen card loses the tag.
func TestLegalCards_ExcludesTaggedNotImplementedByID(t *testing.T) {
	if _, ok := GetCard(ids.StrikeGoldRed).(NotImplemented); !ok {
		t.Skip("Strike Gold [R] is no longer NotImplemented — pick another tagged card or drop this test")
	}
	for _, c := range legalCardsForFormat(format.SilverAge) {
		if c.ID() == ids.StrikeGoldRed {
			t.Fatalf("legal-card pool included Strike Gold [R] despite its NotImplemented tag")
		}
	}
}

// Tests that a known-tagged Unplayable card is absent from LegalCards. Self-retires (skip)
// when the chosen card loses the tag.
func TestLegalCards_ExcludesTaggedUnplayableByID(t *testing.T) {
	if _, ok := GetCard(ids.PotionOfSeeingBlue).(Unplayable); !ok {
		t.Skip("Potion of Seeing [B] is no longer Unplayable — pick another tagged card or drop this test")
	}
	for _, c := range legalCardsForFormat(format.SilverAge) {
		if c.ID() == ids.PotionOfSeeingBlue {
			t.Fatalf("legal-card pool included Potion of Seeing [B] despite its Unplayable tag")
		}
	}
}

// Tests that the implemented Silver Age banned cards stay out of the pool. These four are the
// cards whose only pool-exclusion is the format banlist (no NotImplemented / Unplayable
// marker), so they're the regression guard for legality now flowing through internal/format
// rather than a per-card NotSilverAgeLegal marker. Exercises the real card Name() against the
// banlist, so a banlist apostrophe / spelling drift (e.g. Fiddler's Green) fails here.
func TestLegalCards_ExcludesFormatBanned(t *testing.T) {
	banned := map[string]bool{
		"Sirens of Safe Harbor": true,
		"Plunder Run":           true,
		"Fiddler's Green":       true,
		"Fate Foreseen":         true,
	}
	for _, c := range legalCardsForFormat(format.SilverAge) {
		if banned[c.Name()] {
			t.Errorf("legal-card pool included Silver Age banned card %s", c.DisplayName())
		}
	}
}

// Tests that weapons never leak into LegalCards. Weapons are now full card.Cards sharing
// the CardID pool with deck cards, but they're equipment, not deck cards — they must not be
// reachable by the deck builder / mutations as main-deck entries. Weapon-slot mutations use
// LegalWeapons instead (see TestLegalWeapons_IncludesImplementedWeapon).
func TestLegalCards_ExcludesWeapons(t *testing.T) {
	inPool := make(map[ids.CardID]bool)
	for _, c := range legalCardsForFormat(format.SilverAge) {
		inPool[c.ID()] = true
	}
	for _, w := range AllWeapons {
		wc, ok := w.(interface{ ID() ids.CardID })
		if !ok {
			t.Fatalf("weapon %s exposes no CardID — the platonic weapon card should be a card.Card", w.Name())
		}
		if inPool[wc.ID()] {
			t.Errorf("legal-card pool included weapon %s (CardID %d); weapons must not enter the main-deck pool", w.Name(), wc.ID())
		}
	}
}

// Tests that IllegalCards flags banlisted and wrong-class cards while passing cards legal for
// the hero.
func TestIllegalCards(t *testing.T) {
	rune := fakeHero{class: card.TypeRuneblade, typeSet: card.NewTypeSet(card.TypeRuneblade)}
	legal := GetCard(ids.AetherSlashRed)
	if bad := IllegalCards(format.SilverAge, rune, []deck.Card{legal}); len(bad) != 0 {
		t.Errorf("Aether Slash flagged illegal for a Runeblade: %v", bad)
	}

	id, ok := CardByName("Sink Below")
	if !ok {
		t.Skip("Sink Below not registered — pick another implemented banned card")
	}
	if bad := IllegalCards(format.SilverAge, rune, []deck.Card{GetCard(id)}); len(bad) != 1 {
		t.Errorf("a banned card was not flagged illegal, got %d illegal", len(bad))
	}

	thief := fakeHero{class: card.TypeThief, typeSet: card.NewTypeSet(card.TypeThief)}
	if bad := IllegalCards(format.SilverAge, thief, []deck.Card{legal}); len(bad) != 1 {
		t.Errorf("a Runeblade card was not flagged illegal for a non-Runeblade hero, got %d", len(bad))
	}
}

// fakeRarityCard exposes just the Rarity() probe that rarityLegal asserts on.
type fakeRarityCard struct{ rarity string }

func (c fakeRarityCard) Rarity() string { return c.rarity }

// Tests the format rarity gate: cards carrying Rarity() are filtered by the rule, and values
// without it defer to the other filters.
func TestRarityLegal_GatesOnPrintedRarity(t *testing.T) {
	if rarityLegal(format.SilverAge, fakeRarityCard{"Majestic"}) {
		t.Error("rarityLegal kept a Majestic card; the Silver Age rarity rule should drop it")
	}
	if !rarityLegal(format.SilverAge, fakeRarityCard{"Rare"}) {
		t.Error("rarityLegal dropped a Rare card; Basic/Common/Rare are legal")
	}
	if !rarityLegal(format.SilverAge, struct{}{}) {
		t.Error("rarityLegal dropped a value with no Rarity() method; it should defer to the other filters")
	}
}

// Tests that every card in the pool declares a Basic/Common/Rare rarity — the invariant the
// format rarity rule enforces.
func TestLegalCards_AllRarityLegal(t *testing.T) {
	for _, c := range legalCardsForFormat(format.SilverAge) {
		r, ok := c.(interface{ Rarity() string })
		if !ok {
			t.Errorf("pool card %s exposes no Rarity()", c.DisplayName())
			continue
		}
		if !format.SilverAge.IsRarityLegal(r.Rarity()) {
			t.Errorf("pool card %s has non-Silver-Age rarity %q", c.DisplayName(), r.Rarity())
		}
	}
}

// Tests that an implemented, format-legal weapon still appears in LegalWeapons, so weapon-
// slot mutations keep considering it even though weapons are excluded from the card pool.
func TestLegalWeapons_IncludesImplementedWeapon(t *testing.T) {
	legal := make(map[string]bool)
	for _, w := range (Registry{}).LegalWeaponsFor(format.SilverAge, heroes.Viserai) {
		legal[w.Name()] = true
	}
	// Nebula Blade carries no exclusion marker, so it must be a weapon-slot candidate.
	if !legal["Nebula Blade"] {
		t.Errorf("legal-weapon pool missing Nebula Blade; weapon-slot mutations would never consider it")
	}
}

// Tests that LegalWeapons excludes every weapon carrying the NotImplemented marker.
func TestLegalWeapons_SkipsNotImplemented(t *testing.T) {
	for _, w := range (Registry{}).LegalWeaponsFor(format.SilverAge, heroes.Viserai) {
		if _, ok := w.(NotImplemented); ok {
			t.Errorf("legal-weapon pool included NotImplemented weapon %s", w.Name())
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
	for _, w := range (Registry{}).LegalWeaponsFor(format.SilverAge, heroes.Viserai) {
		if w == tagged {
			t.Fatalf("legal-weapon pool included Annals of Sutcliffe despite its NotImplemented tag")
		}
	}
}
