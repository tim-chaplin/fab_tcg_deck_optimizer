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

	"github.com/tim-chaplin/fab-deck-optimizer/internal/deck"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/gameengine"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/hero"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/hero/heroes"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/textio"
)

// annealCacheCapacity caps the cross-round hand-eval cache. At ~760 bytes / entry
// (TestEvalCache_MemoryPerEntry), 2M ≈ 1.5 GB resident — comfortable on a 16 GB+
// workstation and large enough that multi-hour sessions don't hit eviction.
const annealCacheCapacity = 2_000_000

// annealConfig bundles the knobs runAnneal needs. Built by runAnnealCmd from its flag.FlagSet.
type annealConfig struct {
	// shuffles is the per-eval shuffle budget when adaptive is false (apples-to-apples
	// acceptance, repro flows). Ignored when adaptive is true — the adaptive path stops
	// when SE drops to minImprovement/4 (capped at adaptiveShufflesCap).
	shuffles   int
	adaptive   bool
	matchup    sim.Matchup
	deckSize   int
	maxCopies  int
	seed       int64
	outPath    string
	debug      bool
	reevaluate bool
	// startTemp / tempDecay / minTemp are the SA knobs. startTemp=0 is a classical hill climb.
	startTemp float64
	tempDecay float64
	minTemp   float64
	// minImprovement is the strict (T==0) noise floor: avg must exceed bestAvg by more than
	// this margin. Guards against infinite-loop acceptance of within-noise wins. The
	// probabilistic SA gate ignores it so annealing can still cross ties / dips.
	minImprovement float64
	// quietLoad suppresses the baseline card-list dump in prepareBaseline — useful for
	// wrapper scripts that re-invoke anneal repeatedly on the same deck.
	quietLoad bool
	// maxDuration caps wall-clock time; zero means no cap (until user hits Enter). The
	// deadline aborts like stdin Enter — outstanding round finishes and the deck still saves.
	maxDuration time.Duration
	// pairMutations gates the pair-swap layer in AllMutations. See the -pair-mutations flag
	// help for when to enable it.
	pairMutations bool
}

// defaultDeckNameFor keys the default deck filename on (hero, format, incoming) so different
// regimes don't hill-climb one another's optimum.
func defaultDeckNameFor(h hero.Hero, f GameplayFormat, incoming int) string {
	return fmt.Sprintf("%s_%s_%d_incoming", strings.ToLower(h.Name()), f, incoming)
}

