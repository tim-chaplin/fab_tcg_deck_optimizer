package main

import (
	"math/rand"
	"path/filepath"
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/deck"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/format"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/gameengine"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/hero/heroes"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

// TestPrepareBaseline_AlwaysRescoresLoadedDeck pins that a loaded deck is re-scored against the
// current simulator even when its saved run count already exceeds the budget, so a stale saved avg
// can't seed bestEver. The seed claims 50000 runs; prepareBaseline must return fresh stats at the
// current (smaller) budget rather than trust them.
func TestPrepareBaseline_AlwaysRescoresLoadedDeck(t *testing.T) {
	gameengine.MaxDeckSize = 40
	path := filepath.Join(t.TempDir(), "deck.json")

	rng := rand.New(rand.NewSource(1))
	d := deck.Random(heroes.Viserai, format.SilverAge, 40, 2, rng, registry.Registry{})
	staleStats := deck.Stats{Runs: 50000, Hands: 100, TotalValue: 999} // saved avg 9.99
	if err := writeDeck(d, staleStats, path); err != nil {
		t.Fatalf("seed writeDeck: %v", err)
	}

	cfg := annealConfig{
		shuffles:  50,
		matchup:   sim.Matchup{IncomingPhysicalDamage: 5},
		deckSize:  40,
		maxCopies: 2,
		outPath:   path,
		quietLoad: true,
	}
	var fresh deck.Stats
	captureEvalOutput(t, func() {
		_, fresh, _ = prepareBaseline(cfg, rng)
	})
	if fresh.Runs != cfg.shuffles {
		t.Errorf("loaded deck not re-scored: returned stats.Runs=%d, want %d (saved claimed 50000)",
			fresh.Runs, cfg.shuffles)
	}
}
