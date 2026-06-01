package sim

import (
	"math/rand"
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/deck"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/hero/heroes"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry"
)

// TestPerCardAblation_DeltaPerUniqueCard checks PerCardAblation returns one delta per unique card
// and is deterministic for a fixed single-worker seed.
func TestPerCardAblation_DeltaPerUniqueCard(t *testing.T) {
	read := registry.GetCard(ids.ReadTheRunesRed)
	snatch := registry.GetCard(ids.SnatchRed)
	d := deck.New(heroes.Viserai, nil, []deck.Card{read, read, read, read, snatch, snatch, snatch, snatch})
	mp := Matchup{IncomingPhysicalDamage: 4}

	ev := NewEvaluator()
	deltas := ev.PerCardAblation(d, 40, mp, rand.New(rand.NewSource(1)))

	uniqueIDs, _ := d.UniqueIDs()
	if len(deltas) != len(uniqueIDs) {
		t.Fatalf("got %d deltas, want one per unique card (%d)", len(deltas), len(uniqueIDs))
	}
	for _, id := range uniqueIDs {
		if _, ok := deltas[id]; !ok {
			t.Errorf("missing delta for %s", registry.GetCard(id).Name())
		}
	}

	again := ev.PerCardAblation(d, 40, mp, rand.New(rand.NewSource(1)))
	for id, v := range deltas {
		if again[id] != v {
			t.Errorf("non-deterministic delta for %s: %v vs %v", registry.GetCard(id).Name(), v, again[id])
		}
	}
}
