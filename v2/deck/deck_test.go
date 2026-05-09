package deck

import (
	"math/rand"
	"reflect"
	"slices"
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
)

// fakeWeapon is a tiny Weapon implementation for the package's tests. Hand count drives
// loadout legality; name disambiguates dual-wielded slots.
type fakeWeapon struct {
	name  string
	hands int
}

func (w fakeWeapon) Name() string { return w.name }
func (w fakeWeapon) Hands() int   { return w.hands }

// fakeCard is a tiny Card implementation. ID is the per-printing key the copy-budget
// arithmetic operates on; DisplayName is the canonical sideboard-merge key.
type fakeCard struct {
	id      ids.CardID
	display string
}

func (c fakeCard) ID() ids.CardID      { return c.id }
func (c fakeCard) DisplayName() string { return c.display }

// fakeRegistry is the in-memory registry the tests build against. illegal IDs surface
// the SanitizeNotImplemented branch without needing a real exclusion-marker mechanism.
type fakeRegistry struct {
	cards   []Card
	weapons []Weapon
}

func (r *fakeRegistry) LegalCards() []Card     { return slices.Clone(r.cards) }
func (r *fakeRegistry) LegalWeapons() []Weapon { return slices.Clone(r.weapons) }

func newFakeRegistry() *fakeRegistry {
	return &fakeRegistry{
		cards: []Card{
			fakeCard{id: 1, display: "Aether Slash [R]"},
			fakeCard{id: 2, display: "Aether Slash [Y]"},
			fakeCard{id: 3, display: "Read the Runes [R]"},
			fakeCard{id: 4, display: "Hocus Pocus [B]"},
		},
		weapons: []Weapon{fakeWeapon{name: "Reaping Blade", hands: 2}, fakeWeapon{name: "Scepter of Pain", hands: 1}},
	}
}

// cardByID is a test helper for the SanitizeNotImplemented tests that need to retrieve a
// known legal card by its ID; production code doesn't need GetCard on the Registry.
func (r *fakeRegistry) cardByID(id ids.CardID) Card {
	for _, c := range r.cards {
		if c.ID() == id {
			return c
		}
	}
	return nil
}

// Tests that New panics on a 3-weapon loadout (max 2) and on a 2H + 1H loadout (both must
// be 1H when 2 are equipped).
func TestNew_RejectsIllegalWeaponLoadouts(t *testing.T) {
	cases := []struct {
		name    string
		weapons []Weapon
	}{
		{"three weapons", []Weapon{
			fakeWeapon{name: "a", hands: 1},
			fakeWeapon{name: "b", hands: 1},
			fakeWeapon{name: "c", hands: 1},
		}},
		{"2H plus 1H", []Weapon{
			fakeWeapon{name: "two", hands: 2},
			fakeWeapon{name: "one", hands: 1},
		}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			defer func() {
				if recover() == nil {
					t.Errorf("New(%v) didn't panic", tc.weapons)
				}
			}()
			New(nil, tc.weapons, nil)
		})
	}
}

// testDefaults is the loadout the ApplyDefaults tests merge against. The names are made
// up — generic deck logic doesn't care about specific card identity.
var testDefaults = Defaults{
	Equipment: []string{"Helm A", "Boots A", "Robes A"},
	Sideboard: []SideboardDefault{
		{Name: "Side Card", Count: 1},
		{Name: "Read the Runes [R]", Count: 2},
	},
}

// Tests that ApplyDefaults adds each Equipment entry once and is idempotent on the second
// pass.
func TestApplyDefaults_EquipmentIdempotent(t *testing.T) {
	d := New(nil, nil, nil)
	d.ApplyDefaults(testDefaults)
	want := []string{"Helm A", "Boots A", "Robes A"}
	if !reflect.DeepEqual(d.Equipment, want) {
		t.Errorf("after first ApplyDefaults: Equipment = %v, want %v", d.Equipment, want)
	}
	d.ApplyDefaults(testDefaults)
	if !reflect.DeepEqual(d.Equipment, want) {
		t.Errorf("after second ApplyDefaults: Equipment = %v, want %v (should be a no-op)", d.Equipment, want)
	}
}

// Tests that ApplyDefaults adds each Sideboard entry up to its target count and is
// idempotent on the second pass.
func TestApplyDefaults_SideboardIdempotent(t *testing.T) {
	d := New(nil, nil, nil)
	d.ApplyDefaults(testDefaults)
	want := []string{"Side Card", "Read the Runes [R]", "Read the Runes [R]"}
	if !reflect.DeepEqual(d.Sideboard, want) {
		t.Errorf("after first ApplyDefaults: Sideboard = %v, want %v", d.Sideboard, want)
	}
	d.ApplyDefaults(testDefaults)
	if !reflect.DeepEqual(d.Sideboard, want) {
		t.Errorf("after second ApplyDefaults: Sideboard = %v, want %v (should be a no-op)", d.Sideboard, want)
	}
}

