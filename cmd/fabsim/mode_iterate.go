package main

import (
	"fmt"
	"math/rand"
	"os"
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

		improved := false
		for i, mut := range mutations {
			select {
			case <-stop:
				fmt.Fprintf(os.Stderr, "\nAborted mid-round after %d rounds / %d improvements in %s\n",
					round, improvements, time.Since(start).Truncate(time.Second))
				fmt.Println()
				printBestDeck(best)
				return
			default:
			}

			// Two-stage evaluation: most mutations are neutral or worse, so a cheap shallow screen
			// filters them out. Only mutations that clear bestAvg on the shallow sample graduate to
			// a fresh -deep-shuffles evaluation that decides adoption. bestAvg always reflects
			// deep-shuffles depth (that's what we wrote to disk), so the apples-to-apples check
			// happens against the deep re-eval, not the shallow screen.
			screen := deck.New(mut.Deck.Hero, mut.Deck.Weapons, mut.Deck.Cards)
			shallowAvg := screen.Evaluate(cfg.shallowShuffles, cfg.incoming, rng).Avg()
			if shallowAvg > bestAvg {
				d := deck.New(mut.Deck.Hero, mut.Deck.Weapons, mut.Deck.Cards)
				avg := d.Evaluate(cfg.deepShuffles, cfg.incoming, rng).Avg()
				if avg > bestAvg {
					improvements++
					fmt.Fprintf(os.Stderr, "\r[round %d] improvement at %d/%d: shallow %.3f → deep %.3f beats %.3f (%s), restarting        \n",
						round, i+1, len(mutations), shallowAvg, avg, bestAvg, mut.Description)
					bestAvg = avg
					best = d
					_ = writeDeck(best, cfg.outPath)
					improved = true
					break
				}
				fmt.Fprintf(os.Stderr, "\r[round %d] shallow %.3f at %d/%d (%s) not confirmed by deep %.3f        \n",
					round, shallowAvg, i+1, len(mutations), mut.Description, avg)
			}
			if (i+1)%50 == 0 {
				fmt.Fprintf(os.Stderr, "\r[round %d] %d/%d evaluated, best still %.3f (%s elapsed)        ",
					round, i+1, len(mutations), bestAvg, time.Since(start).Truncate(time.Second))
			}
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
