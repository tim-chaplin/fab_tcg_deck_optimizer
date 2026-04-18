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

	best, bestAvg := loadExisting(cfg.outPath)
	if best == nil {
		// No starting point on disk — bootstrap with a single random deck, evaluated at the same
		// -deep-shuffles depth the hill climb uses so bestAvg is comparable to future mutations.
		fmt.Fprintf(os.Stderr, "no deck at %s; generating a random starting deck\n", cfg.outPath)
		best = deck.Random(hero.Viserai{}, cfg.deckSize, cfg.maxCopies, rng)
		stats := best.Evaluate(cfg.deepShuffles, cfg.incoming, rng)
		bestAvg = stats.Avg()
		_ = writeDeck(best, cfg.outPath)
		fmt.Printf("Starting deck avg %.3f, saved to %s\n", bestAvg, cfg.outPath)
	} else if best.Stats.Runs < cfg.deepShuffles {
		// The loaded baseline was evaluated with a shallower sample than the hill climb will
		// use for each mutation. Re-evaluate at the current depth so bestAvg is apples-to-apples
		// with the mutation scores the climber will compare against.
		fmt.Printf("Loaded best deck (avg %.3f from %d shuffles); re-evaluating at %d shuffles for an apples-to-apples baseline\n",
			bestAvg, best.Stats.Runs, cfg.deepShuffles)
		best = deck.New(best.Hero, best.Weapons, best.Cards)
		stats := best.Evaluate(cfg.deepShuffles, cfg.incoming, rng)
		bestAvg = stats.Avg()
		_ = writeDeck(best, cfg.outPath)
		fmt.Printf("Re-evaluated baseline avg %.3f, saved to %s\n", bestAvg, cfg.outPath)
	} else {
		fmt.Printf("Loaded best deck (avg %.3f) from %s\n", bestAvg, cfg.outPath)
	}
	fmt.Println("Press Enter to abort.")

	// Abort signal: background goroutine reads stdin and cancels ctx on the first keypress. EOF
	// / closed stdin isn't an abort signal (otherwise iterate would exit immediately when stdin
	// isn't a TTY) — only an actual read of at least one byte counts. Cancelling ctx propagates
	// into IterateParallel so an abort takes effect mid-round rather than waiting for the current
	// round to finish.
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() {
		buf := make([]byte, 1)
		if n, err := os.Stdin.Read(buf); err == nil && n > 0 {
			cancel()
		}
	}()

	// Deterministic hill-climb: enumerate every single-slot mutation of the current best. On the
	// first mutation that scores higher, adopt it and restart enumeration from the new best. If
	// we exhaust every mutation with no improvement, we're at a local maximum.
	round := 0
	improvements := 0
	start := time.Now()
	for {
		round++
		mutations := deck.AllMutations(best, cfg.maxCopies)
		fmt.Fprintf(os.Stderr, "\n[round %d] evaluating %d mutations of avg %.3f\n",
			round, len(mutations), bestAvg)

		// deck.IterateParallel runs a persistent worker pool where each worker does both the
		// shallow screen AND the deep confirmation for any mutation it flags — so the expensive
		// deep pass parallelises across workers instead of bottlenecking on the main goroutine.
		// First worker to land a confirmed improvement short-circuits the rest. The ticker off
		// the two atomics renders "tested X/total (Y deep confirms)" so the user sees BOTH the
		// shallow sweep and the deep-confirm churn advancing.
		var tested, deepsDone atomic.Int64
		tickerDone := make(chan struct{})
		go func() {
			t := time.NewTicker(500 * time.Millisecond)
			defer t.Stop()
			for {
				select {
				case <-tickerDone:
					return
				case <-t.C:
					fmt.Fprintf(os.Stderr, "\r[round %d] tested %d/%d (%d deep confirms, %s elapsed)        ",
						round, tested.Load(), len(mutations), deepsDone.Load(), time.Since(start).Truncate(time.Second))
				}
			}
		}()
		d, avg, idx, found := deck.IterateParallel(
			ctx, mutations, bestAvg, cfg.shallowShuffles, cfg.deepShuffles, cfg.incoming, 0,
			rng.Int63(), rng, &tested, &deepsDone,
		)
		close(tickerDone)

		if ctx.Err() != nil {
			fmt.Fprintf(os.Stderr, "\nAborted mid-round after %d rounds / %d improvements in %s\n",
				round, improvements, time.Since(start).Truncate(time.Second))
			fmt.Println()
			printBestDeck(best)
			return
		}
		improved := false
		if found {
			improvements++
			mut := mutations[idx]
			fmt.Fprintf(os.Stderr, "\r[round %d] improvement at %d/%d: deep %.3f beats %.3f (%s), restarting        \n",
				round, idx+1, len(mutations), avg, bestAvg, mut.Description)
			bestAvg = avg
			best = d
			_ = writeDeck(best, cfg.outPath)
			improved = true
		}

		if !improved {
			fmt.Fprintf(os.Stderr, "\nLocal maximum reached after %d rounds / %d improvements in %s\n",
				round, improvements, time.Since(start).Truncate(time.Second))
			fmt.Println()
			printBestDeck(best)
			return
		}
	}
}
