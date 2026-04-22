package main

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"sync/atomic"
	"time"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/deck"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/hand"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/hero"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/mydecks"
)

// runIterateRestarts runs the hill climb N times from independent starts, resolving template's
// single '*' to 1..N for each deck name. Each run loads the deck at mydecks/<name>.json if one
// exists (so interrupted sessions resume without losing work — an already-converged deck
// passes through one no-op round), or generates a fresh random starting deck when missing.
// Prints a final ranking across all restarts.
func runIterateRestarts(cfg config, template string, n int) {
	type result struct {
		deckName string
		avg      float64
	}
	results := make([]result, 0, n)
	for i := 1; i <= n; i++ {
		deckName := strings.Replace(template, "*", strconv.Itoa(i), 1)
		outPath, err := mydecks.Path(deckName)
		if err != nil {
			die("%v", err)
		}
		fmt.Fprintf(os.Stderr, "\n=== Restart %d/%d: %s ===\n", i, n, deckName)
		cfg.outPath = outPath
		avg := runIterate(cfg)
		results = append(results, result{deckName: deckName, avg: avg})
	}

	fmt.Fprintln(os.Stderr, "\n=== Restart ranking ===")
	for _, r := range results {
		fmt.Fprintf(os.Stderr, "  %-40s avg %.3f\n", r.deckName, r.avg)
	}
	var best result
	for i, r := range results {
		if i == 0 || r.avg > best.avg {
			best = r
		}
	}
	fmt.Fprintf(os.Stderr, "\nBest: %s (avg %.3f)\n", best.deckName, best.avg)
}

func runIterate(cfg config) float64 {
	rng := rand.New(rand.NewSource(cfg.seed))

	best, bestAvg := prepareBaseline(cfg, rng)
	fmt.Println("Press Enter to abort.")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	watchStdinForAbort(cancel)

	// Deterministic hill-climb: enumerate every single-slot mutation of the current best, adopt
	// the first mutation that scores higher, and restart enumeration. Exhausting every mutation
	// with no improvement means we're at a local maximum.
	round := 0
	improvements := 0
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
			return bestAvg
		}
		if !found {
			fmt.Fprintf(os.Stderr, "\nLocal maximum reached after %d rounds / %d improvements in %s\n",
				round, improvements, time.Since(start).Truncate(time.Second))
			fmt.Println()
			printBestDeck(best)
			return bestAvg
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