// runAnnealCmd parses anneal's flags from args and dispatches to runAnneal.
func runAnnealCmd(args []string) {
	fs := flag.NewFlagSet("anneal", flag.ExitOnError)
	deckName := fs.String("deck", "", "deck name; resolved to mydecks/<name>.json (\".json\" suffix optional). Defaults to <hero>_<format>_<incoming>_incoming so different (hero, format, -incoming) regimes keep separate deck files. When the named deck exists, anneal resumes from it as a checkpoint.")
	shuffles := fs.Int("shuffles", -1, "per-eval shuffle budget. -1 (default) runs adaptively to -min-improvement. Any non-negative value runs exactly that many shuffles for apples-to-apples acceptance / repro flows.")
	incoming := fs.Int("incoming", 0, "opponent damage per turn (required — different values produce different optimal decks, so this is explicit rather than defaulted)")
	arcaneIncoming := fs.Int("arcane-incoming", 0, "opponent arcane damage per turn (defaults to 0 — the non-arcane matchup; raise it to score cards that gate on incoming arcane)")
	deckSize := fs.Int("deck-size", 40, "number of cards per deck")
	maxCopies := fs.Int("max-copies", defaultMaxCopies, "maximum copies of any single card printing per deck")
	seed := fs.Int64("seed", time.Now().UnixNano(), "RNG seed")
	formatFlag := fs.String("format", string(SilverAge), "constructed format whose banlist restricts the card pool during search (only \"silver_age\" is supported today)")
	debug := fs.Bool("debug", false, "print additional debug info to stdout / stderr")
	reevaluate := fs.Bool("reevaluate", false, "force re-evaluation of the loaded deck's baseline avg, even if its prior run count already matches the current -shuffles budget. Use after adjusting modelling assumptions or fixing bugs that may have shifted the deck's true score.")
	finalize := fs.Bool("finalize", false, "high-precision pass — sets -shuffles to 100000 (fixed) and tightens -min-improvement to 0.01. Use on a deck that's already converged to squeeze out the remaining sub-percent improvements.")
	startTemp := fs.Float64("start-temp", 0, "simulated-annealing starting temperature. 0 (default) runs a pure hill climb. Higher values probabilistically accept worse mutations early; acceptance probability is exp((avg - baseline) / T). Good starting range is ~0.05–0.5 given typical Value units.")
	minImprovement := fs.Float64("min-improvement", 0.1, "noise floor on strict (T==0) acceptance: a mutation's avg must exceed the current avg by more than this margin to be accepted. Guards against infinite-loop acceptance of within-noise wins; raise it for chunkier improvements only, lower it (e.g. 0.01) for fine-grained finalize passes. The probabilistic SA gate at T>0 ignores this margin so annealing can still cross ties. Doubles as the adaptive-eval precision target so per-mutation evals resolve the mean to the granularity acceptance needs.")
	tempDecay := fs.Float64("temp-decay", 0.95, "multiplicative cooling per acceptance — T ← T × decay, floored at -min-temp. Unused when -start-temp is 0.")
	minTemp := fs.Float64("min-temp", 0, "minimum temperature. Once T reaches this floor the climb becomes greedy until a local maximum is found. 0 disables annealing in the converged tail.")
	quietLoad := fs.Bool("quiet-load", false, "skip the baseline card-list dump at startup. Intended for wrapper scripts (e.g. anneal-reanneal.ps1) that re-invoke anneal many times on the same deck — the listing never changes pass-to-pass and floods the log.")
	cpuprofile := fs.String("cpuprofile", "", "if set, write a CPU profile to this path covering the entire anneal run. Pair with -max-duration for a time-boxed profile-driven optimization pass.")
	memprofile := fs.String("memprofile", "", "if set, write a heap profile to this path at exit (after a runtime.GC()).")
	maxDuration := fs.Duration("max-duration", 0, "cap wall-clock duration; the run aborts cleanly at the deadline like a stdin Enter. Zero (default) runs until the user hits Enter.")
	pairMutations := fs.Bool("pair-mutations", false, "include the synergy-pair (-1/-1, +1/+1) swap layer in each round's mutation pool. Off by default because the layer multiplies per-round candidates substantially for little acceptance yield. Flip on when introducing a new synergy pair to give it a chance to land.")
	_ = parseFlagsAnywhere(fs, args)
	if fs.NArg() > 0 {
		die("anneal: unexpected positional argument(s): %v (did you mean -deck %s?)", fs.Args(), fs.Args()[0])
	}
	requireFlag(fs, "anneal", "incoming")

	fmtValue, err := parseGameplayFormat(*formatFlag)
	if err != nil {
		die("%v", err)
	}

	gameengine.OptDebug = *debug

	// -finalize tightens the noise floor so sub-0.1 wins land that the default 0.1
	// -min-improvement gate would reject. Applied post-parse so it composes with explicit flags.
	if *finalize {
		*minImprovement = 0.01
	}

	name := *deckName
	if name == "" {
		name = defaultDeckNameFor(heroes.Viserai, fmtValue, *incoming)
	}
	outPath, err := textio.MydecksPath(name)
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
		debug:          *debug,
		reevaluate:     *reevaluate,
		startTemp:      *startTemp,
		tempDecay:      *tempDecay,
		minTemp:        *minTemp,
		minImprovement: *minImprovement,
		quietLoad:      *quietLoad,
		maxDuration:    *maxDuration,
		pairMutations:  *pairMutations,
	}

	// Wrap in a function that owns the profile lifecycle so deferred StopCPUProfile / heap
	// dump fire on every exit path (clean finish OR abort). os.Exit skips defers, so dispatch
	// the exit code only after the wrapper returns.
	aborted := runAnnealWithProfiling(cfg, *cpuprofile, *memprofile)
	if aborted {
		os.Exit(130)
	}
}