// Tests that the sideboard merge respects SideboardCopyCap against main-deck copies —
// when the main deck already holds the per-card cap, the default sideboard add is skipped.
func TestApplyDefaults_SideboardClampsAgainstMainDeck(t *testing.T) {
	red := fakeCard{id: 3, display: "Read the Runes [R]"}
	d := New(nil, nil, []Card{red, red})
	d.ApplyDefaults(testDefaults)

	for _, name := range d.Sideboard {
		if name == "Read the Runes [R]" {
			t.Errorf("Sideboard added Read the Runes [R] despite main deck already holding %d copies (cap is %d)",
				2, SideboardCopyCap)
		}
	}
}

// Tests that Random builds a deck with the requested size, picks weapons from the
// legal pool, and never exceeds maxCopies for any single printing.
func TestRandom_BuildsLegalDeckWithinCopyBudget(t *testing.T) {
	reg := newFakeRegistry()
	rng := rand.New(rand.NewSource(42))

	d := Random(nil, 8, 2, rng, nil, reg)
	if len(d.Cards) != 8 {
		t.Errorf("len(Cards) = %d, want 8", len(d.Cards))
	}
	counts := map[ids.CardID]int{}
	for _, c := range d.Cards {
		counts[c.ID()]++
	}
	for id, n := range counts {
		if n > 2 {
			t.Errorf("card id %d appears %d times, want ≤ 2 (maxCopies)", id, n)
		}
	}
	if len(d.Weapons) == 0 || len(d.Weapons) > 2 {
		t.Errorf("Weapons count = %d, want 1 or 2", len(d.Weapons))
	}
}

// Tests that Random honours the legal filter — when legal rejects every card, Random
// panics rather than returning a partial deck.
func TestRandom_PanicsWhenLegalFilterRejectsEveryCard(t *testing.T) {
	reg := newFakeRegistry()
	defer func() {
		if recover() == nil {
			t.Errorf("Random didn't panic when legal rejected every card")
		}
	}()
	Random(nil, 4, 2, rand.New(rand.NewSource(1)), func(Card) bool { return false }, reg)
}

// Tests that SanitizeNotImplemented swaps every card not in reg.LegalCards for a
// legal replacement, returning a non-empty replacement list.
func TestSanitizeNotImplemented_ReplacesIllegalSlots(t *testing.T) {
	reg := newFakeRegistry()
	bad := fakeCard{id: 99, display: "Mystery Card"}
	good := reg.cardByID(1)
	d := New(nil, nil, []Card{good, bad, good, bad})
	swaps := d.SanitizeNotImplemented(2, rand.New(rand.NewSource(7)), nil, reg)

	if len(swaps) != 2 {
		t.Errorf("swaps count = %d, want 2 (one per illegal slot)", len(swaps))
	}
	legal := map[ids.CardID]bool{}
	for _, c := range reg.LegalCards() {
		legal[c.ID()] = true
	}
	for i, c := range d.Cards {
		if !legal[c.ID()] {
			t.Errorf("Cards[%d] = %v is still illegal after sanitize", i, c)
		}
	}
}

// Tests that SanitizeNotImplemented is a no-op (empty return slice) when every slot is
// already legal.
func TestSanitizeNotImplemented_NoOpOnCleanDeck(t *testing.T) {
	reg := newFakeRegistry()
	d := New(nil, nil, []Card{reg.cardByID(1), reg.cardByID(2)})
	swaps := d.SanitizeNotImplemented(2, rand.New(rand.NewSource(1)), nil, reg)
	if swaps != nil {
		t.Errorf("swaps = %v, want nil for an already-clean deck", swaps)
	}
}

// Tests that SanitizeNotImplemented panics when every pool entry is already saturated in
// the surviving slots — guarding against the previous infinite-loop failure mode.
func TestSanitizeNotImplemented_PanicsWhenPoolSaturated(t *testing.T) {
	// Two-card pool, both at maxCopies=1 in surviving slots; the illegal slot can't
	// be filled without exceeding the cap.
	reg := &fakeRegistry{
		cards: []Card{
			fakeCard{id: 1, display: "A"},
			fakeCard{id: 2, display: "B"},
		},
		weapons: []Weapon{fakeWeapon{name: "w", hands: 2}},
	}
	bad := fakeCard{id: 99, display: "Bad"}
	d := New(nil, nil, []Card{reg.cardByID(1), reg.cardByID(2), bad})

	defer func() {
		if recover() == nil {
			t.Errorf("SanitizeNotImplemented didn't panic when pool was saturated")
		}
	}()
	d.SanitizeNotImplemented(1, rand.New(rand.NewSource(1)), nil, reg)
}
