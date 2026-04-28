package deckio

import (
	"math/rand"
	"reflect"
	"strings"
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/heroes"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

func TestMarshalUnmarshalRoundTrip(t *testing.T) {
	rng := rand.New(rand.NewSource(1))
	d := sim.Random(heroes.Viserai{}, 40, 2, rng, nil)
	d.Evaluate(50, 4, rng)

	data, err := Marshal(d)
	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}
	got, err := Unmarshal(data)
	if err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}

	if got.Hero.Name() != d.Hero.Name() {
		t.Errorf("hero: got %q want %q", got.Hero.Name(), d.Hero.Name())
	}
	if len(got.Cards) != len(d.Cards) {
		t.Fatalf("cards len: got %d want %d", len(got.Cards), len(d.Cards))
	}
	// Compare as multisets — the JSON form sorts card names, so order isn't preserved across a
	// round trip. What matters is that the same cards (with the same counts) come back.
	wantCounts := map[string]int{}
	for _, c := range d.Cards {
		wantCounts[c.Name()]++
	}
	gotCounts := map[string]int{}
	for _, c := range got.Cards {
		gotCounts[c.Name()]++
	}
	if !reflect.DeepEqual(gotCounts, wantCounts) {
		t.Errorf("card counts: got %v want %v", gotCounts, wantCounts)
	}
	if !reflect.DeepEqual(got.Stats.FirstCycle, d.Stats.FirstCycle) {
		t.Errorf("first cycle: got %+v want %+v", got.Stats.FirstCycle, d.Stats.FirstCycle)
	}
	if got.Stats.Best.Summary.Value != d.Stats.Best.Summary.Value {
		t.Errorf("best value: got %d want %d", got.Stats.Best.Summary.Value, d.Stats.Best.Summary.Value)
	}
}

// TestRoundTrip_PreservesPerCardMarginal pins that the per-card marginal-stats accumulator
// (PresentTotal/PresentHands and AbsentTotal/AbsentHands per unique ids.CardID) survives a
// Marshal/Unmarshal so a re-loaded deck can render the marginal-value table without a
// fresh sim. Compared via the public Marginal() so a regression in any of the four
// underlying fields surfaces.
func TestRoundTrip_PreservesPerCardMarginal(t *testing.T) {
	rng := rand.New(rand.NewSource(13))
	d := sim.Random(heroes.Viserai{}, 40, 2, rng, nil)
	d.Evaluate(50, 4, rng)
	if len(d.Stats.PerCardMarginal) == 0 {
		t.Fatalf("baseline deck produced no PerCardMarginal entries; test can't differentiate good from bad")
	}

	data, err := Marshal(d)
	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}
	got, err := Unmarshal(data)
	if err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}

	if len(got.Stats.PerCardMarginal) != len(d.Stats.PerCardMarginal) {
		t.Fatalf("PerCardMarginal entry count: got %d want %d",
			len(got.Stats.PerCardMarginal), len(d.Stats.PerCardMarginal))
	}
	for id, want := range d.Stats.PerCardMarginal {
		gotEntry, ok := got.Stats.PerCardMarginal[id]
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

// TestRoundTrip_PreservesBestTurnLog pins the on-disk best-turn round-trip: Marshal/Unmarshal
// carry sim.BestTurn.Log verbatim, so a reloaded TurnLog matches section-for-section. Since
// the formatter consumes Log at print time, pinning the structured shape implicitly pins the
// printout shape.
func TestRoundTrip_PreservesBestTurnLog(t *testing.T) {
	rng := rand.New(rand.NewSource(7))
	d := sim.Random(heroes.Viserai{}, 40, 2, rng, nil)
	want := sim.TurnLog{
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
	d.Stats.Best = sim.BestTurn{
		Summary:            sim.TurnSummary{Value: 21},
		StartingRunechants: 0,
		Log:                want,
	}

	data, err := Marshal(d)
	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}
	got, err := Unmarshal(data)
	if err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}

	if !reflect.DeepEqual(got.Stats.Best.Log, want) {
		t.Errorf("Log: got %+v\n want %+v", got.Stats.Best.Log, want)
	}
	if got.Stats.Best.Summary.Value != 21 {
		t.Errorf("Value: got %d want 21", got.Stats.Best.Summary.Value)
	}
}

// TestRoundTrip_PreservesSideboard verifies the user-managed Sideboard field survives a
// Marshal/Unmarshal cycle as a multiset of card names. Sideboard contents don't affect the
// sim — this test pins that the IO layer still round-trips them so `fabsim eval` / `anneal`
// can preserve them across runs even when they'd otherwise drop the data.
func TestRoundTrip_PreservesSideboard(t *testing.T) {
	rng := rand.New(rand.NewSource(42))
	d := sim.Random(heroes.Viserai{}, 40, 2, rng, nil)
	d.Sideboard = []string{"Aether Slash [R]", "Aether Slash [R]", "Arcanic Spike [B]"}

	data, err := Marshal(d)
	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}
	got, err := Unmarshal(data)
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
	d := sim.Random(heroes.Viserai{}, 40, 2, rng, nil)
	data, err := Marshal(d)
	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}
	if strings.Contains(string(data), `"sideboard"`) {
		t.Errorf("empty sideboard should be omitted from JSON; got:\n%s", data)
	}
}

// TestUnmarshal_SideboardAcceptsAnyName pins the lenient behaviour of Sideboard: it's a
// name-only list the sim never reads, so any string the user wrote — including equipment
// pieces and yet-to-be-implemented cards — round-trips verbatim instead of aborting the
// parse.
func TestUnmarshal_SideboardAcceptsAnyName(t *testing.T) {
	const raw = `{
  "hero": "Viserai",
  "weapons": [],
  "cards": [],
  "sideboard": ["Not A Real Card", "Crown of Dichotomy"],
  "pitch": {"red": 0, "yellow": 0, "blue": 0},
  "stats": {"avg": 0, "runs": 0, "hands": 0, "total_value": 0, "first_cycle": {"Hands": 0, "Total": 0}, "second_cycle": {"Hands": 0, "Total": 0}, "best": {"hand": [], "roles": [], "weapons": [], "value": 0, "starting_runechants": 0}}
}`
	d, err := Unmarshal([]byte(raw))
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
	d := sim.Random(heroes.Viserai{}, 40, 2, rng, nil)
	d.Equipment = []string{"Beckoning Haunt", "Nullrune Boots", "Blade Beckoner Helm"}

	data, err := Marshal(d)
	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}
	got, err := Unmarshal(data)
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
	d := sim.Random(heroes.Viserai{}, 40, 2, rng, nil)
	data, err := Marshal(d)
	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}
	if strings.Contains(string(data), `"equipment"`) {
		t.Errorf("empty equipment should be omitted from JSON; got:\n%s", data)
	}
}