// runAnnealWithProfiling wraps runAnneal with optional CPU + heap profile capture. CPU
// profile runs across the whole session; heap profile snapshots once at exit after a forced
// GC so live-only allocations dominate.
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
// user hit Enter; callers propagate this to main so the process exits non-zero and wrapper
// scripts stop their outer loop instead of relaunching.
type annealResult struct {
	bestEverAvg float64
	startingAvg float64
	aborted     bool
}

func runAnneal(cfg annealConfig) annealResult {
	rng := rand.New(rand.NewSource(cfg.seed))

	current, currentStats, currentAvg := prepareBaseline(cfg, rng)
	// All-time best tracks the highest-avg deck since runAnneal started; the saved JSON
	// mirrors it. SA walks through worse states to escape local maxima, but the on-disk
	// artifact always reflects the peak.
	bestEver := current
	bestEverStats := currentStats
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
	// verbose gates the per-round / per-mutation acceptance lines. Classical hill-climb
	// keeps them on as progress indicator; SA suppresses them unless -debug asked.
	verbose := temperature == 0 || cfg.debug
	if temperature > 0 {
		fmt.Fprintf(os.Stderr, "Simulated annealing: startTemp=%.3f decay=%.3f minTemp=%.3f\n",
			cfg.startTemp, cfg.tempDecay, cfg.minTemp)
	}

	// Persistent hand-eval cache shared across rounds — see sim.NewCacheBounded.
	roundCache := sim.NewCacheBounded(annealCacheCapacity)

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

		// Round-scoped start so the ticker's elapsed/ETA reflect this round, not the session.
		var completed atomic.Int64
		roundStart := time.Now()
		stopTicker := startRoundTicker(round, len(mutations), roundStart, &completed,
			temperature, currentAvg, bestEverAvg)
		d, dStats, avg, idx, found := sim.RunMutationRound(
			ctx, mutations, currentAvg, temperature, cfg.minImprovement,
			cfg.shuffles, cfg.matchup, 0, 0,
			rng.Int63(), &completed, cfg.adaptive, cfg.minImprovement, roundCache,
		)
		stopTicker()

		if ctx.Err() != nil {
			return finishAnnealRun(cfg, bestEver, bestEverStats, bestEverAvg, startingAvg,
				fmt.Sprintf("Aborted mid-round after %d rounds / %d acceptances in %s",
					round, acceptances, time.Since(start).Truncate(time.Second)),
				true)
		}
		if !found {
			// Full round with zero acceptances means every mutation (including
			// probabilistically-accepted worse ones) failed. At any T > 0 with thousands of
			// candidates this is vanishingly unlikely unless converged, so treat as local max.
			return finishAnnealRun(cfg, bestEver, bestEverStats, bestEverAvg, startingAvg,
				fmt.Sprintf("Local maximum reached after %d rounds / %d acceptances in %s",
					round, acceptances, time.Since(start).Truncate(time.Second)),
				false)
		}

		acceptances++
		bestEver, bestEverStats, bestEverAvg = applyAcceptedMutation(cfg, round, verbose, tempLabel,
			idx, len(mutations), mutations[idx], d, dStats, avg, currentAvg, bestEver, bestEverStats, bestEverAvg)
		current = d
		currentAvg = avg
		temperature = coolDown(temperature, cfg.tempDecay, cfg.minTemp)
	}
}

// buildRoundMutations enumerates every single-card / weapon mutation and shuffles the
// result so first-improvement climb and SA acceptance don't concentrate on a fixed slice.
func buildRoundMutations(cfg annealConfig, rng *rand.Rand, current *deck.Deck) []deck.Mutation {
	mutations := deck.AllMutations(current, cfg.maxCopies, cfg.pairMutations, registry.Registry{})
	rng.Shuffle(len(mutations), func(i, j int) {
		mutations[i], mutations[j] = mutations[j], mutations[i]
	})
	return mutations
}

