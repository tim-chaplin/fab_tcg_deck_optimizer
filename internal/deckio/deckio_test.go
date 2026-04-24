package deckio

import (
	"math/rand"
	"reflect"
	"strings"
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/cards"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/deck"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/hand"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/hero"
)

func filterHand(line []hand.CardAssignment) []hand.CardAssignment {
	out := make([]hand.CardAssignment, 0, len(line))
	for _, a := range line {
		if a.FromArsenal {
			continue
		}
		out = append(out, a)
	}
	return out
}

func TestMarshalUnmarshalRoundTrip(t *testing.T) {
	rng := rand.New(rand.NewSource(1))
	d := deck.Random(hero.Viserai{}, 40, 2, rng, nil)
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

// TestRoundTrip_PreservesBestTurnContributions locks in that per-card Contribution and the
// AttackChain's per-step Damage / TriggerDamage / AuraTriggerDamage round-trip through
// Marshal/Unmarshal, so a reloaded deck renders with the same per-card numbers the live sim
// produced.
func TestRoundTrip_PreservesBestTurnContributions(t *testing.T) {
	rng := rand.New(rand.NewSource(7))
	d := deck.Random(hero.Viserai{}, 40, 2, rng, nil)
	d.Evaluate(100, 0, rng)

	if len(d.Stats.Best.Summary.BestLine) == 0 {
		t.Skip("evaluation produced no best turn; rerun with a different seed")
	}

	// Locate at least one BestLine entry and one AttackChain entry with non-zero damage/contrib
	// so the assertion below is meaningful — if the sim found nothing of value we'd be checking
	// 0 == 0 and the test wouldn't catch a regression that dropped the fields.
	var haveNonZeroContrib, haveNonZeroDamage bool
	for _, a := range d.Stats.Best.Summary.BestLine {
		if a.Contribution != 0 {
			haveNonZeroContrib = true
		}
	}
	for _, e := range d.Stats.Best.Summary.AttackChain {
		if e.Damage != 0 || e.TriggerDamage != 0 {
			haveNonZeroDamage = true
		}
	}
	if !haveNonZeroContrib || !haveNonZeroDamage {
		t.Skip("evaluation produced an all-zero best turn; rerun with a different seed")
	}

	data, err := Marshal(d)
	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}
	got, err := Unmarshal(data)
	if err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}

	// Arsenal-in entries (FromArsenal=true) belong to a previous turn's hand and are skipped by
	// bestTurnToJSON; compare hand-only entries in parallel order.
	wantLine := filterHand(d.Stats.Best.Summary.BestLine)
	gotLine := filterHand(got.Stats.Best.Summary.BestLine)
	if len(gotLine) != len(wantLine) {
		t.Fatalf("BestLine hand len: got %d want %d", len(gotLine), len(wantLine))
	}
	for i := range wantLine {
		if gotLine[i].Contribution != wantLine[i].Contribution {
			t.Errorf("BestLine[%d].Contribution: got %.3f want %.3f", i, gotLine[i].Contribution, wantLine[i].Contribution)
		}
	}

	wantChain := d.Stats.Best.Summary.AttackChain
	gotChain := got.Stats.Best.Summary.AttackChain
	if len(gotChain) != len(wantChain) {
		t.Fatalf("AttackChain len: got %d want %d", len(gotChain), len(wantChain))
	}
	for i := range wantChain {
		if gotChain[i].Card.Name() != wantChain[i].Card.Name() {
			t.Errorf("AttackChain[%d].Card: got %q want %q", i, gotChain[i].Card.Name(), wantChain[i].Card.Name())
		}
		if gotChain[i].Damage != wantChain[i].Damage {
			t.Errorf("AttackChain[%d].Damage: got %.3f want %.3f", i, gotChain[i].Damage, wantChain[i].Damage)
		}
		if gotChain[i].TriggerDamage != wantChain[i].TriggerDamage {
			t.Errorf("AttackChain[%d].TriggerDamage: got %.3f want %.3f", i, gotChain[i].TriggerDamage, wantChain[i].TriggerDamage)
		}
		if gotChain[i].AuraTriggerDamage != wantChain[i].AuraTriggerDamage {
			t.Errorf("AttackChain[%d].AuraTriggerDamage: got %.3f want %.3f", i, gotChain[i].AuraTriggerDamage, wantChain[i].AuraTriggerDamage)
		}
	}
}

