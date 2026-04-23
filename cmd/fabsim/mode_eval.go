package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
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
	_ = parseFlagsAnywhere(fs, args)
	if fs.NArg() != 1 {
		die("eval: need exactly one positional <deck> (got %d); try `fabsim eval <deck>`", fs.NArg())
	}
	requireFlag(fs, "eval", "incoming")
	fmtValue, err := fmtpkg.Parse(*formatFlag)
	if err != nil {
		die("%v", err)
	}
	runEval(resolveDeckPath(fs.Arg(0)), *deepShuffles, *incoming, *seed, fmtValue)
}

// runEval loads the deck at outPath, simulates it for deepShuffles hands, and prints the
// fresh stats. The file on disk is NOT overwritten: eval is a read-only measurement so a
// deck can be re-scored at a new shuffle depth or different -incoming without clobbering
// the saved stats. Exception: a deck that arrived with card.NotImplemented copies has
// those slots replaced in memory before scoring and the sanitized result is written back,
// because keeping the old file would leave the on-disk avg forever out of sync with the
// cards we can actually simulate.
//
// Prints only the score summary — not the card list — so the freshly-computed mean stays
// visible on a small terminal. The deck's contents haven't changed, so repeating them here
// just scrolls the score off the top; `fabsim print <deck>` is the command for the full dump.
func runEval(outPath string, deepShuffles, incoming int, seed int64, fmtValue fmtpkg.Format) {
	loaded := mustLoadDeck(outPath)
	// Wrap the loaded hero/weapons/cards in a fresh Deck so Evaluate's stats start from zero
	// instead of accumulating on top of the persisted Stats.
	d := deck.New(loaded.Hero, loaded.Weapons, loaded.Cards)
	rng := rand.New(rand.NewSource(seed))
	savedAvg := loaded.Stats.Mean()
	maxCopies := inferMaxCopies(d.Cards)
	replaced := sanitizeLoadedDeck(d, maxCopies, rng, fmtValue.IsLegal)
	d.Evaluate(deepShuffles, incoming, rng)
	if len(replaced) > 0 {
		fmt.Fprintf(os.Stderr, "warning: sanitized avg %.3f → %.3f (delta %+.3f); rewriting %s\n",
			savedAvg, d.Stats.Mean(), d.Stats.Mean()-savedAvg, outPath)
		if err := writeDeck(d, outPath); err != nil {
			die("%v", err)
		}
	}
	printDeckSummary(d)
}

// inferMaxCopies returns the highest per-printing count present in cs. Used by eval to
// derive a copies-cap for in-place sanitization when the caller didn't supply one: a deck
// already at N copies of some card stays legal under a cap of N, so picking max(count)
// keeps the existing distribution valid post-sanitize.
func inferMaxCopies(cs []card.Card) int {
	counts := map[card.ID]int{}
	maxCount := 1
	for _, c := range cs {
		counts[c.ID()]++
		if counts[c.ID()] > maxCount {
			maxCount = counts[c.ID()]
		}
	}
	return maxCount
}
