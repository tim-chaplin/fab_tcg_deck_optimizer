package main

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"sync/atomic"
	"time"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/deck"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/hand"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/hero"
)

// annealResult carries the outcome of a single runAnneal pass. aborted is true when the
// user hit Enter (stdin watcher fired) — callers propagate this up to main so the process
// exits non-zero and wrapper scripts like anneal-reanneal.ps1 stop their outer loop instead
// of immediately launching another pass.
type annealResult struct {
	bestEverAvg float64
	startingAvg float64
	aborted     bool
}

func runAnneal(cfg config) annealResult {
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
		tempLabel := ""
		if temperature > 0 {
			tempLabel = fmt.Sprintf(" (T=%.4f)", temperature)
		}
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
			fmt.Fprintf(os.Stderr, "\nAborted mid-round after %d rounds / %d acceptances in %s\n",
				round, acceptances, time.Since(start).Truncate(time.Second))
			fmt.Println()
			if shouldPrintFinalDeck(cfg.startTemp, bestEverAvg, startingAvg) {
				printBestDeck(bestEver)
			}
			return annealResult{bestEverAvg: bestEverAvg, startingAvg: startingAvg, aborted: true}
		}
		if !found {
			// A full round with zero acceptances means every mutation — including the
			// probabilistically-accepted worse ones — failed the gate. At any T > 0 with
			// thousands of mutations this is vanishingly unlikely unless we've genuinely
			// converged, so treat it as a local maximum regardless of the current temperature.
			fmt.Fprintf(os.Stderr, "\nLocal maximum reached after %d rounds / %d acceptances in %s\n",
				round, acceptances, time.Since(start).Truncate(time.Second))
			fmt.Println()
			if shouldPrintFinalDeck(cfg.startTemp, bestEverAvg, startingAvg) {
				printBestDeck(bestEver)
			}
			return annealResult{bestEverAvg: bestEverAvg, startingAvg: startingAvg}
		}

		acceptances++
		mut := mutations[idx]
		verb := "improvement"
		if avg <= currentAvg {
			verb = "annealing step"
		}
		if verbose {
			fmt.Fprintf(os.Stderr, "\r[round %d] %s at %d/%d: deep %.3f vs %.3f (%s)%s       \n",
				round, verb, idx+1, len(mutations), avg, currentAvg, mut.Description, tempLabel)
		}
		current = d
		currentAvg = avg
		if avg > bestEverAvg {
			if !verbose {
				// Surface every new all-time best in non-verbose annealing mode so long
				// reanneal sessions still show visible forward motion without the per-round
				// spam. \r + trailing padding overwrites the ticker line cleanly; the \n at
				// the end promotes this to a persistent entry above the next round's ticker.
				fmt.Fprintf(os.Stderr, "\r[round %d] new best %.3f (was %.3f, +%.3f)%s                                \n",
					round, avg, bestEverAvg, avg-bestEverAvg, tempLabel)
			}
			bestEver = d
			bestEverAvg = avg
			_ = writeDeck(bestEver, cfg.outPath)
		}
		temperature = coolDown(temperature, cfg.tempDecay, cfg.minTemp)
	}
}

// shouldPrintFinalDeck decides whether to dump the full deck listing at the end of a run.
// Annealing sessions that exit with no net improvement over the starting deck (a common
// outcome for long reanneal loops that are just probing) would otherwise reprint the same
// cards every session and bury the session-summary line in noise. Classical mode keeps
// printing regardless — a no-improvement run there is a single round and the listing is the
// user's confirmation of what was evaluated.
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

// prepareBaseline returns the starting deck for the hill climb with its deep-shuffles avg. Four
// cases: no deck on disk (generate random + evaluate), loaded deck under deepShuffles
// (re-evaluate for an apples-to-apples baseline), -reevaluate set (force re-evaluation even if
// the run count already matches — for when modelling assumptions changed), or deck already
// deep-evaluated (use as-is).
func prepareBaseline(cfg config, rng *rand.Rand) (*deck.Deck, float64) {
	best, bestAvg := loadExisting(cfg.outPath)
	if best == nil {
		fmt.Fprintf(os.Stderr, "no deck at %s; generating a random starting deck\n", cfg.outPath)
		best = deck.Random(hero.Viserai{}, cfg.deckSize, cfg.maxCopies, rng, cfg.legalFilter())
		bestAvg = best.Evaluate(cfg.deepShuffles, cfg.incoming, rng).Mean()
		_ = writeDeck(best, cfg.outPath)
		fmt.Printf("Starting deck avg %.3f, saved to %s\n", bestAvg, cfg.outPath)
		fmt.Println()
		printCardList(best)
		return best, bestAvg
	}
	if cfg.reevaluate || best.Stats.Runs < cfg.deepShuffles {
		reason := fmt.Sprintf("from %d shuffles", best.Stats.Runs)
		if cfg.reevaluate && best.Stats.Runs >= cfg.deepShuffles {
			reason = "-reevaluate forced"
		}
		fmt.Printf("Loaded best deck (avg %.3f %s); re-evaluating at %d shuffles for an apples-to-apples baseline\n",
			bestAvg, reason, cfg.deepShuffles)
		best = deck.New(best.Hero, best.Weapons, best.Cards)
		bestAvg = best.Evaluate(cfg.deepShuffles, cfg.incoming, rng).Mean()
		_ = writeDeck(best, cfg.outPath)
		fmt.Printf("Re-evaluated baseline avg %.3f, saved to %s\n", bestAvg, cfg.outPath)
		fmt.Println()
		printCardList(best)
		return best, bestAvg
	}
	fmt.Printf("Loaded best deck (avg %.3f) from %s\n", bestAvg, cfg.outPath)
	fmt.Println()
	printCardList(best)
	return best, bestAvg
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
