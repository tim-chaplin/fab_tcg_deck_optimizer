package main

import (
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

	// Signal channel: background goroutine reads stdin and sends on stop. EOF / closed stdin
	// isn't an abort signal (otherwise iterate would exit immediately when stdin isn't a TTY) —
	// only an actual read of at least one byte counts.
	stop := make(chan struct{}, 1)
	go func() {
		buf := make([]byte, 1)
		if n, err := os.Stdin.Read(buf); err == nil && n > 0 {
			stop <- struct{}{}
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

		// Streaming parallel shallow screen. deck.IterateParallel runs a persistent pool of
		// shallow workers, publishes each result on a per-mutation ready channel, and
		// short-circuits as soon as the main-goroutine deep confirmation lands an improvement —
		// so on a round where an early mutation wins we don't burn cycles screening the tail.
		// A ticker goroutine prints live "tested N/total" progress off the shared atomic so the
		// user sees the worker pool moving even during long rounds; the ticker exits cleanly once
		// IterateParallel returns.
		var tested atomic.Int64
		tickerDone := make(chan struct{})
		go func() {
			t := time.NewTicker(500 * time.Millisecond)
			defer t.Stop()
			for {
				select {
				case <-tickerDone:
					return
				case <-t.C:
					fmt.Fprintf(os.Stderr, "\r[round %d] tested %d/%d (%s elapsed)        ",
						round, tested.Load(), len(mutations), time.Since(start).Truncate(time.Second))
				}
			}
		}()
		d, avg, idx, found := deck.IterateParallel(
			mutations, bestAvg, cfg.shallowShuffles, cfg.deepShuffles, cfg.incoming, 0,
			rng.Int63(), rng, &tested,
		)
		close(tickerDone)

		improved := false
		select {
		case <-stop:
			fmt.Fprintf(os.Stderr, "\nAborted mid-round after %d rounds / %d improvements in %s\n",
				round, improvements, time.Since(start).Truncate(time.Second))
			fmt.Println()
			printBestDeck(best)
			return
		default:
		}
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
