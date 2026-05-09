package sim_test

import (
	. "github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
	"math/rand"
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/heroes"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/weapons"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/deck"
)

// TestRandom_FilterExcludesRejected confirms the legal predicate is actually applied to the
// candidate pool: a filter that blocks Plunder Run (all variants) should produce decks that
// never contain any Plunder Run printing, even across many samples.
func TestRandom_FilterExcludesRejected(t *testing.T) {
	bannedIDs := map[ids.CardID]bool{
		ids.PlunderRunRed:    true,
		ids.PlunderRunYellow: true,
		ids.PlunderRunBlue:   true,
	}
	legal := func(c deck.Card) bool { return !bannedIDs[c.ID()] }
	rng := rand.New(rand.NewSource(1))
	for i := 0; i < 20; i++ {
		d := deck.Random(heroes.Viserai{}, 40, 2, rng, legal, registry.Registry{})
		for j, c := range d.Cards {
			if bannedIDs[c.ID()] {
				t.Errorf("sample %d: card[%d] = %s was in the banlist", i, j, c.(Card).Name())
			}
		}
	}
}

// Tests that NotImplemented cards never land in the registry's legal pool, with or
// without a legal predicate.
func TestLegalPool_SkipsNotImplemented(t *testing.T) {
	for _, pred := range []func(deck.Card) bool{nil, func(deck.Card) bool { return true }} {
		for _, c := range (registry.Registry{}).LegalCards() {
			if pred != nil && !pred(c) {
				continue
			}
			if _, ok := c.(NotImplemented); ok {
				t.Errorf("LegalPool included NotImplemented card %s", c.(Card).Name())
			}
		}
	}
}

// Tests that Random never samples a NotImplemented card across many seeds.
func TestRandom_ExcludesNotImplemented(t *testing.T) {
	rng := rand.New(rand.NewSource(1))
	for i := 0; i < 20; i++ {
		d := deck.Random(heroes.Viserai{}, 40, 2, rng, nil, registry.Registry{})
		for j, c := range d.Cards {
			if _, ok := c.(NotImplemented); ok {
				t.Errorf("sample %d card[%d] = %s implements NotImplemented", i, j, c.(Card).Name())
			}
		}
	}
}

// Tests that a known-tagged NotImplemented card (Strike Gold [R]) is absent from LegalPool —
// gives the property test teeth so a broken marker interface fails loudly. Self-retires when
// the card loses the tag.
func TestLegalPool_ExcludesTaggedCardsByID(t *testing.T) {
	if _, ok := GetCard(ids.StrikeGoldRed).(NotImplemented); !ok {
		t.Skip("Strike Gold [R] is no longer NotImplemented — pick another tagged card or drop this test")
	}
	for _, c := range (registry.Registry{}).LegalCards() {
		if c.ID() == ids.StrikeGoldRed {
			t.Fatalf("LegalPool included Strike Gold [R] despite its NotImplemented tag")
		}
	}
}

// Tests that Potion of Seeing [B] (a known-tagged Unplayable card) is absent from the
// registry's legal pool. Self-retires when the card loses the tag.
func TestLegalPool_ExcludesUnplayableByID(t *testing.T) {
	if _, ok := GetCard(ids.PotionOfSeeingBlue).(Unplayable); !ok {
		t.Skip("Potion of Seeing [B] no longer Unplayable; pick another tagged card or drop test")
	}
	for _, c := range (registry.Registry{}).LegalCards() {
		if c.ID() == ids.PotionOfSeeingBlue {
			t.Fatalf("LegalPool included Potion of Seeing [B] despite its Unplayable tag")
		}
	}
}

