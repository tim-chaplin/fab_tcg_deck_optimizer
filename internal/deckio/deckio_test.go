package deckio

import (
	"math/rand"
	"reflect"
	"strings"
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/deck"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/hero"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/registry"
)

func TestMarshalUnmarshalRoundTrip(t *testing.T) {
	rng := rand.New(rand.NewSource(1))
	d := deck.Random(hero.Viserai{}, 40, 2, rng, nil, registry.Registry{})
	stats := sim.NewEvaluator().Evaluate(d, 50, sim.Matchup{IncomingDamage: 4}, rng)

	data, err := Marshal(d, stats)
	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}
	got, gotStats, err := Unmarshal(data)
	if err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}

	if got.Hero.(hero.Hero).Name() != d.Hero.(hero.Hero).Name() {
		t.Errorf("hero: got %q want %q", got.Hero.(hero.Hero).Name(), d.Hero.(hero.Hero).Name())
	}
	if got.Size() != d.Size() {
		t.Fatalf("cards len: got %d want %d", got.Size(), d.Size())
	}
	// Compare as multisets — the JSON form sorts card names, so order isn't preserved across a
	// round trip. What matters is that the same cards (with the same counts) come back.
	wantCounts := d.NameCounts()
	gotCounts := got.NameCounts()
	if !reflect.DeepEqual(gotCounts, wantCounts) {
		t.Errorf("card counts: got %v want %v", gotCounts, wantCounts)
	}
	if !reflect.DeepEqual(gotStats.FirstCycle, stats.FirstCycle) {
		t.Errorf("first cycle: got %+v want %+v", gotStats.FirstCycle, stats.FirstCycle)
	}
	if gotStats.Best.Value != stats.Best.Value {
		t.Errorf("best value: got %d want %d", gotStats.Best.Value, stats.Best.Value)
	}
}

// Tests that PerCardMarginal (Present/Absent totals and hand counts per CardID) survives
// Marshal/Unmarshal.
func TestRoundTrip_PreservesPerCardMarginal(t *testing.T) {
	rng := rand.New(rand.NewSource(13))
	d := deck.Random(hero.Viserai{}, 40, 2, rng, nil, registry.Registry{})
	stats := sim.NewEvaluator().Evaluate(d, 50, sim.Matchup{IncomingDamage: 4}, rng)
	if len(stats.PerCardMarginal) == 0 {
		t.Fatalf("baseline deck produced no PerCardMarginal entries; test can't differentiate good from bad")
	}

	data, err := Marshal(d, stats)
	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}
	_, gotStats, err := Unmarshal(data)
	if err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}

	if len(gotStats.PerCardMarginal) != len(stats.PerCardMarginal) {
		t.Fatalf("PerCardMarginal entry count: got %d want %d",
			len(gotStats.PerCardMarginal), len(stats.PerCardMarginal))
	}
	for id, want := range stats.PerCardMarginal {
		gotEntry, ok := gotStats.PerCardMarginal[id]
		if !ok {
			t.Errorf("PerCardMarginal missing entry for %s after round trip", registry.GetCard(id).Name())
			continue
		}
		if gotEntry.PresentHands != want.PresentHands || gotEntry.AbsentHands != want.AbsentHands {
			t.Errorf("%s bucket counts: got present=%d absent=%d, want present=%d absent=%d",
				registry.GetCard(id).Name(),
				gotEntry.PresentHands, gotEntry.AbsentHands,
				want.PresentHands, want.AbsentHands)
		}
		if gotEntry.PresentTotal != want.PresentTotal || gotEntry.AbsentTotal != want.AbsentTotal {
			t.Errorf("%s bucket totals: got present=%v absent=%v, want present=%v absent=%v",
				registry.GetCard(id).Name(),
				gotEntry.PresentTotal, gotEntry.AbsentTotal,
				want.PresentTotal, want.AbsentTotal)
		}
		if gotEntry.Marginal() != want.Marginal() {
			t.Errorf("%s Marginal(): got %v want %v", registry.GetCard(id).Name(),
				gotEntry.Marginal(), want.Marginal())
		}
	}
}

// Tests that deck.BestTurn.Log round-trips through Marshal/Unmarshal verbatim.
func TestRoundTrip_PreservesBestTurnLog(t *testing.T) {
	rng := rand.New(rand.NewSource(7))
	d := deck.Random(hero.Viserai{}, 40, 2, rng, nil, registry.Registry{})
	want := deck.TurnLog{
		StartOfTurn: []string{
			"Hand: Hocus Pocus [B], Consuming Volition [R]",
			"Arsenal: Sigil of the Arknight [B]",
			"Auras: 1 Runechant",
			"Sigil of the Arknight [B]: drew Hit the High Notes [R] into hand",
		},
		MyTurn: []string{
			"Hocus Pocus [B]: PITCH",
			"Consuming Volition [R]: ATTACK (+4)",
			"Viserai created a runechant (+1)",
		},
		EndOfTurn: []string{
			"Hand: Hit the High Notes [R]",
			"Auras: 1 Runechant",
		},
	}
	stats := deck.Stats{
		Best: deck.BestTurn{
			Value: 21,
			Log:   want,
		},
	}

	data, err := Marshal(d, stats)
	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}
	_, gotStats, err := Unmarshal(data)
	if err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}

	if !reflect.DeepEqual(gotStats.Best.Log, want) {
		t.Errorf("Log: got %+v\n want %+v", gotStats.Best.Log, want)
	}
	if gotStats.Best.Value != 21 {
		t.Errorf("Value: got %d want 21", gotStats.Best.Value)
	}
}

