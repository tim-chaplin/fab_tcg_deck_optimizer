package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/deck"
	fmtpkg "github.com/tim-chaplin/fab-deck-optimizer/internal/format"
)

// runEvalCmd parses eval's flags and dispatches to runEval. eval always operates on an
// existing deck passed positionally; flags cover only re-simulation knobs.
func runEvalCmd(args []string) {
	fs := flag.NewFlagSet("eval", flag.ExitOnError)
	fs.Usage = func() {
		fmt.Fprintln(fs.Output(), "Usage: fabsim eval <deck> [flags]")
		fmt.Fprintln(fs.Output())
		fmt.Fprintln(fs.Output(), "Flags:")
		fs.PrintDefaults()
	}
	deepShuffles := fs.Int("deep-shuffles", 10000, "shuffles per deck used to re-score the deck")
	incoming := fs.Int("incoming", 0, "opponent damage per turn (required unless -print-only is set — must match the value the deck was annealed at for comparable numbers)")
	seed := fs.Int64("seed", time.Now().UnixNano(), "RNG seed")
	formatFlag := fs.String("format", string(fmtpkg.SilverAge), "constructed format predicate applied to replacement picks when the loaded deck contains NotImplemented cards")
	maxCopies := fs.Int("max-copies", defaultMaxCopies, "maximum copies of any single card printing per deck, applied when replacing NotImplemented cards in the loaded deck")
	printOnly := fs.Bool("print-only", false, "load the deck and print the stats from the last run without simulating or rewriting the on-disk .json / .txt")
	brief := fs.Bool("brief", false, "print only the score summary (no card list, per-card stats, or best turn)")
	_ = parseFlagsAnywhere(fs, args)
	if fs.NArg() != 1 {
		die("eval: need exactly one positional <deck> (got %d); try `fabsim eval <deck>`", fs.NArg())
	}
	if !*printOnly {
		requireFlag(fs, "eval", "incoming")
	}
	fmtValue, err := fmtpkg.Parse(*formatFlag)
	if err != nil {
		die("%v", err)
	}
	runEval(resolveDeckPath(fs.Arg(0)), *deepShuffles, *incoming, *maxCopies, *seed, fmtValue, *printOnly, *brief)
}

// runEval loads the deck at outPath and prints its stats. Default behaviour (printOnly=false)
// re-simulates the deck for deepShuffles hands against incoming and writes the fresh stats
// back to disk — both the JSON and the sibling fabrary .txt — so the on-disk copy always
// reflects the latest binary's modelling. printOnly=true skips the simulation and the rewrite
// entirely: the loaded stats are printed as-is, which is what you want for a quick look at a
// saved deck without spending shuffles or mutating the file.
//
// The sanitize pass (replacing any card.NotImplemented copies with legal substitutes) runs
// only in the simulate path: a sanitized deck has to be rewritten anyway so the on-disk avg
// stays in sync with the cards we can actually simulate, and skipping it in print-only
// preserves the "don't touch the file" promise.
//
// Output shape is controlled by brief:
//   - brief=false (default): full printBestDeck dump — summary, card list, best-turn block,
//     per-card stats.
//   - brief=true: score summary only. Good for scripted re-scoring where the card list and
//     best turn are noise.
func runEval(outPath string, deepShuffles, incoming, maxCopies int, seed int64, fmtValue fmtpkg.Format, printOnly, brief bool) {
	loaded := mustLoadDeck(outPath)
	if printOnly {
		printLoadedDeck(loaded, brief)
		return
	}
	// Wrap the loaded hero/weapons/cards in a fresh Deck so Evaluate's stats start from zero
	// instead of accumulating on top of the persisted Stats. Sideboard and Equipment carry
	// over verbatim — the sim ignores both, but the post-eval writeDeck round-trips them
	// back to disk so the user's hand-managed lists aren't dropped by a re-score.
	d := deck.New(loaded.Hero, loaded.Weapons, loaded.Cards)
	d.Sideboard = loaded.Sideboard
	d.Equipment = loaded.Equipment
	rng := rand.New(rand.NewSource(seed))
	savedAvg := loaded.Stats.Mean()
	sanitizeLoadedDeck(d, maxCopies, rng, fmtValue.IsLegal)
	d.Evaluate(deepShuffles, incoming, rng)
	fmt.Fprintf(os.Stderr, "eval: avg %.3f → %.3f (delta %+.3f); rewriting %s\n",
		savedAvg, d.Stats.Mean(), d.Stats.Mean()-savedAvg, outPath)
	if err := writeDeck(d, outPath); err != nil {
		die("%v", err)
	}
	printLoadedDeck(d, brief)
}

// printLoadedDeck dispatches between the brief summary and the full printBestDeck dump;
// used by both the simulate path and -print-only.
func printLoadedDeck(d *deck.Deck, brief bool) {
	if brief {
		printDeckSummary(d)
		return
	}
	printBestDeck(d)
}
