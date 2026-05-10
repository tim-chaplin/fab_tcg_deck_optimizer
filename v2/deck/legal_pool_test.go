package deck_test

import (
	"math/rand"
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/heroes"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
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
		uniq, _ := d.UniqueIDs()
		for _, id := range uniq {
			if bannedIDs[id] {
				t.Errorf("sample %d: %s was in the banlist", i, registry.GetCard(id).DisplayName())
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
			if _, ok := c.(sim.NotImplemented); ok {
				t.Errorf("LegalPool included NotImplemented card %s", c.(sim.Card).Name())
			}
		}
	}
}

// Tests that Random never samples a NotImplemented card across many seeds.
func TestRandom_ExcludesNotImplemented(t *testing.T) {
	rng := rand.New(rand.NewSource(1))
	for i := 0; i < 20; i++ {
		d := deck.Random(heroes.Viserai{}, 40, 2, rng, nil, registry.Registry{})
		uniq, _ := d.UniqueIDs()
		for _, id := range uniq {
			if _, ok := registry.GetCard(id).(sim.NotImplemented); ok {
				t.Errorf("sample %d: %s implements NotImplemented", i, registry.GetCard(id).DisplayName())
			}
		}
	}
}

// Tests that a known-tagged NotImplemented card (Strike Gold [R]) is absent from LegalPool —
// gives the property test teeth so a broken marker interface fails loudly. Self-retires when
// the card loses the tag.
func TestLegalPool_ExcludesTaggedCardsByID(t *testing.T) {
	if _, ok := registry.GetCard(ids.StrikeGoldRed).(sim.NotImplemented); !ok {
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
	if _, ok := registry.GetCard(ids.PotionOfSeeingBlue).(sim.Unplayable); !ok {
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
	if _, ok := registry.GetCard(ids.StrikeGoldRed).(sim.NotImplemented); !ok {
		t.Skip("Strike Gold [R] is no longer NotImplemented — pick another tagged card or drop this test")
	}
	tagged := registry.GetCard(ids.StrikeGoldRed)
	safe := registry.GetCard(ids.ArcanicCrackleRed)
	if _, t2 := safe.(sim.NotImplemented); t2 {
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
	uniq, _ := d.UniqueIDs()
	for _, id := range uniq {
		if _, ok := registry.GetCard(id).(sim.NotImplemented); ok {
			t.Errorf("%s still implements NotImplemented", registry.GetCard(id).DisplayName())
		}
	}
	for name, count := range d.NameCounts() {
		if count > 2 {
			t.Errorf("%s appears %d times, exceeds maxCopies=2", name, count)
		}
	}
	for _, r := range replaced {
		if r.From.ID() != ids.StrikeGoldRed {
			t.Errorf("replacement From = %s, want Strike Gold [R]", r.From.(sim.Card).Name())
		}
		if _, ok := r.To.(sim.NotImplemented); ok {
			t.Errorf("replacement To = %s implements NotImplemented", r.To.(sim.Card).Name())
		}
	}
}

// TestSanitizeNotImplemented_NoOpOnCleanDeck confirms the sanitizer is an identity operation
// when the deck already has no NotImplemented cards: no replacements, no composition change.
func TestSanitizeNotImplemented_NoOpOnCleanDeck(t *testing.T) {
	a := registry.GetCard(ids.ArcanicCrackleRed)
	if _, tagged := a.(sim.NotImplemented); tagged {
		t.Fatal("ArcanicCrackleRed gained a NotImplemented marker — pick another implemented sentinel")
	}
	d := deck.New(heroes.Viserai{}, []deck.Weapon{weapons.NebulaBlade{}}, []deck.Card{a, a, a, a})
	beforeFingerprint := d.Fingerprint()

	rng := rand.New(rand.NewSource(1))
	replaced := d.SanitizeNotImplemented(2, rng, nil, registry.Registry{})

	if len(replaced) != 0 {
		t.Errorf("replacements on clean deck = %d, want 0", len(replaced))
	}
	if got := d.Fingerprint(); got != beforeFingerprint {
		t.Errorf("composition mutated: fingerprint %q → %q", beforeFingerprint, got)
	}
}

// Tests that no single-slot mutation introduces a NotImplemented card.
func TestAllMutations_ExcludesNotImplementedAdditions(t *testing.T) {
	a := registry.GetCard(ids.ArcanicCrackleRed)
	if _, tagged := a.(sim.NotImplemented); tagged {
		t.Fatal("ArcanicCrackleRed gained a NotImplemented marker — pick another implemented sentinel for this test")
	}
	d := deck.New(heroes.Viserai{}, []deck.Weapon{weapons.NebulaBlade{}}, []deck.Card{a, a, a, a})
	for _, m := range deck.AllMutations(d, 2, registry.Registry{}, nil) {
		uniq, _ := m.Deck.UniqueIDs()
		for _, id := range uniq {
			if _, ok := registry.GetCard(id).(sim.NotImplemented); ok {
				t.Errorf("%s introduced NotImplemented card %s", m.Description, registry.GetCard(id).DisplayName())
			}
		}
	}
}

// TestLegalWeapons_SkipsNotImplemented pins the weapon-side counterpart to
// TestLegalPool_SkipsNotImplemented: every weapon LegalWeapons surfaces must not implement
// NotImplemented, regardless of which weapons currently carry the tag.
func TestLegalWeapons_SkipsNotImplemented(t *testing.T) {
	for _, w := range (registry.Registry{}).LegalWeapons() {
		if _, ok := w.(sim.NotImplemented); ok {
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
			if _, ok := w.(sim.NotImplemented); ok {
				t.Errorf("sample %d weapon[%d] = %s implements NotImplemented", i, j, w.Name())
			}
		}
	}
}

// Tests that Annals of Sutcliffe (a known-tagged NotImplemented weapon) is absent from
// LegalWeapons. Self-retires when the weapon loses the tag.
func TestLegalWeapons_ExcludesTaggedWeaponByID(t *testing.T) {
	tagged := weapons.AnnalsOfSutcliffe{}
	if _, ok := any(tagged).(sim.NotImplemented); !ok {
		t.Skip("Annals of Sutcliffe is no longer NotImplemented — pick another tagged weapon or drop this test")
	}
	for _, w := range (registry.Registry{}).LegalWeapons() {
		if w.(sim.Weapon).ID() == tagged.ID() {
			t.Fatalf("LegalWeapons included Annals of Sutcliffe despite its NotImplemented tag")
		}
	}
}

// Tests that no weapon-loadout mutation proposes a loadout containing a NotImplemented
// weapon.
func TestAllMutations_ExcludesNotImplementedWeaponLoadouts(t *testing.T) {
	a := registry.GetCard(ids.ArcanicCrackleRed)
	if _, tagged := a.(sim.NotImplemented); tagged {
		t.Fatal("ArcanicCrackleRed gained a NotImplemented marker — pick another implemented sentinel for this test")
	}
	d := deck.New(heroes.Viserai{}, []deck.Weapon{weapons.NebulaBlade{}}, []deck.Card{a, a, a, a})
	for _, m := range deck.AllMutations(d, 2, registry.Registry{}, nil) {
		for _, w := range m.Deck.Weapons {
			if _, ok := w.(sim.NotImplemented); ok {
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

	cs := registry.GetCard(ids.CriticalStrikeRed)
	other := registry.GetCard(ids.AetherSlashRed)
	d := deck.New(heroes.Viserai{}, []deck.Weapon{weapons.NebulaBlade{}}, []deck.Card{cs, cs, other, other})

	for i, m := range deck.AllMutations(d, 2, registry.Registry{}, legal) {
		bannedIn := 0
		counts := m.Deck.NameCounts()
		for id := range bannedIDs {
			bannedIn += counts[registry.GetCard(id).DisplayName()]
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
