package fabrary

import (
	"math/rand"
	"reflect"
	"strings"
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/deck"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/hero"
)

// TestMarshalUnmarshalRoundTrip exercises a random deck through Marshal → Unmarshal and checks that
// weapons, cards, and hero all come back intact (stats are intentionally not round-tripped).
func TestMarshalUnmarshalRoundTrip(t *testing.T) {
	rng := rand.New(rand.NewSource(1))
	d := deck.Random(hero.Viserai{}, 40, 2, rng, nil)

	text := Marshal(d)
	got, skipped, err := Unmarshal(text)
	if err != nil {
		t.Fatalf("Unmarshal: %v\n---\n%s", err, text)
	}
	if len(skipped) != 0 {
		t.Errorf("unexpected skipped cards on registered-only round trip: %v", skipped)
	}

	if got.Hero.Name() != d.Hero.Name() {
		t.Errorf("hero: got %q want %q", got.Hero.Name(), d.Hero.Name())
	}
	gotW, wantW := weaponNameCounts(got), weaponNameCounts(d)
	if !reflect.DeepEqual(gotW, wantW) {
		t.Errorf("weapon counts: got %v want %v", gotW, wantW)
	}
	wantCards := cardNameCounts(d)
	gotCards := cardNameCounts(got)
	if !reflect.DeepEqual(gotCards, wantCards) {
		t.Errorf("card counts: got %v want %v", gotCards, wantCards)
	}
}

// TestMarshalFormat pins the output shape: header, Arena section, Deck section, lowercase color
// suffix, and sorted card lines. A change here means downstream fabrary compatibility may break —
// update consciously.
func TestMarshalFormat(t *testing.T) {
	rng := rand.New(rand.NewSource(1))
	d := deck.Random(hero.Viserai{}, 40, 2, rng, nil)
	text := Marshal(d)

	wantPrefix := "Name: Viserai\nHero: Viserai\nFormat: Silver Age\n\nArena cards\n"
	if !strings.HasPrefix(text, wantPrefix) {
		t.Errorf("header mismatch:\ngot:\n%s\nwant prefix:\n%s", text, wantPrefix)
	}
	if !strings.Contains(text, "\n\nDeck cards\n") {
		t.Errorf("missing Deck cards section:\n%s", text)
	}
	// Pitch color should be lowercased in parens.
	if strings.Contains(text, "(Red)") || strings.Contains(text, "(Yellow)") || strings.Contains(text, "(Blue)") {
		t.Errorf("expected lowercase pitch colors; got:\n%s", text)
	}
}

// TestMarshalRendersAppliedDefaults pins that when a caller runs ApplyDefaults (as
// writeDeck does), Marshal's output carries the default equipment in Arena and the default
// sideboard entries in Sideboard. The default lists themselves are pinned in the deck
// package's tests; this one only checks the fabrary text picks them up verbatim.
//
// Builds a Deck directly instead of going through deck.Random so the test is stable across
// pool changes (e.g. cards getting tagged NotImplemented and dropping out of the random pool,
// causing a seed-1 deck to roll different cards that collide with sideboard-default cap).
func TestMarshalRendersAppliedDefaults(t *testing.T) {
	d := &deck.Deck{Hero: hero.Viserai{}}
	d.ApplyDefaults()
	text := Marshal(d)

	// Pick representative entries from each default list — full-list coverage belongs in
	// deck.TestApplyDefaults_*.
	for _, want := range []string{
		"1x Beckoning Haunt\n",
		"1x Blade Beckoner Helm\n",
		"1x Crown of Dichotomy\n",
		"2x Read the Runes (red)\n",
	} {
		if !strings.Contains(text, want) {
			t.Errorf("expected %q in output; got:\n%s", strings.TrimSuffix(want, "\n"), text)
		}
	}
}

// TestUnmarshalSample parses the exact sample the user supplied (verbatim from fabrary.net's
// plain-text export) to prove the parser tolerates the real output, including the footer and the
// mix of weapons + non-weapon equipment in the Arena section. Non-weapon arena entries land
// in Deck.Equipment so they round-trip on re-export.
func TestUnmarshalSample(t *testing.T) {
	const sample = `Name: Viserai
Hero: Viserai
Format: Silver Age

Arena cards
1x Blade Beckoner Boots
1x Blade Beckoner Helm
1x Blossom of Spring
1x Crown of Dichotomy
1x Nullrune Boots
1x Nullrune Gloves
1x Reaping Blade
1x Runebleed Robe
1x Runehold Release

Deck cards
2x Arcane Polarity (red)
2x Condemn to Slaughter (red)
2x Deathly Duet (red)
2x Drowning Dire (red)
2x Malefic Incantation (red)
2x Mauvrion Skies (red)
2x Meat and Greet (red)
2x Reduce to Runechant (red)
2x Rune Flash (red)
2x Runeblood Incantation (red)
2x Runerager Swarm (red)
2x Runic Reaping (red)
2x Shrill of Skullform (red)
2x Sigil of Suffering (red)
2x Spellblade Assault (red)
2x Weeping Battleground (red)
2x Malefic Incantation (yellow)
2x Deathly Duet (blue)
2x Malefic Incantation (blue)
2x Meat and Greet (blue)
2x Rune Flash (blue)
2x Runerager Swarm (blue)
2x Shrill of Skullform (blue)

Made with ❤️ at the FaBrary
See the full deck @ https://fabrary.net/decks/01KP1AZ5SAS425YN30WB779M41
`
	d, skipped, err := Unmarshal(sample)
	if err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}
	if d.Hero.Name() != "Viserai" {
		t.Errorf("hero: got %q want %q", d.Hero.Name(), "Viserai")
	}
	// Exactly one weapon in the sample maps to a registered weapon ("Reaping Blade"); the
	// other non-weapon Arena lines are equipment the optimizer doesn't model — they land in
	// Deck.Equipment as opaque strings.
	if len(d.Weapons) != 1 || d.Weapons[0].Name() != "Reaping Blade" {
		names := make([]string, len(d.Weapons))
		for i, w := range d.Weapons {
			names[i] = w.Name()
		}
		t.Errorf("weapons: got %v want [Reaping Blade]", names)
	}
	// Count how many of each equipment name landed in Deck.Equipment, and spot-check a
	// couple of the ones in the sample.
	gotEquipment := map[string]int{}
	for _, e := range d.Equipment {
		gotEquipment[e]++
	}
	for _, want := range []string{"Blade Beckoner Boots", "Crown of Dichotomy", "Nullrune Boots", "Runehold Release"} {
		if gotEquipment[want] != 1 {
			t.Errorf("equipment[%q] = %d, want 1", want, gotEquipment[want])
		}
	}
	if len(d.Cards) == 0 {
		t.Fatalf("expected deck cards, got none")
	}
	// Every card in this sample is now in the registry (some as NotImplemented stubs); the
	// unknown-card skip path is covered by TestUnmarshalUnknownCardSkipped instead.
	if len(skipped) != 0 {
		t.Errorf("expected no skipped cards on this sample; got %v", skipped)
	}
}

