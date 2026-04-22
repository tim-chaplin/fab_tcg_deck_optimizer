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

func runIterate(cfg config) float64 {
	rng := rand.New(rand.NewSource(cfg.seed))

	current, currentAvg := prepareBaseline(cfg, rng)
	// All-time best tracks the highest-avg deck seen since runIterate started. The saved JSON
	// mirrors this — simulated annealing intentionally walks through worse states to escape
	// local maxima, but the on-disk artifact should always reflect the peak reached so far.
	bestEver := current
	bestEverAvg := currentAvg
	fmt.Println("Press Enter to abort.")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	watchStdinForAbort(cancel)

	temperature := cfg.startTemp
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
		tempLabel := ""
		if temperature > 0 {
			tempLabel = fmt.Sprintf(" (T=%.4f)", temperature)
		}
		fmt.Fprintf(os.Stderr, "\n[round %d] evaluating %d mutations of avg %.3f%s (best ever %.3f)\n",
			round, len(mutations), currentAvg, tempLabel, bestEverAvg)

		var tested, deepsDone atomic.Int64
		stopTicker := startRoundTicker(round, len(mutations), start, &tested, &deepsDone)
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
			printBestDeck(bestEver)
			return bestEverAvg
		}
		if !found {
			// A full round with zero acceptances means every mutation — including the
			// probabilistically-accepted worse ones — failed the gate. At any T > 0 with
			// thousands of mutations this is vanishingly unlikely unless we've genuinely
			// converged, so treat it as a local maximum regardless of the current temperature.
			fmt.Fprintf(os.Stderr, "\nLocal maximum reached after %d rounds / %d acceptances in %s\n",
				round, acceptances, time.Since(start).Truncate(time.Second))
			fmt.Println()
			printBestDeck(bestEver)
			return bestEverAvg
		}

		acceptances++
		mut := mutations[idx]
		verb := "improvement"
		if avg <= currentAvg {
			verb = "annealing step"
		}
		fmt.Fprintf(os.Stderr, "\r[round %d] %s at %d/%d: deep %.3f vs %.3f (%s)%s       \n",
			round, verb, idx+1, len(mutations), avg, currentAvg, mut.Description, tempLabel)
		current = d
		currentAvg = avg
		if avg > bestEverAvg {
			bestEver = d
			bestEverAvg = avg
			_ = writeDeck(bestEver, cfg.outPath)
		}
		temperature = coolDown(temperature, cfg.tempDecay, cfg.minTemp)
	}
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

// startRoundTicker launches a 500ms ticker that renders the round's shallow-screen and
// deep-confirm counts to stderr on a \r-terminated line so the user sees the worker pool moving
// during long rounds. Returns a stop function the caller must call when the round finishes.
func startRoundTicker(round, total int, start time.Time, tested, deepsDone *atomic.Int64) func() {
	done := make(chan struct{})
	go func() {
		t := time.NewTicker(500 * time.Millisecond)
		defer t.Stop()
		for {
			select {
			case <-done:
				return
			case <-t.C:
				fmt.Fprintf(os.Stderr, "\r[round %d] tested %d/%d (%d deep confirms, %s elapsed)        ",
					round, tested.Load(), total, deepsDone.Load(), time.Since(start).Truncate(time.Second))
			}
		}
	}()
	return func() { close(done) }
}