// formatTempLabel is the " (T=…)" suffix; empty string when temperature is 0.
func formatTempLabel(temperature float64) string {
	if temperature <= 0 {
		return ""
	}
	return fmt.Sprintf(" (T=%.4f)", temperature)
}

// applyAcceptedMutation runs the after-success branch: logs the acceptance, logs a new
// all-time best in non-verbose mode, persists the deck when avg exceeds bestEverAvg, and
// returns the possibly-updated bestEver / bestEverAvg.
func applyAcceptedMutation(cfg annealConfig, round int, verbose bool, tempLabel string,
	idx, total int, mut deck.Mutation, d *deck.Deck, dStats deck.Stats, avg, currentAvg float64,
	bestEver *deck.Deck, bestEverStats deck.Stats, bestEverAvg float64) (*deck.Deck, deck.Stats, float64) {
	verb := "improvement"
	if avg <= currentAvg {
		verb = "annealing step"
	}
	if verbose {
		fmt.Fprintf(os.Stderr, "\r[round %d] %s at %d/%d: deep %.3f vs %.3f (%s)%s       \n",
			round, verb, idx+1, total, avg, currentAvg, mut.Description, tempLabel)
	}
	if avg <= bestEverAvg {
		return bestEver, bestEverStats, bestEverAvg
	}
	if !verbose {
		// Surface new all-time bests in non-verbose mode so long reanneal sessions show
		// forward motion. \r + padding overwrites the ticker line; the trailing \n promotes
		// the entry above the next round's ticker.
		fmt.Fprintf(os.Stderr, "\r[round %d] new best %.3f (was %.3f, +%.3f)%s                                \n",
			round, avg, bestEverAvg, avg-bestEverAvg, tempLabel)
	}
	if err := writeDeck(d, dStats, cfg.outPath); err != nil {
		die("%v", err)
	}
	return d, dStats, avg
}