// TestUnmarshalUnknownCardSkipped pins the lenient behaviour: unknown deck-section cards do NOT
// abort the parse; they're reported in the returned skip map so the caller can warn. fabrary
// decks routinely reference cards the optimizer hasn't implemented yet, so a strict failure would
// make import unusable in practice.
func TestUnmarshalUnknownCardSkipped(t *testing.T) {
	const text = `Name: Viserai
Hero: Viserai
Format: Silver Age

Arena cards
1x Reaping Blade

Deck cards
2x Not A Real Card (red)
`
	d, skipped, err := Unmarshal(text)
	if err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}
	if len(d.Cards) != 0 {
		t.Errorf("expected 0 known cards, got %d", len(d.Cards))
	}
	if skipped["Not A Real Card [R]"] != 2 {
		t.Errorf("skipped should contain Not A Real Card [R] x2; got %v", skipped)
	}
}

// TestMarshalSideboardSection verifies that a deck with an explicit Sideboard renders a
// trailing "Sideboard" section with count-and-name lines in the same shape as the Deck section,
// placed after Deck cards.
func TestMarshalSideboardSection(t *testing.T) {
	rng := rand.New(rand.NewSource(1))
	d := deck.Random(hero.Viserai{}, 40, 2, rng, nil)

	// Use Mauvrion Skies [R] — its pitch-color suffix exercises the toFabraryCardName
	// lowercase conversion. Sideboard is a string list; names are stored in canonical form.
	d.Sideboard = []string{"Mauvrion Skies [R]", "Mauvrion Skies [R]"}
	text := Marshal(d)
	if !strings.Contains(text, "\nSideboard\n") {
		t.Errorf("populated sideboard should emit a Sideboard section; got:\n%s", text)
	}
	if !strings.Contains(text, "2x Mauvrion Skies (red)") {
		t.Errorf("expected '2x Mauvrion Skies (red)' in sideboard section; got:\n%s", text)
	}
	if strings.Index(text, "Sideboard") < strings.Index(text, "Deck cards") {
		t.Errorf("Sideboard must come after Deck cards; got:\n%s", text)
	}
}

// TestUnmarshalSideboardRoundTrip pins the import path: a fabrary-style text with a Sideboard
// section parses into Deck.Sideboard as a multiset, separate from the main card list.
func TestUnmarshalSideboardRoundTrip(t *testing.T) {
	const sample = `Name: Viserai
Hero: Viserai
Format: Silver Age

Arena cards
1x Reaping Blade

Deck cards
2x Aether Slash (red)

Sideboard
2x Mauvrion Skies (red)
1x Runic Reaping (blue)
`
	d, skipped, err := Unmarshal(sample)
	if err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}
	if len(skipped) != 0 {
		t.Errorf("unexpected skipped cards: %v", skipped)
	}
	wantMain := map[string]int{"Aether Slash [R]": 2}
	wantSide := map[string]int{"Mauvrion Skies [R]": 2, "Runic Reaping [B]": 1}
	if got := cardNameCounts(d); !reflect.DeepEqual(got, wantMain) {
		t.Errorf("main cards: got %v want %v", got, wantMain)
	}
	gotSide := map[string]int{}
	for _, name := range d.Sideboard {
		gotSide[name]++
	}
	if !reflect.DeepEqual(gotSide, wantSide) {
		t.Errorf("sideboard: got %v want %v", gotSide, wantSide)
	}
}

// TestUnmarshalUnknownHeroFails guards the one remaining hard failure: without a known hero, we
// can't build a deck at all, so this must still be an error (not a silent drop).
func TestUnmarshalUnknownHeroFails(t *testing.T) {
	const text = `Name: Someone
Hero: Not A Hero
Format: Silver Age

Arena cards

Deck cards
`
	_, _, err := Unmarshal(text)
	if err == nil {
		t.Fatal("expected error for unknown hero, got nil")
	}
}

func cardNameCounts(d *deck.Deck) map[string]int {
	m := map[string]int{}
	for _, c := range d.Cards {
		m[card.DisplayName(c)]++
	}
	return m
}

func weaponNameCounts(d *deck.Deck) map[string]int {
	m := map[string]int{}
	for _, w := range d.Weapons {
		m[card.DisplayName(w)]++
	}
	return m
}
