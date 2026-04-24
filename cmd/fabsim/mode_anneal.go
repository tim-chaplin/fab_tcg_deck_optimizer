package main

import (
	"context"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"sync/atomic"
	"time"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/deck"
	fmtpkg "github.com/tim-chaplin/fab-deck-optimizer/internal/format"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/hand"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/hero"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/mydecks"
)

// annealConfig bundles the knobs runAnneal needs. Built by runAnnealCmd from its flag.FlagSet.
type annealConfig struct {
	shallowShuffles int
	deepShuffles    int
	incoming        int
	deckSize        int
	maxCopies       int
	seed            int64
	outPath         string
	format          fmtpkg.Format
	debug           bool
	reevaluate      bool
	// startTemp / tempDecay / minTemp are the simulated-annealing knobs. startTemp of 0
	// degenerates to the classical hill-climb (strict > baseline acceptance).
	startTemp float64
	tempDecay float64
	minTemp   float64
	// quietLoad suppresses the baseline card-list dump in prepareBaseline. Set by wrapper
	// scripts that re-invoke anneal repeatedly on the same deck; the listing is unchanging
	// noise after the first pass.
	quietLoad bool
}

// legalFilter returns the card-pool predicate for this run's format. anneal always runs under a
// format, so this is non-nil; the deck package accepts nil for "no filtering" generally.
func (c annealConfig) legalFilter() func(card.Card) bool {
	return c.format.IsLegal
}

// defaultDeckNameFor returns the deck name when -deck isn't supplied, keyed by hero, format, and
// -incoming. Different regimes produce different optimal decks, so each gets its own file to
// avoid hill-climbing one regime's best under another regime's objective.
func defaultDeckNameFor(h hero.Hero, f fmtpkg.Format, incoming int) string {
	return fmt.Sprintf("%s_%s_%d_incoming", strings.ToLower(h.Name()), f, incoming)
}

