package main

import (
	"context"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"runtime"
	"runtime/pprof"
	"strings"
	"sync/atomic"
	"time"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/deckformat"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/heroes"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/mydecks"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

// annealConfig bundles the knobs runAnneal needs. Built by runAnnealCmd from its flag.FlagSet.
type annealConfig struct {
	// shuffles is the per-eval shuffle budget when adaptive is false (apples-to-apples
	// acceptance, repro flows). Ignored when adaptive is true — the deck package's
	// adaptive path uses its own SE target and cap.
	shuffles   int
	adaptive   bool
	matchup    sim.Matchup
	deckSize   int
	maxCopies  int
	seed       int64
	outPath    string
	format     deckformat.Format
	debug      bool
	reevaluate bool
	// startTemp / tempDecay / minTemp are the simulated-annealing knobs. startTemp of 0
	// degenerates to the classical hill-climb (strict > baseline acceptance).
	startTemp float64
	tempDecay float64
	minTemp   float64
	// minImprovement is the noise-floor a strict (T==0) acceptance must clear: avg must
	// exceed bestAvg by more than this margin. Guards against infinite acceptance loops where
	// shuffle noise lets a stream of near-zero "wins" keep firing. The probabilistic SA gate
	// ignores this floor so annealing can still cross ties / dips.
	minImprovement float64
	// quietLoad suppresses the baseline card-list dump in prepareBaseline. Set by wrapper
	// scripts that re-invoke anneal repeatedly on the same deck; the listing is unchanging
	// noise after the first pass.
	quietLoad bool
	// maxDuration caps the run's wall-clock time. Zero means no cap (anneal runs until the user
	// hits Enter). The deadline aborts the same way the stdin watcher does, so any outstanding
	// round finishes evaluating before the loop exits and the deck still saves.
	maxDuration time.Duration
}

// legalFilter returns the card-pool predicate for this run's format. anneal always runs under a
// format, so this is non-nil; the deck package accepts nil for "no filtering" generally.
func (c annealConfig) legalFilter() func(sim.Card) bool {
	return c.format.IsLegal
}

// defaultDeckNameFor returns the deck name when -deck isn't supplied, keyed by hero, format, and
// -incoming. Different regimes produce different optimal decks, so each gets its own file to
// avoid hill-climbing one regime's best under another regime's objective.
func defaultDeckNameFor(h sim.Hero, f deckformat.Format, incoming int) string {
	return fmt.Sprintf("%s_%s_%d_incoming", strings.ToLower(h.Name()), f, incoming)
}

// runAnnealCmd parses anneal's flags from args and dispatches to runAnneal. Every anneal-only
// knob lives on this FlagSet so `fabsim anneal -help` shows exactly the flags that apply here.
func runAnnealCmd(args []string) {
	fs := flag.NewFlagSet("anneal", flag.ExitOnError)
	deckName := fs.String("deck", "", "deck name; resolved to mydecks/<name>.json (\".json\" suffix optional). Defaults to <hero>_<format>_<incoming>_incoming so different (hero, format, -incoming) regimes keep separate deck files. When the named deck exists, anneal resumes from it as a checkpoint.")
	shuffles := fs.Int("shuffles", -1, "per-eval shuffle budget. -1 (default) runs adaptively, stopping once the per-turn mean's standard error drops below the built-in target. Any non-negative value runs exactly that many shuffles for apples-to-apples acceptance / repro flows.")
	incoming := fs.Int("incoming", 0, "opponent damage per turn (required — different values produce different optimal decks, so this is explicit rather than defaulted)")
	arcaneIncoming := fs.Int("arcane-incoming", 0, "opponent arcane damage per turn (defaults to 0 — the non-arcane matchup; raise it to score cards that gate on incoming arcane)")
	deckSize := fs.Int("deck-size", 40, "number of cards per deck")
	maxCopies := fs.Int("max-copies", defaultMaxCopies, "maximum copies of any single card printing per deck")
	seed := fs.Int64("seed", time.Now().UnixNano(), "RNG seed")
	formatFlag := fs.String("format", string(deckformat.SilverAge), "constructed format whose banlist restricts the card pool during search (only \"silver_age\" is supported today)")
	debug := fs.Bool("debug", false, "force per-round logs even when annealing is on (T>0 normally hides them); also prints every Opt() outcome to stdout as it fires")
	reevaluate := fs.Bool("reevaluate", false, "force re-evaluation of the loaded deck's baseline avg, even if its prior run count already matches the current -shuffles budget. Use after adjusting modelling assumptions or fixing bugs that may have shifted the deck's true score.")
	finalize := fs.Bool("finalize", false, "high-precision pass — sets -shuffles to 100000 (fixed) and tightens -min-improvement to 0.01. Use on a deck that's already converged to squeeze out the remaining sub-percent improvements.")
	startTemp := fs.Float64("start-temp", 0, "simulated-annealing starting temperature. 0 (default) runs a pure hill climb. Higher values probabilistically accept worse mutations early; acceptance probability is exp((avg - baseline) / T). Good starting range is ~0.05–0.5 given typical Value units.")
	minImprovement := fs.Float64("min-improvement", 0.1, "noise floor on strict (T==0) acceptance: a mutation's avg must exceed the current avg by more than this margin to be accepted. Guards against infinite-loop acceptance of within-noise wins; raise it for chunkier improvements only, lower it (e.g. 0.01) for fine-grained finalize passes. The probabilistic SA gate at T>0 ignores this margin so annealing can still cross ties.")
	tempDecay := fs.Float64("temp-decay", 0.95, "multiplicative cooling per acceptance — T ← T × decay, floored at -min-temp. Unused when -start-temp is 0.")
	minTemp := fs.Float64("min-temp", 0, "minimum temperature. Once T reaches this floor the climb becomes greedy until a local maximum is found. 0 disables annealing in the converged tail.")
	quietLoad := fs.Bool("quiet-load", false, "skip the baseline card-list dump at startup. Intended for wrapper scripts (e.g. anneal-reanneal.ps1) that re-invoke anneal many times on the same deck — the listing never changes pass-to-pass and floods the log.")
	cpuprofile := fs.String("cpuprofile", "", "if set, write a CPU profile to this path covering the entire anneal run. Pair with -max-duration for a time-boxed profile-driven optimization pass.")
	memprofile := fs.String("memprofile", "", "if set, write a heap profile to this path at exit (after a runtime.GC()).")
	maxDuration := fs.Duration("max-duration", 0, "cap wall-clock duration; the run aborts cleanly at the deadline like a stdin Enter. Zero (default) runs until the user hits Enter.")
	_ = parseFlagsAnywhere(fs, args)
	if fs.NArg() > 0 {
		die("anneal: unexpected positional argument(s): %v (did you mean -deck %s?)", fs.Args(), fs.Args()[0])
	}
	requireFlag(fs, "anneal", "incoming")

	fmtValue, err := deckformat.Parse(*formatFlag)
	if err != nil {
		die("%v", err)
	}

	sim.OptDebug = *debug

	// -finalize bundles the high-precision overrides — pinned shuffle count plus a tighter
	// noise floor so sub-0.1 wins land that the default 0.1 -min-improvement gate would reject.
	// Applied as a post-parse override so it composes cleanly with explicit per-flag values.
	if *finalize {
		*shuffles = 100000
		*minImprovement = 0.01
	}

	name := *deckName
	if name == "" {
		name = defaultDeckNameFor(heroes.Viserai{}, fmtValue, *incoming)
	}
	outPath, err := mydecks.Path(name)
	if err != nil {
		die("%v", err)
	}

	cfg := annealConfig{
		shuffles:       *shuffles,
		adaptive:       *shuffles < 0,
		matchup:        sim.Matchup{IncomingDamage: *incoming, ArcaneIncomingDamage: *arcaneIncoming},
		deckSize:       *deckSize,
		maxCopies:      *maxCopies,
		seed:           *seed,
		outPath:        outPath,
		format:         fmtValue,
		debug:          *debug,
		reevaluate:     *reevaluate,
		startTemp:      *startTemp,
		tempDecay:      *tempDecay,
		minTemp:        *minTemp,
		minImprovement: *minImprovement,
		quietLoad:      *quietLoad,
		maxDuration:    *maxDuration,
	}

	// Run inside a wrapper that owns the profile lifecycle: deferred StopCPUProfile / heap dump
	// must fire on every exit path (clean finish OR abort) but os.Exit skips defers, so do the
	// exit-code dispatch only after the wrapper returns.
	aborted := runAnnealWithProfiling(cfg, *cpuprofile, *memprofile)
	if aborted {
		os.Exit(130)
	}
}

// runAnnealWithProfiling wraps runAnneal with optional CPU + heap profile capture. CPU profile
// runs across the entire anneal session; heap profile snapshots once at exit after a forced GC
// so live-only allocations dominate the result. Returns the aborted flag so the caller can
// pick the right exit code.
func runAnnealWithProfiling(cfg annealConfig, cpuprofile, memprofile string) bool {
	if cpuprofile != "" {
		f, err := os.Create(cpuprofile)
		if err != nil {
			die("create cpuprofile: %v", err)
		}
		defer f.Close()
		if err := pprof.StartCPUProfile(f); err != nil {
			die("start cpuprofile: %v", err)
		}
		defer pprof.StopCPUProfile()
	}
	res := runAnneal(cfg)
	fmt.Fprintf(os.Stderr, "\nSession summary: avg %.3f → %.3f (%+.3f)\n",
		res.startingAvg, res.bestEverAvg, res.bestEverAvg-res.startingAvg)
	if memprofile != "" {
		f, err := os.Create(memprofile)
		if err != nil {
			die("create memprofile: %v", err)
		}
		defer f.Close()
		runtime.GC()
		if err := pprof.WriteHeapProfile(f); err != nil {
			die("write memprofile: %v", err)
		}
	}
	return res.aborted
}

// annealResult carries the outcome of a single runAnneal pass. aborted is true when the
// user hit Enter (stdin watcher fired) — callers propagate this up to main so the process
// exits non-zero and wrapper scripts like anneal-reanneal.ps1 stop their outer loop instead
// of immediately launching another pass.
type annealResult struct {
	bestEverAvg float64
	startingAvg float64
	aborted     bool
}

func runAnneal(cfg annealConfig) annealResult {
	rng := rand.New(rand.NewSource(cfg.seed))

	current, currentAvg := prepareBaseline(cfg, rng)
	// All-time best tracks the highest-avg deck seen since runAnneal started. The saved JSON
	// mirrors this — simulated annealing intentionally walks through worse states to escape
	// local maxima, but the on-disk artifact should always reflect the peak reached so far.
	bestEver := current
	bestEverAvg := currentAvg
	startingAvg := currentAvg
	fmt.Println("Press Enter to abort.")

	var ctx context.Context
	var cancel context.CancelFunc
	if cfg.maxDuration > 0 {
		ctx, cancel = context.WithTimeout(context.Background(), cfg.maxDuration)
	} else {
		ctx, cancel = context.WithCancel(context.Background())
	}
	defer cancel()
	watchStdinForAbort(cancel)

	temperature := cfg.startTemp
	// verbose gates the noisy per-round headers and per-mutation acceptance lines. Classical
	// hill-climb (T==0) keeps them on because the log is the user's progress indicator; under
	// annealing they flood the terminal with hundreds of near-identical lines, so they're
	// suppressed unless -debug asked for them.
	verbose := temperature == 0 || cfg.debug
	if temperature > 0 {
		fmt.Fprintf(os.Stderr, "Simulated annealing: startTemp=%.3f decay=%.3f minTemp=%.3f\n",
			cfg.startTemp, cfg.tempDecay, cfg.minTemp)
	}

	round := 0
	acceptances := 0
	start := time.Now()
	for {
		round++
		mutations := buildRoundMutations(cfg, rng, current)
		tempLabel := formatTempLabel(temperature)
		if verbose {
			fmt.Fprintf(os.Stderr, "\n[round %d] evaluating %d mutations of avg %.3f%s (best ever %.3f)\n",
				round, len(mutations), currentAvg, tempLabel, bestEverAvg)
		}

		// Round-scoped start so the ticker's elapsed/ETA reflect this round's progress
		// rather than the cumulative session.
		var completed atomic.Int64
		roundStart := time.Now()
		stopTicker := startRoundTicker(round, len(mutations), roundStart, &completed,
			temperature, currentAvg, bestEverAvg)
		d, avg, idx, found := sim.IterateParallel(
			ctx, mutations, currentAvg, temperature, cfg.minImprovement,
			cfg.shuffles, cfg.matchup, 0, 0,
			rng.Int63(), &completed, cfg.adaptive,
		)
		stopTicker()

		if ctx.Err() != nil {
			return finishAnnealRun(cfg, bestEver, bestEverAvg, startingAvg,
				fmt.Sprintf("Aborted mid-round after %d rounds / %d acceptances in %s",
					round, acceptances, time.Since(start).Truncate(time.Second)),
				true)
		}
		if !found {
			// A full round with zero acceptances means every mutation — including the
			// probabilistically-accepted worse ones — failed the gate. At any T > 0 with
			// thousands of mutations this is vanishingly unlikely unless we've genuinely
			// converged, so treat it as a local maximum regardless of the current temperature.
			return finishAnnealRun(cfg, bestEver, bestEverAvg, startingAvg,
				fmt.Sprintf("Local maximum reached after %d rounds / %d acceptances in %s",
					round, acceptances, time.Since(start).Truncate(time.Second)),
				false)
		}

		acceptances++
		bestEver, bestEverAvg = applyAcceptedMutation(cfg, round, verbose, tempLabel,
			idx, len(mutations), mutations[idx], d, avg, currentAvg, bestEver, bestEverAvg)
		current = d
		currentAvg = avg
		temperature = coolDown(temperature, cfg.tempDecay, cfg.minTemp)
	}
}

// buildRoundMutations produces the per-round mutation list: enumerates every
// single-card/weapon mutation and shuffles the order so exploration is unbiased.
//
// AllMutations returns a ids.CardID-sorted slice for stability; the unconditional shuffle here
// is what keeps the first-improvement classical climb from sampling the head of the slice
// disproportionately, and what keeps the probabilistic SA gate from concentrating its
// acceptances on a fixed slice of the solution space.
func buildRoundMutations(cfg annealConfig, rng *rand.Rand, current *sim.Deck) []sim.Mutation {
	mutations := sim.AllMutations(current, cfg.maxCopies, cfg.legalFilter())
	rng.Shuffle(len(mutations), func(i, j int) {
		mutations[i], mutations[j] = mutations[j], mutations[i]
	})
	return mutations
}

// formatTempLabel is the " (T=…)" suffix the round header and acceptance lines share. Empty
// string when temperature is 0 so classical hill-climb logs stay clean.
func formatTempLabel(temperature float64) string {
	if temperature <= 0 {
		return ""
	}
	return fmt.Sprintf(" (T=%.4f)", temperature)
}

// applyAcceptedMutation runs the after-success branch of a round: logs the acceptance (against
// the pre-mutation currentAvg so the "improvement" / "annealing step" label reflects whether
// this specific step walked uphill), logs a new all-time best in non-verbose mode, persists
// the deck to disk when avg exceeds bestEverAvg, and returns the possibly-updated bestEver /
// bestEverAvg. The current deck and its avg stay owned by the caller.
func applyAcceptedMutation(cfg annealConfig, round int, verbose bool, tempLabel string,
	idx, total int, mut sim.Mutation, d *sim.Deck, avg, currentAvg float64,
	bestEver *sim.Deck, bestEverAvg float64) (*sim.Deck, float64) {
	verb := "improvement"
	if avg <= currentAvg {
		verb = "annealing step"
	}
	if verbose {
		fmt.Fprintf(os.Stderr, "\r[round %d] %s at %d/%d: deep %.3f vs %.3f (%s)%s       \n",
			round, verb, idx+1, total, avg, currentAvg, mut.Description, tempLabel)
	}
	if avg <= bestEverAvg {
		return bestEver, bestEverAvg
	}
	if !verbose {
		// Surface every new all-time best in non-verbose annealing mode so long
		// reanneal sessions still show visible forward motion without the per-round
		// spam. \r + trailing padding overwrites the ticker line cleanly; the \n at
		// the end promotes this to a persistent entry above the next round's ticker.
		fmt.Fprintf(os.Stderr, "\r[round %d] new best %.3f (was %.3f, +%.3f)%s                                \n",
			round, avg, bestEverAvg, avg-bestEverAvg, tempLabel)
	}
	if err := writeDeck(d, cfg.outPath); err != nil {
		die("%v", err)
	}
	return d, avg
}

// finishAnnealRun emits the terminal status line (abort / converged), optionally prints the
// full best-ever deck listing, and builds the annealResult the top-level command surfaces as
// exit code and session summary. aborted is threaded through as-is because runAnnealCmd keys
// exit code 130 off it.
func finishAnnealRun(cfg annealConfig, bestEver *sim.Deck, bestEverAvg, startingAvg float64,
	statusLine string, aborted bool) annealResult {
	fmt.Fprintln(os.Stderr, "\n"+statusLine)
	fmt.Println()
	if shouldPrintFinalDeck(cfg.startTemp, bestEverAvg, startingAvg) {
		printBestDeck(bestEver)
	}
	return annealResult{bestEverAvg: bestEverAvg, startingAvg: startingAvg, aborted: aborted}
}

// shouldPrintFinalDeck decides whether to dump the full deck listing at the end of a run.
// Annealing sessions that exit with no net improvement (common in long reanneal loops that
// are just probing) suppress the listing to keep the session-summary line visible. Classical
// mode keeps printing regardless: a no-improvement run is a single round and the listing is
// the user's confirmation of what was evaluated.
func shouldPrintFinalDeck(startTemp, bestEverAvg, startingAvg float64) bool {
	return startTemp == 0 || bestEverAvg > startingAvg
}

// coolDown applies one round of geometric cooling, clamped at minTemp so the classical-mode
// hill climb is fully recovered once temperature reaches the floor.
func coolDown(temperature, decay, minTemp float64) float64 {
	next := temperature * decay
	if next < minTemp {
		return minTemp
	}
	return next
}

// baselineEvaluate runs the eval used for every prepareBaseline path. Adaptive when
// cfg.adaptive is true (no -shuffles pinned); fixed-shuffles otherwise. The two paths
// return the same Stats shape; Stats.Runs reflects the actual shuffle count so the next
// prepareBaseline call's "already evaluated" check still works (an adaptive run may finish
// below the cap and prompt a re-evaluation next session, which is fine — adaptive runs
// are cheap). Routes through evaluateParallel so the once-per-session baseline benefits
// from the same DefaultWorkers fan-out as iterate's per-mutation evals.
func baselineEvaluate(d *sim.Deck, cfg annealConfig, rng *rand.Rand) sim.Stats {
	shuffles := cfg.shuffles
	if cfg.adaptive {
		shuffles = -1
	}
	stats, _ := evaluateParallel(d, shuffles, cfg.matchup, rng)
	return stats
}

// prepareBaseline returns the starting deck for the hill climb with its baseline avg.
// Four cases: no deck on disk (generate random + evaluate); loaded deck under the current
// shuffle budget (re-evaluate for an apples-to-apples baseline); -reevaluate set (force
// re-evaluation even if the run count already matches); or deck already evaluated at the
// budget (use as-is). File exists but doesn't parse → die loudly rather than silently
// overwrite a corrupt checkpoint.
//
// A loaded deck that contains sim.NotImplemented copies (e.g. a pre-tag deck recovered
// from disk) is sanitized before any of the above branches: the tagged slots are replaced
// with random legal picks and the run always takes the re-evaluate path so the baseline
// reflects the new card list.
func prepareBaseline(cfg annealConfig, rng *rand.Rand) (*sim.Deck, float64) {
	best, bestAvg, err := loadExisting(cfg.outPath)
	if err != nil {
		die("%v", err)
	}
	if best == nil {
		fmt.Fprintf(os.Stderr, "no deck at %s; generating a random starting deck\n", cfg.outPath)
		best = sim.Random(heroes.Viserai{}, cfg.deckSize, cfg.maxCopies, rng, cfg.legalFilter())
		bestAvg = baselineEvaluate(best, cfg, rng).Mean()
		if err := writeDeck(best, cfg.outPath); err != nil {
			die("%v", err)
		}
		fmt.Printf("Starting deck avg %.3f, saved to %s\n", bestAvg, cfg.outPath)
		maybePrintBaselineCards(cfg, best)
		return best, bestAvg
	}
	// Sanitize in place before the evaluation-branch decision. Any swap forces the
	// re-evaluate path below by zeroing Stats.Runs so the saved run count can't satisfy
	// the "already deep-evaluated" check against a now-different card list. bestAvg stays
	// at the loaded value for the "saved avg → current avg" delta display below.
	sanitized := sanitizeLoadedDeck(best, cfg.maxCopies, rng, cfg.legalFilter())
	if len(sanitized) > 0 {
		best.Stats.Runs = 0
	}
	// Re-evaluate when the saved deck was scored at fewer shuffles than the current budget,
	// or when -reevaluate forces it. Adaptive runs always take the re-evaluate path because
	// the recorded Stats.Runs reflects whatever count the previous adaptive run terminated
	// at, which carries no precision guarantee for the new run's target.
	needReeval := cfg.adaptive || cfg.reevaluate || best.Stats.Runs < cfg.shuffles
	if needReeval {
		best, bestAvg = reevaluateBaseline(cfg, rng, best, bestAvg, sanitized)
		maybePrintBaselineCards(cfg, best)
		return best, bestAvg
	}
	fmt.Printf("Loaded best deck (avg %.3f) from %s\n", bestAvg, cfg.outPath)
	maybePrintBaselineCards(cfg, best)
	return best, bestAvg
}

// reevaluateBaseline rebuilds the loaded deck against the current shuffle budget and writes
// the refreshed Stats back to disk. Picks an explanatory reason label (sanitized cards
// replaced, -reevaluate forced, or stale shuffle count), reconstructs the deck so any saved
// Stats are dropped (Sideboard and Equipment are preserved), runs baselineEvaluate, and
// persists the result. Returns the rebuilt deck and its fresh avg.
func reevaluateBaseline(cfg annealConfig, rng *rand.Rand, loaded *sim.Deck, savedAvg float64, sanitized []sim.NotImplementedReplacement) (*sim.Deck, float64) {
	reason := fmt.Sprintf("from %d shuffles", loaded.Stats.Runs)
	if cfg.reevaluate && loaded.Stats.Runs >= cfg.shuffles {
		reason = "-reevaluate forced"
	}
	if len(sanitized) > 0 {
		reason = fmt.Sprintf("%d NotImplemented card(s) replaced", len(sanitized))
	}
	// Label the loaded number "saved avg" so it can't be mistaken for the re-evaluated
	// score. Decks scored under older simulation logic can have saved avgs that diverge
	// substantially from what today's simulator produces.
	budgetLabel := fmt.Sprintf("%d shuffles", cfg.shuffles)
	if cfg.adaptive {
		budgetLabel = "adaptive shuffles"
	}
	fmt.Printf("Loaded best deck (saved avg %.3f, %s); re-evaluating at %s for an apples-to-apples baseline\n",
		savedAvg, reason, budgetLabel)
	// Sideboard and Equipment are user-managed and don't feed the sim — preserve them across
	// the stats reset so the re-evaluated deck writes back unchanged.
	sideboard := loaded.Sideboard
	equipment := loaded.Equipment
	rebuilt := sim.New(loaded.Hero, loaded.Weapons, loaded.Cards)
	rebuilt.Sideboard = sideboard
	rebuilt.Equipment = equipment
	freshAvg := baselineEvaluate(rebuilt, cfg, rng).Mean()
	if err := writeDeck(rebuilt, cfg.outPath); err != nil {
		die("%v", err)
	}
	// Show saved→current so the delta from any simulation-logic drift is visible at a glance,
	// instead of the user guessing which of the two printed numbers is the fresh one.
	fmt.Printf("Re-evaluated baseline: %.3f → %.3f, saved to %s\n", savedAvg, freshAvg, cfg.outPath)
	return rebuilt, freshAvg
}

// maybePrintBaselineCards emits the startup card-list dump unless -quiet-load suppressed it. The
// leading blank line is part of the listing block, so it's also gated — otherwise -quiet-load
// would leave a lone empty line hanging after the baseline avg summary.
func maybePrintBaselineCards(cfg annealConfig, d *sim.Deck) {
	if cfg.quietLoad {
		return
	}
	fmt.Println()
	printCardList(d)
}

// watchStdinForAbort spawns a background goroutine that calls cancel() on the first keypress.
// EOF / closed stdin isn't an abort (so iterate doesn't exit immediately on non-TTY stdin); only
// a successful read of at least one byte counts. Cancellation propagates into IterateParallel so
// an abort takes effect mid-round.
func watchStdinForAbort(cancel context.CancelFunc) {
	go func() {
		buf := make([]byte, 1)
		if n, err := os.Stdin.Read(buf); err == nil && n > 0 {
			cancel()
		}
	}()
}

// startRoundTicker launches a 500ms ticker that renders round progress plus the annealing
// state (T, current avg, best-ever) to a \r-terminated stderr line. Runs in both classical
// and annealing modes: in annealing mode it's effectively the only ongoing progress indicator
// (the per-round / per-mutation logs are suppressed without -debug), so keeping the snapshot
// rich enough to track the walk is what makes silent mode bearable. ETA is the per-round
// time-to-finish projection assuming all `total` mutations are evaluated at the current
// rate; it under-estimates by the fraction of mutations cut short by an early acceptance,
// but over-estimates dominate in long converged-tail runs. Returns a stop function the
// caller must call when the round finishes.
func startRoundTicker(round, total int, roundStart time.Time, completed *atomic.Int64, temperature, currentAvg, bestEverAvg float64) func() {
	done := make(chan struct{})
	go func() {
		t := time.NewTicker(500 * time.Millisecond)
		defer t.Stop()
		for {
			select {
			case <-done:
				return
			case <-t.C:
				tempLabel := ""
				if temperature > 0 {
					tempLabel = fmt.Sprintf("  T=%.4f", temperature)
				}
				elapsed := time.Since(roundStart)
				fmt.Fprintf(os.Stderr, "\r[round %d] tested %d/%d  cur %.3f  best %.3f%s  %s elapsed%s        ",
					round, completed.Load(), total, currentAvg, bestEverAvg, tempLabel,
					elapsed.Truncate(time.Second), formatETA(elapsed, completed.Load(), int64(total)))
			}
		}
	}()
	return func() { close(done) }
}

// formatETA returns the " ETA <duration>" suffix for a progress line, or "" when no
// estimate is yet meaningful (no work done, or work already complete). Duration is the
// projected time remaining at the current rate, truncated to seconds.
func formatETA(elapsed time.Duration, done, total int64) string {
	if done <= 0 || done >= total {
		return ""
	}
	remaining := total - done
	eta := time.Duration(int64(elapsed) * remaining / done)
	return "  ETA " + eta.Truncate(time.Second).String()
}
