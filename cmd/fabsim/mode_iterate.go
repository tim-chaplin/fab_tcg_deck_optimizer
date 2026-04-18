package main

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"sync/atomic"
	"time"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/deck"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/hero"
)

func runIterate(cfg config) {
	rng := rand.New(rand.NewSource(cfg.seed))

	best, bestAvg := prepareBaseline(cfg, rng)
	fmt.Println("Press Enter to abort.")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	watchStdinForAbort(cancel)

	// Deterministic hill-climb: enumerate every single-slot mutation of the current best. On the
	// first mutation that scores higher, adopt it and restart enumeration from the new best. If
	// we exhaust every mutation with no improvement, we're at a local maximum.
	round := 0
	improvements := 0
	start := time.Now()
	for {
		round++
		mutations := deck.AllMutations(best, cfg.maxCopies, cfg.legalFilter())
		fmt.Fprintf(os.Stderr, "\n[round %d] evaluating %d mutations of avg %.3f\n",
			round, len(mutations), bestAvg)

		var tested, deepsDone atomic.Int64
		stopTicker := startRoundTicker(round, len(mutations), start, &tested, &deepsDone)
		d, avg, idx, found := deck.IterateParallel(
			ctx, mutations, bestAvg, cfg.shallowShuffles, cfg.deepShuffles, cfg.incoming, 0,
			rng.Int63(), &tested, &deepsDone,
		)
		stopTicker()

		if ctx.Err() != nil {
			fmt.Fprintf(os.Stderr, "\nAborted mid-round after %d rounds / %d improvements in %s\n",
				round, improvements, time.Since(start).Truncate(time.Second))
			fmt.Println()
			printBestDeck(best)
			return
		}
		if !found {
			fmt.Fprintf(os.Stderr, "\nLocal maximum reached after %d rounds / %d improvements in %s\n",
				round, improvements, time.Since(start).Truncate(time.Second))
			fmt.Println()
			printBestDeck(best)
			return
		}

		improvements++
		mut := mutations[idx]
		fmt.Fprintf(os.Stderr, "\r[round %d] improvement at %d/%d: deep %.3f beats %.3f (%s), restarting        \n",
			round, idx+1, len(mutations), avg, bestAvg, mut.Description)
		bestAvg = avg
		best = d
		_ = writeDeck(best, cfg.outPath)
	}
}

// prepareBaseline returns the starting deck for the hill climb along with its deep-shuffles avg.
// Three cases: no deck on disk (generate random + evaluate), loaded deck evaluated at fewer than
// deepShuffles (re-evaluate at current depth for apples-to-apples), or loaded deck already deep-
// evaluated (use as-is).
func prepareBaseline(cfg config, rng *rand.Rand) (*deck.Deck, float64) {
	best, bestAvg := loadExisting(cfg.outPath)
	if best == nil {
		fmt.Fprintf(os.Stderr, "no deck at %s; generating a random starting deck\n", cfg.outPath)
		best = deck.Random(hero.Viserai{}, cfg.deckSize, cfg.maxCopies, rng, cfg.legalFilter())
		bestAvg = best.Evaluate(cfg.deepShuffles, cfg.incoming, rng).Avg()
		_ = writeDeck(best, cfg.outPath)
		fmt.Printf("Starting deck avg %.3f, saved to %s\n", bestAvg, cfg.outPath)
		return best, bestAvg
	}
	if best.Stats.Runs < cfg.deepShuffles {
		fmt.Printf("Loaded best deck (avg %.3f from %d shuffles); re-evaluating at %d shuffles for an apples-to-apples baseline\n",
			bestAvg, best.Stats.Runs, cfg.deepShuffles)
		best = deck.New(best.Hero, best.Weapons, best.Cards)
		bestAvg = best.Evaluate(cfg.deepShuffles, cfg.incoming, rng).Avg()
		_ = writeDeck(best, cfg.outPath)
		fmt.Printf("Re-evaluated baseline avg %.3f, saved to %s\n", bestAvg, cfg.outPath)
		return best, bestAvg
	}
	fmt.Printf("Loaded best deck (avg %.3f) from %s\n", bestAvg, cfg.outPath)
	return best, bestAvg
}

// watchStdinForAbort spawns a background goroutine that reads stdin and calls cancel() on the
// first keypress. EOF / closed stdin isn't an abort signal (otherwise iterate would exit
// immediately when stdin isn't a TTY) — only an actual read of at least one byte counts.
// Cancelling the context propagates into IterateParallel so an abort takes effect mid-round
// rather than waiting for the current round to finish.
func watchStdinForAbort(cancel context.CancelFunc) {
	go func() {
		buf := make([]byte, 1)
		if n, err := os.Stdin.Read(buf); err == nil && n > 0 {
			cancel()
		}
	}()
}

// startRoundTicker launches a 500ms ticker that renders the current round's progress — both the
// shallow-screen count and the deep-confirm count — to stderr on a CR-terminated line, so the
// user sees the worker pool moving even during long rounds. Returns a stop function the caller
// must call once the round finishes.
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
