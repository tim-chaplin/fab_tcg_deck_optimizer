package fabrary

import (
	"math/rand"
	"reflect"
	"strings"
	"testing"

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

// TestMarshalIncludesDefaultArenaPackage pins the convenience behaviour: Marshal emits the
// user's fixed equipment loadout (defaultArenaPackage) alongside the deck's weapons, so the
// exported .txt can be pasted into fabrary without hand-editing equipment slots.
func TestMarshalIncludesDefaultArenaPackage(t *testing.T) {
	rng := rand.New(rand.NewSource(1))
	d := deck.Random(hero.Viserai{}, 40, 2, rng, nil)
	text := Marshal(d)

	for _, name := range defaultArenaPackage {
		want := "1x " + name + "\n"
		if !strings.Contains(text, want) {
			t.Errorf("expected %q in Arena section; output:\n%s", strings.TrimSuffix(want, "\n"), text)
		}
	}
}

// TestUnmarshalSample parses the exact sample the user supplied (verbatim from fabrary.net's
// plain-text export) to prove the parser tolerates the real output, including the footer and the
// mix of weapons + non-weapon equipment in the Arena section.
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
	// Exactly one weapon in the sample maps to a registered weapon ("Reaping Blade"); the other
	// Arena lines are equipment the optimizer doesn't model and should be silently skipped.
	if len(d.Weapons) != 1 || d.Weapons[0].Name() != "Reaping Blade" {
		names := make([]string, len(d.Weapons))
		for i, w := range d.Weapons {
			names[i] = w.Name()
		}
		t.Errorf("weapons: got %v want [Reaping Blade]", names)
	}
	if len(d.Cards) == 0 {
		t.Fatalf("expected deck cards, got none")
	}
	// "Arcane Polarity" isn't in the registry yet; the sample has 2 red copies, so the skip map
	// should report it. If/when it gets implemented this expectation needs to move to whichever
	// card remains unimplemented in the sample.
	if skipped["Arcane Polarity (Red)"] != 2 {
		t.Errorf("skipped map should report Arcane Polarity (Red) x2; got %v", skipped)
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
	if skipped["Not A Real Card (Red)"] != 2 {
		t.Errorf("skipped should contain Not A Real Card (Red) x2; got %v", skipped)
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
		m[c.Name()]++
	}
	return m
}

func weaponNameCounts(d *deck.Deck) map[string]int {
	m := map[string]int{}
	for _, w := range d.Weapons {
		m[w.Name()]++
	}
	return m
}