// TestRoundTrip_PreservesStartOfTurnAuras locks in that the best turn's StartOfTurnAuras list
// (the auras that were in play at the top of the captured turn) survives Marshal/Unmarshal by
// card name, preserving duplicates and order. Without the round-trip, a reloaded deck's best
// turn would lose its "Auras in play at start of turn" header line.
func TestRoundTrip_PreservesStartOfTurnAuras(t *testing.T) {
	rng := rand.New(rand.NewSource(7))
	d := deck.Random(hero.Viserai{}, 40, 2, rng, nil)
	// Seed a best turn by hand so the assertion doesn't depend on the sim organically
	// producing carryover auras — synthetic input is enough to pin the round trip.
	d.Stats.Best = deck.BestTurn{
		Summary: hand.TurnSummary{
			BestLine: []hand.CardAssignment{{Card: cards.Get(card.MaleficIncantationRed), Role: hand.Attack}},
			StartOfTurnAuras: []card.Card{
				cards.Get(card.MaleficIncantationRed),
				cards.Get(card.MaleficIncantationRed),
				cards.Get(card.SigilOfTheArknightBlue),
			},
			Value: 1,
		},
	}

	data, err := Marshal(d)
	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}
	got, err := Unmarshal(data)
	if err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}

	wantNames := []string{"Malefic Incantation (Red)", "Malefic Incantation (Red)", "Sigil of the Arknight (Blue)"}
	gotAuras := got.Stats.Best.Summary.StartOfTurnAuras
	if len(gotAuras) != len(wantNames) {
		t.Fatalf("StartOfTurnAuras len: got %d want %d", len(gotAuras), len(wantNames))
	}
	for i := range wantNames {
		if gotAuras[i].Name() != wantNames[i] {
			t.Errorf("StartOfTurnAuras[%d]: got %q want %q", i, gotAuras[i].Name(), wantNames[i])
		}
	}
}

// TestRoundTrip_PreservesArsenalIn pins the arsenal-in entry round-trip: a BestLine slot
// with FromArsenal=true survives Marshal/Unmarshal so the reloaded deck can re-render the
// "(from arsenal)" tag. Without this, `fabsim eval -print-only` would lose the tag on any
// saved best turn whose winning play chain included the arsenal-in card.
func TestRoundTrip_PreservesArsenalIn(t *testing.T) {
	rng := rand.New(rand.NewSource(11))
	d := deck.Random(hero.Viserai{}, 40, 2, rng, nil)
	arsenalCard := cards.Get(card.MauvrionSkiesRed)
	handCard := cards.Get(card.HitTheHighNotesRed)
	// Seed a best turn by hand with an arsenal-in entry at the tail (bestUncached's
	// convention: hand cards at indices [0,n); arsenal-in at n). The chain includes both
	// so rebuildAttackChain reconstructs them on load.
	d.Stats.Best = deck.BestTurn{
		Summary: hand.TurnSummary{
			BestLine: []hand.CardAssignment{
				{Card: handCard, Role: hand.Attack, Contribution: 6},
				{Card: arsenalCard, Role: hand.Attack, Contribution: 3, FromArsenal: true},
			},
			AttackChain: []hand.AttackChainEntry{
				{Card: arsenalCard, Damage: 3},
				{Card: handCard, Damage: 6},
			},
			Value: 9,
		},
	}

	data, err := Marshal(d)
	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}
	got, err := Unmarshal(data)
	if err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}

	var arsenalEntry hand.CardAssignment
	var found bool
	for _, a := range got.Stats.Best.Summary.BestLine {
		if a.FromArsenal {
			arsenalEntry = a
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("BestLine missing arsenal-in entry after round-trip; got %+v", got.Stats.Best.Summary.BestLine)
	}
	if arsenalEntry.Card.ID() != arsenalCard.ID() {
		t.Errorf("arsenal-in card: got %q want %q", arsenalEntry.Card.Name(), arsenalCard.Name())
	}
	if arsenalEntry.Role != hand.Attack {
		t.Errorf("arsenal-in role: got %v want Attack", arsenalEntry.Role)
	}
	if arsenalEntry.Contribution != 3 {
		t.Errorf("arsenal-in contribution: got %v want 3", arsenalEntry.Contribution)
	}
}

