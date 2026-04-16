package deckio

import (
	"math/rand"
	"reflect"
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/deck"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/hero"
)

func TestMarshalUnmarshalRoundTrip(t *testing.T) {
	rng := rand.New(rand.NewSource(1))
	d := deck.Random(hero.Viserai{}, 40, 2, rng)
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
	if got.Stats.Best.Play.Value != d.Stats.Best.Play.Value {
		t.Errorf("best value: got %d want %d", got.Stats.Best.Play.Value, d.Stats.Best.Play.Value)
	}
}