// Tests SanitizeNotImplemented: replaces tagged slots, preserves deck size, respects
// maxCopies, and reports one swap per replaced card.
func TestSanitizeNotImplemented_ReplacesTaggedSlotsAndKeepsSizeLegal(t *testing.T) {
	if _, ok := GetCard(ids.StrikeGoldRed).(NotImplemented); !ok {
		t.Skip("Strike Gold [R] is no longer NotImplemented — pick another tagged card or drop this test")
	}
	tagged := GetCard(ids.StrikeGoldRed)
	safe := GetCard(ids.ArcanicCrackleRed)
	if _, t2 := safe.(NotImplemented); t2 {
		t.Fatal("ArcanicCrackleRed gained a NotImplemented marker — pick another implemented keeper for this test")
	}
	d := deck.New(heroes.Viserai{}, []deck.Weapon{weapons.NebulaBlade{}},
		[]deck.Card{safe, safe, tagged, tagged})

	rng := rand.New(rand.NewSource(1))
	replaced := d.SanitizeNotImplemented(2, rng, nil, registry.Registry{})

	if len(replaced) != 2 {
		t.Errorf("replaced %d slots, want 2", len(replaced))
	}
	if d.Size() != 4 {
		t.Errorf("card count after sanitize = %d, want 4", d.Size())
	}
	for i, c := range d.Cards {
		if _, ok := c.(NotImplemented); ok {
			t.Errorf("card[%d] = %s still implements NotImplemented", i, c.(Card).Name())
		}
	}
	counts := map[ids.CardID]int{}
	for _, c := range d.Cards {
		counts[c.ID()]++
		if counts[c.ID()] > 2 {
			t.Errorf("%s appears %d times, exceeds maxCopies=2", c.(Card).Name(), counts[c.ID()])
		}
	}
	for _, r := range replaced {
		if r.From.ID() != ids.StrikeGoldRed {
			t.Errorf("replacement From = %s, want Strike Gold [R]", r.From.(Card).Name())
		}
		if _, ok := r.To.(NotImplemented); ok {
			t.Errorf("replacement To = %s implements NotImplemented", r.To.(Card).Name())
		}
	}
}

// TestSanitizeNotImplemented_NoOpOnCleanDeck confirms the sanitizer is an identity operation
// when the deck already has no NotImplemented cards: no replacements, no mutations to
// Cards.
func TestSanitizeNotImplemented_NoOpOnCleanDeck(t *testing.T) {
	a := GetCard(ids.ArcanicCrackleRed)
	if _, tagged := a.(NotImplemented); tagged {
		t.Fatal("ArcanicCrackleRed gained a NotImplemented marker — pick another implemented sentinel")
	}
	d := deck.New(heroes.Viserai{}, []deck.Weapon{weapons.NebulaBlade{}}, []deck.Card{a, a, a, a})
	before := append([]Card(nil), a, a, a, a)

	rng := rand.New(rand.NewSource(1))
	replaced := d.SanitizeNotImplemented(2, rng, nil, registry.Registry{})

	if len(replaced) != 0 {
		t.Errorf("replacements on clean deck = %d, want 0", len(replaced))
	}
	for i, c := range d.Cards {
		if c.ID() != before[i].ID() {
			t.Errorf("card[%d] mutated: %s → %s", i, before[i].Name(), c.(Card).Name())
		}
	}
}

// Tests that no single-slot mutation introduces a NotImplemented card.
func TestAllMutations_ExcludesNotImplementedAdditions(t *testing.T) {
	a := GetCard(ids.ArcanicCrackleRed)
	if _, tagged := a.(NotImplemented); tagged {
		t.Fatal("ArcanicCrackleRed gained a NotImplemented marker — pick another implemented sentinel for this test")
	}
	d := deck.New(heroes.Viserai{}, []deck.Weapon{weapons.NebulaBlade{}}, []deck.Card{a, a, a, a})
	for _, m := range AllMutations(d, 2, registry.Registry{}, nil) {
		for _, c := range m.Deck.Cards {
			if _, ok := c.(NotImplemented); ok {
				t.Errorf("%s introduced NotImplemented card %s", m.Description, c.(Card).Name())
			}
		}
	}
}

