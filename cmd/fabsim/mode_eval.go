package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/deck"
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
	sim.OptDebug = *debug
	runEval(resolveDeckPath(fs.Arg(0)), *shuffles, *precision, sim.Matchup{IncomingDamage: *incoming, ArcaneIncomingDamage: *arcaneIncoming}, *seed, *printOnly, *brief, *debug)
}

// runEval loads the deck at outPath and prints its stats. Default behaviour (printOnly=false)
// first re-simulates the deck for deepShuffles hands against incoming and writes the fresh
// stats back to disk — both the JSON and the sibling fabrary .txt — so the on-disk copy
// always reflects the latest binary's modelling. printOnly=true skips that step and just
// loads-and-prints, which is what you want for a quick look at a saved deck without spending
// shuffles or mutating the file. Both branches share the same load-and-print path so the
// output is identical regardless of whether a sim ran first.
//
// Output shape is controlled by brief:
//   - brief=false (default): full printBestDeck dump — summary, card list, best-turn block,
//     per-card stats.
//   - brief=true: score summary only. Good for scripted re-scoring where the card list and
//     best turn are noise.
//
// debug=true prints extra telemetry to stderr after the run — currently the hand-eval
// cache hit rate. Only meaningful when a fresh simulation actually ran (printOnly=false);
// otherwise the Evaluator never spun up.
func runEval(outPath string, shuffles int, precision float64, mp sim.Matchup, seed int64, printOnly, brief, debug bool) {
	if !printOnly {
		evaluateAndPersist(outPath, shuffles, precision, mp, seed, debug)
	}
	d, stats := mustLoadDeck(outPath)
	printLoadedDeck(d, stats, brief)
}

// evaluateAndPersist runs the deck eval — adaptive when shuffles is negative (capped at
// adaptiveShufflesCap), fixed otherwise — then writes the fresh stats back to disk
// (.json + sibling fabrary .txt). Returns the simulated deck so callers can print its
// stats. The stderr summary lets the operator see the re-score happening before the
// printed output appears.
//
// Always uses a dedicated Evaluator (rather than the package-level shared one) so the
// per-Evaluator cache stats are always available. debug=true prints them after the run;
// otherwise they're computed-but-discarded — the cache itself runs unconditionally because
// it speeds up the eval regardless of whether the operator wants the telemetry.
func evaluateAndPersist(outPath string, shuffles int, precision float64, mp sim.Matchup, seed int64, debug bool) (*deck.Deck, sim.DeckStats) {
	loaded, loadedStats := mustLoadDeck(outPath)
	// Wrap the loaded hero/weapons/cards in a fresh Deck so the eval's stats start from zero
	// instead of accumulating on top of the persisted ones. Sideboard and Equipment carry
	// over verbatim — the sim ignores both, but the post-eval writeDeck round-trips them
	// back to disk so the user's hand-managed lists aren't dropped by a re-score.
	d := loaded.Copy()
	d.Sideboard = loaded.Sideboard
	d.Equipment = loaded.Equipment
	rng := rand.New(rand.NewSource(seed))
	savedAvg := loadedStats.Mean()
	// Parallel-shuffle eval: workers fan the shuffle loop across all available cores,
	// sharing the cache via the RWMutex-protected lookup path. fabsim eval is the
	// flagship single-deck workload — getting from 1.8s to ~0.5s on 8 workers cuts
	// re-score wall-clock noticeably.
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

// printCacheStats writes the hand-eval cache counters to stderr as a single annotated
// line. Lives in this file because eval is the only mode that exposes per-Evaluator stats
// today — iterate / anneal use a worker pool with one Evaluator per worker, so a single
// cache-stats line wouldn't capture the workload.
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

// safePct returns num/denom or 0 when denom is 0. Spares the printCacheStats format string
// from peppering checks at every call.
func safePct(num, denom int) float64 {
	if denom == 0 {
		return 0
	}
	return float64(num) / float64(denom)
}

// printLoadedDeck dispatches between the brief summary and the full printBestDeck dump;
// used by both the simulate path and -print-only.
func printLoadedDeck(d *deck.Deck, s sim.DeckStats, brief bool) {
	if brief {
		printDeckSummary(d, s)
		return
	}
	printBestDeck(d, s)
}