// TestRoundTrip_PreservesTriggersFromLastTurn pins the carryover-AuraTrigger round-trip
// including the Revealed-into-hand attribution. Sigil of the Arknight fires at start of
// action phase with Damage=0 and reveals the deck top; a reloaded deck must still render
// the "drew X into hand" line, which requires both the aura and its revealed card to
// round-trip by name.
func TestRoundTrip_PreservesTriggersFromLastTurn(t *testing.T) {
	rng := rand.New(rand.NewSource(13))
	d := deck.Random(hero.Viserai{}, 40, 2, rng, nil)
	d.Stats.Best = deck.BestTurn{
		Summary: hand.TurnSummary{
			BestLine: []hand.CardAssignment{{Card: cards.Get(card.MaleficIncantationRed), Role: hand.Attack}},
			TriggersFromLastTurn: []hand.TriggerContribution{
				{Card: cards.Get(card.SigilOfTheArknightBlue), Revealed: cards.Get(card.HitTheHighNotesRed)},
				{Card: cards.Get(card.MaleficIncantationRed), Damage: 2},
			},
			Value: 3,
		},
	}

	data, err := Marshal(d)
	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}
	got, err := Unmarshal(data)
	if err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}

	trigs := got.Stats.Best.Summary.TriggersFromLastTurn
	if len(trigs) != 2 {
		t.Fatalf("TriggersFromLastTurn len: got %d want 2 (%+v)", len(trigs), trigs)
	}
	if trigs[0].Card.Name() != "Sigil of the Arknight (Blue)" {
		t.Errorf("trigs[0].Card: got %q want Sigil of the Arknight (Blue)", trigs[0].Card.Name())
	}
	if trigs[0].Revealed == nil || trigs[0].Revealed.Name() != "Hit the High Notes (Red)" {
		t.Errorf("trigs[0].Revealed: got %v want Hit the High Notes (Red)", trigs[0].Revealed)
	}
	if trigs[0].Damage != 0 {
		t.Errorf("trigs[0].Damage: got %d want 0", trigs[0].Damage)
	}
	if trigs[1].Card.Name() != "Malefic Incantation (Red)" || trigs[1].Damage != 2 || trigs[1].Revealed != nil {
		t.Errorf("trigs[1]: got %+v want {Malefic Incantation (Red), Damage:2, Revealed:nil}", trigs[1])
	}
}

// TestRoundTrip_PreservesSideboard verifies the user-managed Sideboard field survives a
// Marshal/Unmarshal cycle as a multiset of card names. Sideboard contents don't affect the
// sim — this test pins that the IO layer still round-trips them so `fabsim eval` / `anneal`
// can preserve them across runs even when they'd otherwise drop the data.
func TestRoundTrip_PreservesSideboard(t *testing.T) {
	rng := rand.New(rand.NewSource(42))
	d := deck.Random(hero.Viserai{}, 40, 2, rng, nil)
	d.Sideboard = []string{"Aether Slash (Red)", "Aether Slash (Red)", "Arcanic Spike (Blue)"}

	data, err := Marshal(d)
	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}
	got, err := Unmarshal(data)
	if err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}

	want := map[string]int{"Aether Slash (Red)": 2, "Arcanic Spike (Blue)": 1}
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
	d := deck.Random(hero.Viserai{}, 40, 2, rng, nil)
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
	d := deck.Random(hero.Viserai{}, 40, 2, rng, nil)
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
	d := deck.Random(hero.Viserai{}, 40, 2, rng, nil)
	data, err := Marshal(d)
	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}
	if strings.Contains(string(data), `"equipment"`) {
		t.Errorf("empty equipment should be omitted from JSON; got:\n%s", data)
	}
}