// TestLegalWeapons_SkipsNotImplemented pins the weapon-side counterpart to
// TestLegalPool_SkipsNotImplemented: every weapon LegalWeapons surfaces must not implement
// NotImplemented, regardless of which weapons currently carry the tag.
func TestLegalWeapons_SkipsNotImplemented(t *testing.T) {
	for _, w := range (registry.Registry{}).LegalWeapons() {
		if _, ok := w.(NotImplemented); ok {
			t.Errorf("LegalWeapons included NotImplemented weapon %s", w.Name())
		}
	}
}

// TestRandom_ExcludesNotImplementedWeapons confirms no sampled random deck equips a weapon
// tagged with NotImplemented. Mirrors TestRandom_ExcludesNotImplemented for the weapon side.
func TestRandom_ExcludesNotImplementedWeapons(t *testing.T) {
	rng := rand.New(rand.NewSource(1))
	for i := 0; i < 20; i++ {
		d := deck.Random(heroes.Viserai{}, 40, 2, rng, nil, registry.Registry{})
		for j, w := range d.Weapons {
			if _, ok := w.(NotImplemented); ok {
				t.Errorf("sample %d weapon[%d] = %s implements NotImplemented", i, j, w.Name())
			}
		}
	}
}

// Tests that Annals of Sutcliffe (a known-tagged NotImplemented weapon) is absent from
// LegalWeapons. Self-retires when the weapon loses the tag.
func TestLegalWeapons_ExcludesTaggedWeaponByID(t *testing.T) {
	tagged := weapons.AnnalsOfSutcliffe{}
	if _, ok := any(tagged).(NotImplemented); !ok {
		t.Skip("Annals of Sutcliffe is no longer NotImplemented — pick another tagged weapon or drop this test")
	}
	for _, w := range (registry.Registry{}).LegalWeapons() {
		if w.(Weapon).ID() == tagged.ID() {
			t.Fatalf("LegalWeapons included Annals of Sutcliffe despite its NotImplemented tag")
		}
	}
}

// Tests that no weapon-loadout mutation proposes a loadout containing a NotImplemented
// weapon.
func TestAllMutations_ExcludesNotImplementedWeaponLoadouts(t *testing.T) {
	a := GetCard(ids.ArcanicCrackleRed)
	if _, tagged := a.(NotImplemented); tagged {
		t.Fatal("ArcanicCrackleRed gained a NotImplemented marker — pick another implemented sentinel for this test")
	}
	d := deck.New(heroes.Viserai{}, []deck.Weapon{weapons.NebulaBlade{}}, []deck.Card{a, a, a, a})
	for _, m := range AllMutations(d, 2, registry.Registry{}, nil) {
		for _, w := range m.Deck.Weapons {
			if _, ok := w.(NotImplemented); ok {
				t.Errorf("%s introduced NotImplemented weapon %s", m.Description, w.Name())
			}
		}
	}
}

// Tests that banned cards never appear as swap-in candidates, while remaining valid removal
// targets when present in the starting deck.
func TestAllMutations_FilterExcludesRejectedAdditions(t *testing.T) {
	bannedIDs := map[ids.CardID]bool{
		ids.CriticalStrikeRed: true,
	}
	legal := func(c deck.Card) bool { return !bannedIDs[c.ID()] }

	cs := GetCard(ids.CriticalStrikeRed)
	other := GetCard(ids.AetherSlashRed)
	d := deck.New(heroes.Viserai{}, []deck.Weapon{weapons.NebulaBlade{}}, []deck.Card{cs, cs, other, other})

	for i, m := range AllMutations(d, 2, registry.Registry{}, legal) {
		bannedIn := 0
		for _, c := range m.Deck.Cards {
			if bannedIDs[c.ID()] {
				bannedIn++
			}
		}
		// The starting deck has 2 copies of Critical Strike [R]. A mutation that removes one
		// leaves 1; a mutation that removes the other leaves 1; a weapon-only mutation leaves
		// all 2. No mutation should ADD another copy.
		if bannedIn > 2 {
			t.Errorf("mutation %d (%s): has %d banned copies, want <=2 (no additions allowed)",
				i, m.Description, bannedIn)
		}
	}
}
