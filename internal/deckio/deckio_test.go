package deckio

import (
	"math/rand"
	"reflect"
	"testing"

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

// TestRoundTrip_PreservesBestTurnContributions pins the fix for "loaded deck's best-turn display
// shows +0 everywhere": the JSON form now serialises per-card Contribution alongside Hand/Roles
// and an ordered AttackChain with per-step Damage / TriggerDamage, so re-loading a deck gives
// FormatBestTurn the same per-card numbers the live sim produced.
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
	}
}