// runAnnealCmd parses anneal's flags from args and dispatches to runAnneal. Every anneal-only
// knob lives on this FlagSet so `fabsim anneal -help` shows exactly the flags that apply here.
func runAnnealCmd(args []string) {
	fs := flag.NewFlagSet("anneal", flag.ExitOnError)
	deckName := fs.String("deck", "", "deck name; resolved to mydecks/<name>.json (\".json\" suffix optional). Defaults to <hero>_<format>_<incoming>_incoming so different (hero, format, -incoming) regimes keep separate deck files. When the named deck exists, anneal resumes from it as a checkpoint.")
	shallowShuffles := fs.Int("shallow-shuffles", 100, "shuffles per deck used to screen mutations before deep confirmation")
	deepShuffles := fs.Int("deep-shuffles", 10000, "shuffles per deck used to confirm improvements and to baseline loaded decks")
	incoming := fs.Int("incoming", 0, "opponent damage per turn (required — different values produce different optimal decks, so this is explicit rather than defaulted)")
	deckSize := fs.Int("deck-size", 40, "number of cards per deck")
	maxCopies := fs.Int("max-copies", defaultMaxCopies, "maximum copies of any single card printing per deck")
	seed := fs.Int64("seed", time.Now().UnixNano(), "RNG seed")
	formatFlag := fs.String("format", string(fmtpkg.SilverAge), "constructed format whose banlist restricts the card pool during search (only \"silver_age\" is supported today)")
	debug := fs.Bool("debug", false, "emit extra diagnostic output (e.g. memo cache size between rounds)")
	reevaluate := fs.Bool("reevaluate", false, "force re-evaluation of the loaded deck's baseline avg, even if its prior run count already matches -deep-shuffles. Use after adjusting modelling assumptions or fixing bugs that may have shifted the deck's true score.")
	finalize := fs.Bool("finalize", false, "high-precision pass — overrides -shallow-shuffles to 10000 and -deep-shuffles to 100000. Use on a deck that's already converged to squeeze out the remaining sub-percent improvements.")
	startTemp := fs.Float64("start-temp", 0, "simulated-annealing starting temperature. 0 (default) runs a pure hill climb. Higher values probabilistically accept worse mutations early; acceptance probability is exp((avg - baseline) / T). Good starting range is ~0.05–0.5 given typical Value units.")
	tempDecay := fs.Float64("temp-decay", 0.95, "multiplicative cooling per acceptance — T ← T × decay, floored at -min-temp. Unused when -start-temp is 0.")
	minTemp := fs.Float64("min-temp", 0, "minimum temperature. Once T reaches this floor the climb becomes greedy until a local maximum is found. 0 disables annealing in the converged tail.")
	quietLoad := fs.Bool("quiet-load", false, "skip the baseline card-list dump at startup. Intended for wrapper scripts (e.g. anneal-reanneal.ps1) that re-invoke anneal many times on the same deck — the listing never changes pass-to-pass and floods the log.")
	_ = parseFlagsAnywhere(fs, args)
	if fs.NArg() > 0 {
		die("anneal: unexpected positional argument(s): %v (did you mean -deck %s?)", fs.Args(), fs.Args()[0])
	}
	requireFlag(fs, "anneal", "incoming")

	fmtValue, err := fmtpkg.Parse(*formatFlag)
	if err != nil {
		die("%v", err)
	}

	// -finalize is a shuffle-count shorthand, so apply it as a post-parse override rather than
	// threading it into annealConfig.
	if *finalize {
		*shallowShuffles = 10000
		*deepShuffles = 100000
	}

	name := *deckName
	if name == "" {
		name = defaultDeckNameFor(hero.Viserai{}, fmtValue, *incoming)
	}
	outPath, err := mydecks.Path(name)
	if err != nil {
		die("%v", err)
	}

	cfg := annealConfig{
		shallowShuffles: *shallowShuffles,
		deepShuffles:    *deepShuffles,
		incoming:        *incoming,
		deckSize:        *deckSize,
		maxCopies:       *maxCopies,
		seed:            *seed,
		outPath:         outPath,
		format:          fmtValue,
		debug:           *debug,
		reevaluate:      *reevaluate,
		startTemp:       *startTemp,
		tempDecay:       *tempDecay,
		minTemp:         *minTemp,
		quietLoad:       *quietLoad,
	}

	// Print the session-level delta (starting best vs final best) on any exit path, then
	// surface abort via a non-zero exit so wrapper scripts (anneal-reanneal.ps1 et al.) can
	// tell Enter-initiated termination from natural convergence and stop looping.
	res := runAnneal(cfg)
	fmt.Fprintf(os.Stderr, "\nSession summary: avg %.3f → %.3f (%+.3f)\n",
		res.startingAvg, res.bestEverAvg, res.bestEverAvg-res.startingAvg)
	if res.aborted {
		os.Exit(130)
	}
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

	ctx, cancel := context.WithCancel(context.Background())
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
		mutations := buildRoundMutations(cfg, rng, current, round)
		tempLabel := formatTempLabel(temperature)
		if verbose {
			fmt.Fprintf(os.Stderr, "\n[round %d] evaluating %d mutations of avg %.3f%s (best ever %.3f)\n",
				round, len(mutations), currentAvg, tempLabel, bestEverAvg)
		}

		var tested, deepsDone atomic.Int64
		stopTicker := startRoundTicker(round, len(mutations), start, &tested, &deepsDone,
			temperature, currentAvg, bestEverAvg)
		d, avg, idx, found := deck.IterateParallel(
			ctx, mutations, currentAvg, temperature,
			cfg.shallowShuffles, cfg.deepShuffles, cfg.incoming, 0,
			rng.Int63(), &tested, &deepsDone,
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

// buildRoundMutations produces the per-round mutation list: clears the shared hand memo (so the
// next round starts with a bounded cache), enumerates every single-card/weapon mutation, and
// under annealing shuffles the order so probabilistic acceptances aren't concentrated on the
// weakest card. Also emits the debug memo-size line when -debug is set.
func buildRoundMutations(cfg annealConfig, rng *rand.Rand, current *deck.Deck, round int) []deck.Mutation {
	// Drop the shared hand memo between rounds. Within a round the memo is load-bearing
	// (same hand shapes recur across thousands of shuffles), but cross-round hit rate is
	// near zero and unbounded growth would OOM long hill-climbs.
	if cfg.debug {
		fmt.Fprintf(os.Stderr, "[memo] clearing %d entries before round %d\n", hand.MemoLen(), round)
	}
	hand.ClearMemo()
	mutations := deck.AllMutations(current, cfg.maxCopies, cfg.legalFilter())
	if cfg.startTemp > 0 {
		// AllMutations sorts weakest-card-first so a first-found classical climb tries the
		// highest-expected-gain swaps early. Under annealing that same bias means the
		// probabilistic acceptances disproportionately hit mutations against the weakest
		// card in the deck — shrinking the slice of the solution space the walk actually
		// explores. Shuffling each round gives every mutation an even shot at being the
		// first one accepted at the current temperature.
		rng.Shuffle(len(mutations), func(i, j int) {
			mutations[i], mutations[j] = mutations[j], mutations[i]
		})
	}
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
	idx, total int, mut deck.Mutation, d *deck.Deck, avg, currentAvg float64,
	bestEver *deck.Deck, bestEverAvg float64) (*deck.Deck, float64) {
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
func finishAnnealRun(cfg annealConfig, bestEver *deck.Deck, bestEverAvg, startingAvg float64,
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

// prepareBaseline returns the starting deck for the hill climb with its deep-shuffles avg.
// Four cases: no deck on disk (generate random + evaluate); loaded deck under deepShuffles
// (re-evaluate for an apples-to-apples baseline); -reevaluate set (force re-evaluation even
// if the run count already matches); or deck already deep-evaluated (use as-is). File
// exists but doesn't parse → die loudly rather than silently overwrite a corrupt checkpoint.
//
// A loaded deck that contains card.NotImplemented copies (e.g. a pre-tag deck recovered
// from disk) is sanitized before any of the above branches: the tagged slots are replaced
// with random legal picks and the run always takes the re-evaluate path so the baseline
// reflects the new card list.
func prepareBaseline(cfg annealConfig, rng *rand.Rand) (*deck.Deck, float64) {
	best, bestAvg, err := loadExisting(cfg.outPath)
	if err != nil {
		die("%v", err)
	}
	if best == nil {
		fmt.Fprintf(os.Stderr, "no deck at %s; generating a random starting deck\n", cfg.outPath)
		best = deck.Random(hero.Viserai{}, cfg.deckSize, cfg.maxCopies, rng, cfg.legalFilter())
		bestAvg = best.Evaluate(cfg.deepShuffles, cfg.incoming, rng).Mean()
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
	if cfg.reevaluate || best.Stats.Runs < cfg.deepShuffles {
		reason := fmt.Sprintf("from %d shuffles", best.Stats.Runs)
		if cfg.reevaluate && best.Stats.Runs >= cfg.deepShuffles {
			reason = "-reevaluate forced"
		}
		if len(sanitized) > 0 {
			reason = fmt.Sprintf("%d NotImplemented card(s) replaced", len(sanitized))
		}
		// Label the loaded number "saved avg" so it can't be mistaken for the re-evaluated
		// score. Decks scored under older simulation logic can have saved avgs that diverge
		// substantially from what today's simulator produces.
		fmt.Printf("Loaded best deck (saved avg %.3f, %s); re-evaluating at %d shuffles for an apples-to-apples baseline\n",
			bestAvg, reason, cfg.deepShuffles)
		savedAvg := bestAvg
		// Sideboard and Equipment are user-managed and don't feed the sim — preserve
		// them across the stats reset so the re-evaluated deck writes back unchanged.
		sideboard := best.Sideboard
		equipment := best.Equipment
		best = deck.New(best.Hero, best.Weapons, best.Cards)
		best.Sideboard = sideboard
		best.Equipment = equipment
		bestAvg = best.Evaluate(cfg.deepShuffles, cfg.incoming, rng).Mean()
		if err := writeDeck(best, cfg.outPath); err != nil {
			die("%v", err)
		}
		// Show saved→current so the delta from any simulation-logic drift is visible at a
		// glance, instead of the user guessing which of the two printed numbers is the fresh
		// one.
		fmt.Printf("Re-evaluated baseline: %.3f → %.3f, saved to %s\n", savedAvg, bestAvg, cfg.outPath)
		maybePrintBaselineCards(cfg, best)
		return best, bestAvg
	}
	fmt.Printf("Loaded best deck (avg %.3f) from %s\n", bestAvg, cfg.outPath)
	maybePrintBaselineCards(cfg, best)
	return best, bestAvg
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
// rich enough to track the walk is what makes silent mode bearable. Returns a stop function
// the caller must call when the round finishes.
func startRoundTicker(round, total int, start time.Time, tested, deepsDone *atomic.Int64, temperature, currentAvg, bestEverAvg float64) func() {
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
				fmt.Fprintf(os.Stderr, "\r[round %d] tested %d/%d  cur %.3f  best %.3f%s  %s elapsed        ",
					round, tested.Load(), total, currentAvg, bestEverAvg, tempLabel,
					time.Since(start).Truncate(time.Second))
			}
		}
	}()
	return func() { close(done) }
}
