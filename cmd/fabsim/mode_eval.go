package main

import (
	"fmt"
	"math/rand"
	"os"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/deck"
)

// runEval loads the deck at cfg.outPath, simulates it for cfg.deepShuffles hands, and prints
// the fresh stats. The file on disk is NOT overwritten — eval is a read-only measurement so a
// deck can be re-scored at a new shuffle depth (or a different -incoming) without clobbering
// the saved stats.
func runEval(cfg config) {
	loaded, _ := loadExisting(cfg.outPath)
	if loaded == nil {
		fmt.Fprintf(os.Stderr, "could not load deck from %s\n", cfg.outPath)
		os.Exit(1)
	}
	// Wrap the loaded hero/weapons/cards in a fresh Deck so Evaluate's stats start from zero
	// instead of accumulating on top of the persisted Stats.
	d := deck.New(loaded.Hero, loaded.Weapons, loaded.Cards)
	rng := rand.New(rand.NewSource(cfg.seed))
	d.Evaluate(cfg.deepShuffles, cfg.incoming, rng)
	printBestDeck(d)
}
