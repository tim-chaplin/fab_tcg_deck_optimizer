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
	incoming := fs.Int("incoming", 0, "opponent damage per turn (required — must match the value the deck was annealed at for comparable numbers)")
	seed := fs.Int64("seed", time.Now().UnixNano(), "RNG seed")
	formatFlag := fs.String("format", string(fmtpkg.SilverAge), "constructed format predicate applied to replacement picks when the loaded deck contains NotImplemented cards")
	maxCopies := fs.Int("max-copies", defaultMaxCopies, "maximum copies of any single card printing per deck, applied when replacing NotImplemented cards in the loaded deck")
	reevaluate := fs.Bool("reevaluate", false, "overwrite the on-disk .json / .txt with this run's fresh stats")
	brief := fs.Bool("brief", false, "print only the score summary (no card list, per-card stats, or best turn)")
	_ = parseFlagsAnywhere(fs, args)
	if fs.NArg() != 1 {
		die("eval: need exactly one positional <deck> (got %d); try `fabsim eval <deck>`", fs.NArg())
	}
	requireFlag(fs, "eval", "incoming")
	fmtValue, err := fmtpkg.Parse(*formatFlag)
	if err != nil {
		die("%v", err)
	}
	runEval(resolveDeckPath(fs.Arg(0)), *deepShuffles, *incoming, *maxCopies, *seed, fmtValue, *reevaluate, *brief)
}

// runEval loads the deck at outPath, simulates it for deepShuffles hands, and prints the
// fresh stats. By default the file on disk is NOT overwritten — eval is a read-only
// measurement so a deck can be re-scored at a new shuffle depth or different -incoming
// without clobbering the saved stats. Two cases force a rewrite:
//
//   - A deck that arrived with card.NotImplemented copies has those slots replaced in
//     memory before scoring and the sanitized result is written back, because keeping the
//     old file would leave the on-disk avg forever out of sync with the cards we can
//     actually simulate.
//   - The caller passed -reevaluate to explicitly refresh the persisted stats. Use this
//     to bring a saved file in sync with the current binary's best-turn output.
//
// Output shape is controlled by brief:
//   - brief=false (default): full printBestDeck dump — summary, card list, best-turn block,
//     per-card stats. Same shape as `fabsim print`.
//   - brief=true: score summary only. Good for scripted re-scoring where the card list and
//     best turn are noise.
func runEval(outPath string, deepShuffles, incoming, maxCopies int, seed int64, fmtValue fmtpkg.Format, reevaluate, brief bool) {
	loaded := mustLoadDeck(outPath)
	// Wrap the loaded hero/weapons/cards in a fresh Deck so Evaluate's stats start from zero
	// instead of accumulating on top of the persisted Stats. Sideboard carries over verbatim
	// — the sim ignores it but the post-eval writeDeck round-trips it back to disk.
	d := deck.New(loaded.Hero, loaded.Weapons, loaded.Cards)
	d.Sideboard = loaded.Sideboard
	rng := rand.New(rand.NewSource(seed))
	savedAvg := loaded.Stats.Mean()
	replaced := sanitizeLoadedDeck(d, maxCopies, rng, fmtValue.IsLegal)
	d.Evaluate(deepShuffles, incoming, rng)
	var reason string
	switch {
	case len(replaced) > 0:
		reason = "warning: sanitized"
	case reevaluate:
		reason = "reevaluate:"
	}
	if reason != "" {
		fmt.Fprintf(os.Stderr, "%s avg %.3f → %.3f (delta %+.3f); rewriting %s\n",
			reason, savedAvg, d.Stats.Mean(), d.Stats.Mean()-savedAvg, outPath)
		if err := writeDeck(d, outPath); err != nil {
			die("%v", err)
		}
	}
	if brief {
		printDeckSummary(d)
		return
	}
	printBestDeck(d)
}
