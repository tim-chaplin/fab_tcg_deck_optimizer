package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/deck"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/gameengine"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
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
	shuffles := fs.Int("shuffles", -1, "per-eval shuffle budget. -1 (default) runs adaptively to -precision. Any non-negative value runs exactly that many shuffles (apples-to-apples re-scores, repro flows).")
	precision := fs.Float64("precision", 0.1, "adaptive-eval precision target — stop once the per-turn mean's standard error falls to precision/4 (≈ ±precision/2 on the reported value with 95% confidence). Only relevant when -shuffles is negative (adaptive mode).")
	incoming := fs.Int("incoming", 0, "opponent damage per turn (required unless -print-only is set — must match the value the deck was annealed at for comparable numbers)")
	arcaneIncoming := fs.Int("arcane-incoming", 0, "opponent arcane damage per turn (defaults to 0 — the non-arcane matchup; raise it to score cards that gate on incoming arcane)")
	seed := fs.Int64("seed", time.Now().UnixNano(), "RNG seed")
	printOnly := fs.Bool("print-only", false, "load the deck and print the stats from the last run without simulating or rewriting the on-disk .json / .txt")
	brief := fs.Bool("brief", false, "print only the score summary (no card list, per-card stats, or best turn)")
	debug := fs.Bool("debug", false, "print additional debug info to stdout / stderr")
	_ = parseFlagsAnywhere(fs, args)
	if fs.NArg() != 1 {
		die("eval: need exactly one positional <deck> (got %d); try `fabsim eval <deck>`", fs.NArg())
	}
	if !*printOnly {
		requireFlag(fs, "eval", "incoming")
	}
	gameengine.OptDebug = *debug
	runEval(resolveDeckPath(fs.Arg(0)), *shuffles, *precision, sim.Matchup{IncomingDamage: *incoming, ArcaneIncomingDamage: *arcaneIncoming}, *seed, *printOnly, *brief, *debug)
}

// runEval loads the deck at outPath and prints its stats. With printOnly=false it first
// re-simulates against mp and writes the fresh stats back to disk (.json + sibling fabrary
// .txt). brief=true prints the score summary only; brief=false prints the full
// printBestDeck dump. debug=true prints cache telemetry to stderr after a fresh sim.
func runEval(outPath string, shuffles int, precision float64, mp sim.Matchup, seed int64, printOnly, brief, debug bool) {
	if !printOnly {
		// Print from the in-memory stats — the disk round-trip would drop Stats.PrintBest.
		d, stats := evaluateAndPersist(outPath, shuffles, precision, mp, seed, debug)
		printLoadedDeck(d, stats, brief)
		return
	}
	d, stats := mustLoadDeck(outPath)
	printLoadedDeck(d, stats, brief)
}

// evaluateAndPersist runs the eval (adaptive when shuffles<0, fixed otherwise), writes the
// fresh stats back to disk (.json + sibling fabrary .txt), and returns the simulated deck +
// stats. Prints a one-line before/after summary to stderr; debug=true also prints the
// dedicated Evaluator's cache stats.
func evaluateAndPersist(outPath string, shuffles int, precision float64, mp sim.Matchup, seed int64, debug bool) (*deck.Deck, deck.Stats) {
	loaded, loadedStats := mustLoadDeck(outPath)
	// Fresh Deck so the eval's stats start from zero; carry Sideboard and Equipment so
	// writeDeck round-trips them back to disk verbatim.
	d := loaded.Copy()
	d.Sideboard = loaded.Sideboard
	d.Equipment = loaded.Equipment
	rng := rand.New(rand.NewSource(seed))
	savedAvg := loadedStats.Mean()
	start := time.Now()
	stats, ev := evaluateParallel(d, shuffles, precision, mp, rng)
	elapsed := time.Since(start)
	freshAvg := stats.Mean()
	fmt.Fprintf(os.Stderr, "eval: avg %.3f → %.3f (delta %+.3f) in %s (%s shuffles); rewriting %s\n",
		savedAvg, freshAvg, freshAvg-savedAvg, elapsed.Round(time.Millisecond), commaInt(stats.Runs), outPath)
	if debug {
		printCacheStats(ev.CacheStats())
	}
	if err := writeDeck(d, stats, outPath); err != nil {
		die("%v", err)
	}
	return d, stats
}

// printCacheStats writes the hand-eval cache counters to stderr as a single annotated line.
func printCacheStats(s sim.CacheStats) {
	total := s.Hits + s.Misses
	if total == 0 {
		fmt.Fprintln(os.Stderr, "cache: no Best calls recorded")
		return
	}
	fmt.Fprintf(os.Stderr, "cache: %d calls — %.1f%% hits (%d), %.1f%% misses (%d); %d entries; %.1f%% of misses uncacheable\n",
		total,
		100*s.HitRate(), s.Hits,
		100*float64(s.Misses)/float64(total), s.Misses,
		s.Entries,
		100*safePct(s.Uncacheable, s.Misses),
	)
}

// safePct returns num/denom, or 0 when denom is 0.
func safePct(num, denom int) float64 {
	if denom == 0 {
		return 0
	}
	return float64(num) / float64(denom)
}

// printLoadedDeck dispatches between the brief summary and the full printBestDeck dump.
func printLoadedDeck(d *deck.Deck, s deck.Stats, brief bool) {
	if brief {
		printDeckSummary(d, s)
		return
	}
	printBestDeck(d, s)
}
