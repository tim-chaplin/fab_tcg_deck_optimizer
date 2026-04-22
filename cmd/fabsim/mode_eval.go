package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/deck"
)

// runEvalCmd parses eval's flags and dispatches to runEval. eval always operates on a specific
// existing deck, so the deck is a positional arg rather than a -deck flag; only knobs that
// shape the re-simulation (-deep-shuffles, -incoming, -seed) get flags.
func runEvalCmd(args []string) {
	fs := flag.NewFlagSet("eval", flag.ExitOnError)
	fs.Usage = func() {
		fmt.Fprintln(fs.Output(), "Usage: fabsim eval <deck> [flags]")
		fmt.Fprintln(fs.Output())
		fmt.Fprintln(fs.Output(), "Flags:")
		fs.PrintDefaults()
	}
	deepShuffles := fs.Int("deep-shuffles", 10000, "shuffles per deck used to re-score the deck")
	incoming := fs.Int("incoming", 0, "opponent damage per turn")
	seed := fs.Int64("seed", time.Now().UnixNano(), "RNG seed")
	_ = fs.Parse(args)
	if fs.NArg() != 1 {
		die("eval: need exactly one positional <deck> (got %d); try `fabsim eval <deck>`", fs.NArg())
	}
	runEval(resolveDeckPath(fs.Arg(0)), *deepShuffles, *incoming, *seed)
}

// runEval loads the deck at outPath, simulates it for deepShuffles hands, and prints the fresh
// stats. The file on disk is NOT overwritten — eval is a read-only measurement so a deck can be
// re-scored at a new shuffle depth (or a different -incoming) without clobbering the saved
// stats.
func runEval(outPath string, deepShuffles, incoming int, seed int64) {
	loaded, _ := loadExisting(outPath)
	if loaded == nil {
		fmt.Fprintf(os.Stderr, "could not load deck from %s\n", outPath)
		os.Exit(1)
	}
	// Wrap the loaded hero/weapons/cards in a fresh Deck so Evaluate's stats start from zero
	// instead of accumulating on top of the persisted Stats.
	d := deck.New(loaded.Hero, loaded.Weapons, loaded.Cards)
	rng := rand.New(rand.NewSource(seed))
	d.Evaluate(deepShuffles, incoming, rng)
	printBestDeck(d)
}