// Tests that the user-managed Sideboard field round-trips through Marshal/Unmarshal as a
// multiset of names.
func TestRoundTrip_PreservesSideboard(t *testing.T) {
	rng := rand.New(rand.NewSource(42))
	d := deck.Random(hero.Viserai{}, 40, 2, rng, nil, registry.Registry{})
	d.Sideboard = []string{"Aether Slash [R]", "Aether Slash [R]", "Arcanic Spike [B]"}

	data, err := Marshal(d, deck.Stats{})
	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}
	got, _, err := Unmarshal(data)
	if err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}

	want := map[string]int{"Aether Slash [R]": 2, "Arcanic Spike [B]": 1}
	gotCounts := map[string]int{}
	for _, name := range got.Sideboard {
		gotCounts[name]++
	}
	if !reflect.DeepEqual(gotCounts, want) {
		t.Errorf("sideboard counts: got %v want %v", gotCounts, want)
	}
}

// TestMarshal_OmitsEmptySideboard pins the omitempty behaviour: a deck with no sideboard
// doesn't emit the field at all, so existing on-disk JSON files stay byte-identical after
// a re-serialize.
func TestMarshal_OmitsEmptySideboard(t *testing.T) {
	rng := rand.New(rand.NewSource(1))
	d := deck.Random(hero.Viserai{}, 40, 2, rng, nil, registry.Registry{})
	data, err := Marshal(d, deck.Stats{})
	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}
	if strings.Contains(string(data), `"sideboard"`) {
		t.Errorf("empty sideboard should be omitted from JSON; got:\n%s", data)
	}
}

// Tests that Sideboard accepts any name verbatim — the sim never reads it, so unknown
// names don't abort Unmarshal.
func TestUnmarshal_SideboardAcceptsAnyName(t *testing.T) {
	const raw = `{
  "hero": "Viserai",
  "weapons": [],
  "cards": [],
  "sideboard": ["Not A Real Card", "Crown of Dichotomy"],
  "pitch": {"red": 0, "yellow": 0, "blue": 0},
  "stats": {"avg": 0, "runs": 0, "hands": 0, "total_value": 0, "first_cycle": {"Hands": 0, "Total": 0}, "second_cycle": {"Hands": 0, "Total": 0}, "best": {"hand": [], "roles": [], "weapons": [], "value": 0, "starting_runechants": 0}}
}`
	d, _, err := Unmarshal([]byte(raw))
	if err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}
	want := []string{"Not A Real Card", "Crown of Dichotomy"}
	if !reflect.DeepEqual(d.Sideboard, want) {
		t.Errorf("sideboard: got %v want %v", d.Sideboard, want)
	}
}

// TestRoundTrip_PreservesEquipment pins the Equipment round-trip: the field survives
// Marshal/Unmarshal and accepts names the card registry doesn't cover (the sim doesn't
// model equipment pieces).
func TestRoundTrip_PreservesEquipment(t *testing.T) {
	rng := rand.New(rand.NewSource(1))
	d := deck.Random(hero.Viserai{}, 40, 2, rng, nil, registry.Registry{})
	d.Equipment = []string{"Beckoning Haunt", "Nullrune Boots", "Blade Beckoner Helm"}

	data, err := Marshal(d, deck.Stats{})
	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}
	got, _, err := Unmarshal(data)
	if err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}
	// Sorted copy — Marshal sorts for stable on-disk diff.
	want := []string{"Beckoning Haunt", "Blade Beckoner Helm", "Nullrune Boots"}
	if !reflect.DeepEqual(got.Equipment, want) {
		t.Errorf("equipment: got %v want %v", got.Equipment, want)
	}
}

// TestMarshal_OmitsEmptyEquipment pins omitempty: a deck with no equipment doesn't emit
// the field at all, keeping existing files byte-identical after a re-serialize.
func TestMarshal_OmitsEmptyEquipment(t *testing.T) {
	rng := rand.New(rand.NewSource(1))
	d := deck.Random(hero.Viserai{}, 40, 2, rng, nil, registry.Registry{})
	data, err := Marshal(d, deck.Stats{})
	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}
	if strings.Contains(string(data), `"equipment"`) {
		t.Errorf("empty equipment should be omitted from JSON; got:\n%s", data)
	}
}
