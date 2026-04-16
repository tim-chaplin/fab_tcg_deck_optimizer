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
	for i := range d.Cards {
		if got.Cards[i].Name() != d.Cards[i].Name() {
			t.Errorf("card[%d]: got %q want %q", i, got.Cards[i].Name(), d.Cards[i].Name())
		}
	}
	if !reflect.DeepEqual(got.Stats.FirstCycle, d.Stats.FirstCycle) {
		t.Errorf("first cycle: got %+v want %+v", got.Stats.FirstCycle, d.Stats.FirstCycle)
	}
	if got.Stats.Best.Play.Value != d.Stats.Best.Play.Value {
		t.Errorf("best value: got %d want %d", got.Stats.Best.Play.Value, d.Stats.Best.Play.Value)
	}
}