// finishAnnealRun emits the terminal status line, optionally prints the best-ever deck
// listing, and builds the annealResult that drives the command's exit code and summary.
func finishAnnealRun(cfg annealConfig, bestEver *deck.Deck, bestEverStats deck.Stats, bestEverAvg, startingAvg float64,
	statusLine string, aborted bool) annealResult {
	fmt.Fprintln(os.Stderr, "\n"+statusLine)
	fmt.Println()
	if shouldPrintFinalDeck(cfg.startTemp, bestEverAvg, startingAvg) {
		printBestDeck(bestEver, bestEverStats)
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
// return the same DeckStats shape; DeckStats.Runs reflects the actual shuffle count so
// the next prepareBaseline call's "already evaluated" check still works (an adaptive run
// may finish below the cap and prompt a re-evaluation next session, which is fine —
// adaptive runs are cheap). Routes through evaluateParallel so the once-per-session
// baseline benefits from the same DefaultWorkers fan-out as iterate's per-mutation evals.
func baselineEvaluate(d *deck.Deck, cfg annealConfig, rng *rand.Rand) deck.Stats {
	shuffles := cfg.shuffles
	if cfg.adaptive {
		shuffles = -1
	}
	stats, _ := evaluateParallel(d, shuffles, cfg.minImprovement, cfg.matchup, rng)
	return stats
}

// prepareBaseline returns the starting deck for the hill climb with its baseline avg.
// Four cases: no deck on disk (generate random + evaluate); loaded deck under the current
// shuffle budget (re-evaluate for an apples-to-apples baseline); -reevaluate set (force
// re-evaluation even if the run count already matches); or deck already evaluated at the
// budget (use as-is). File exists but doesn't parse → die loudly rather than silently
// overwrite a corrupt checkpoint.
func prepareBaseline(cfg annealConfig, rng *rand.Rand) (*deck.Deck, deck.Stats, float64) {
	best, bestStats, err := loadExisting(cfg.outPath)
	if err != nil {
		die("%v", err)
	}
	if best == nil {
		fmt.Fprintf(os.Stderr, "no deck at %s; generating a random starting deck\n", cfg.outPath)
		best = deck.Random(heroes.Viserai, cfg.deckSize, cfg.maxCopies, rng, registry.Registry{})
		bestStats = baselineEvaluate(best, cfg, rng)
		bestAvg := bestStats.Mean()
		if err := writeDeck(best, bestStats, cfg.outPath); err != nil {
			die("%v", err)
		}
		fmt.Printf("Starting deck avg %.3f, saved to %s\n", bestAvg, cfg.outPath)
		maybePrintBaselineCards(cfg, best)
		return best, bestStats, bestAvg
	}
	bestAvg := bestStats.Mean()
	// Re-evaluate when the saved deck was scored at fewer shuffles than the current budget,
	// or when -reevaluate forces it. Adaptive runs always take the re-evaluate path because
	// the recorded Runs reflects whatever count the previous adaptive run terminated at,
	// which carries no precision guarantee for the new run's target.
	needReeval := cfg.adaptive || cfg.reevaluate || bestStats.Runs < cfg.shuffles
	if needReeval {
		var freshStats deck.Stats
		best, freshStats, bestAvg = reevaluateBaseline(cfg, rng, best, bestStats, bestAvg)
		maybePrintBaselineCards(cfg, best)
		return best, freshStats, bestAvg
	}
	fmt.Printf("Loaded best deck (avg %.3f) from %s\n", bestAvg, cfg.outPath)
	maybePrintBaselineCards(cfg, best)
	return best, bestStats, bestAvg
}

// reevaluateBaseline rebuilds the loaded deck against the current shuffle budget and writes
// the refreshed stats back to disk. Picks an explanatory reason label (-reevaluate forced or
// stale shuffle count), reconstructs the deck (Sideboard and Equipment preserved), runs
// baselineEvaluate, and persists the result. Returns the rebuilt deck, its fresh stats, and
// avg.
func reevaluateBaseline(cfg annealConfig, rng *rand.Rand, loaded *deck.Deck, loadedStats deck.Stats, savedAvg float64) (*deck.Deck, deck.Stats, float64) {
	reason := fmt.Sprintf("from %d shuffles", loadedStats.Runs)
	if cfg.reevaluate && loadedStats.Runs >= cfg.shuffles {
		reason = "-reevaluate forced"
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
	rebuilt := loaded.Copy()
	rebuilt.Sideboard = sideboard
	rebuilt.Equipment = equipment
	freshStats := baselineEvaluate(rebuilt, cfg, rng)
	freshAvg := freshStats.Mean()
	if err := writeDeck(rebuilt, freshStats, cfg.outPath); err != nil {
		die("%v", err)
	}
	// Show saved→current so the delta from any simulation-logic drift is visible at a glance,
	// instead of the user guessing which of the two printed numbers is the fresh one.
	fmt.Printf("Re-evaluated baseline: %.3f → %.3f, saved to %s\n", savedAvg, freshAvg, cfg.outPath)
	return rebuilt, freshStats, freshAvg
}

// maybePrintBaselineCards emits the startup card-list dump unless -quiet-load suppressed it. The
// leading blank line is part of the listing block, so it's also gated — otherwise -quiet-load
// would leave a lone empty line hanging after the baseline avg summary.
func maybePrintBaselineCards(cfg annealConfig, d *deck.Deck) {
	if cfg.quietLoad {
		return
	}
	fmt.Println()
	printCardList(d)
}

// watchStdinForAbort spawns a background goroutine that calls cancel() on the first keypress.
// EOF / closed stdin isn't an abort (so anneal doesn't exit immediately on non-TTY stdin); only
// a successful read of at least one byte counts. Cancellation propagates into RunMutationRound
// so an abort takes effect mid-round.
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
