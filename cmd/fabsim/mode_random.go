package main

import (
	"fmt"
	"math/rand"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/deck"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/hero"
)

func runRandom(cfg config) {
	rng := rand.New(rand.NewSource(cfg.seed))
	candidates := sampleShallowCandidates(cfg, rng)

	avgs := make([]float64, len(candidates))
	for i, c := range candidates {
		avgs[i] = c.avg
	}
	min, median, max := summarize(avgs)
	fmt.Printf("Phase 1: %d decks, %d shuffles each, incoming=%d, seed=%d\n",
		cfg.numDecks, cfg.shallowShuffles, cfg.incoming, cfg.seed)
	fmt.Printf("Deck value distribution: min %.3f  median %.3f  max %.3f\n", min, median, max)
	fmt.Println()

	// Select top N by shallow average for phase 2.
	sort.Slice(candidates, func(i, j int) bool {
		return candidates[i].avg > candidates[j].avg
	})
	topN := cfg.topN
	if topN > len(candidates) {
		topN = len(candidates)
	}
	bestDeck := evaluateFinalistsDeep(cfg, rng, candidates[:topN])

	fmt.Printf("Phase 2: re-evaluated top %d decks with %d shuffles\n", topN, cfg.deepShuffles)
	fmt.Println()
	printBestDeck(bestDeck)

	if cfg.outPath != "" {
		saveIfBetter(bestDeck, cfg.outPath)
	}
}

// shallowCandidate pairs a generated deck with its phase-1 (shallow) average so phase 2 can
// re-evaluate the top finalists without re-sampling.
type shallowCandidate struct {
	deck *deck.Deck
	avg  float64
}

// sampleShallowCandidates is phase 1: generate cfg.numDecks random decks and score each at
// cfg.shallowShuffles. The progress bar renders to stderr so piping stdout yields clean output.
func sampleShallowCandidates(cfg config, rng *rand.Rand) []shallowCandidate {
	candidates := make([]shallowCandidate, 0, cfg.numDecks)
	fmt.Fprintf(os.Stderr, "Phase 1: evaluating %d decks (%d shuffles each)\n", cfg.numDecks, cfg.shallowShuffles)
	start := time.Now()
	for i := 0; i < cfg.numDecks; i++ {
		d := deck.Random(hero.Viserai{}, cfg.deckSize, cfg.maxCopies, rng, cfg.legalFilter())
		stats := d.Evaluate(cfg.shallowShuffles, cfg.incoming, rng)
		candidates = append(candidates, shallowCandidate{deck: d, avg: stats.Avg()})
		printProgress(i+1, cfg.numDecks, time.Since(start))
	}
	fmt.Fprintln(os.Stderr)
	return candidates
}

// evaluateFinalistsDeep is phase 2: re-evaluate each finalist at cfg.deepShuffles against a
// fresh Deck (zeroed Stats) so phase-1 noise doesn't leak in, and return whichever scores highest.
func evaluateFinalistsDeep(cfg config, rng *rand.Rand, finalists []shallowCandidate) *deck.Deck {
	fmt.Fprintf(os.Stderr, "Phase 2: re-evaluating top %d decks (%d shuffles each)\n", len(finalists), cfg.deepShuffles)
	var bestDeck *deck.Deck
	bestAvg := -1.0
	start := time.Now()
	for i, c := range finalists {
		d := deck.New(c.deck.Hero, c.deck.Weapons, c.deck.Cards)
		stats := d.Evaluate(cfg.deepShuffles, cfg.incoming, rng)
		avg := stats.Avg()
		if avg > bestAvg {
			bestAvg = avg
			bestDeck = d
		}
		printProgress(i+1, len(finalists), time.Since(start))
	}
	fmt.Fprintln(os.Stderr)
	return bestDeck
}

// summarize returns (min, median, max) of vs. Panics if vs is empty. For even-length vs, median
// is the mean of the two middle elements.
func summarize(vs []float64) (min, median, max float64) {
	sorted := make([]float64, len(vs))
	copy(sorted, vs)
	sort.Float64s(sorted)
	n := len(sorted)
	min = sorted[0]
	max = sorted[n-1]
	if n%2 == 1 {
		median = sorted[n/2]
	} else {
		median = (sorted[n/2-1] + sorted[n/2]) / 2
	}
	return
}

// printProgress renders a single-line \r-overwriting progress bar to stderr. Shows deck count,
// percent, elapsed time, and ETA from average per-deck time so far.
func printProgress(done, total int, elapsed time.Duration) {
	const width = 30
	frac := float64(done) / float64(total)
	filled := int(frac * float64(width))
	bar := strings.Repeat("=", filled) + strings.Repeat(" ", width-filled)
	var eta time.Duration
	if done > 0 {
		eta = time.Duration(float64(elapsed) * float64(total-done) / float64(done))
	}
	fmt.Fprintf(os.Stderr, "\r[%s] %d/%d (%.0f%%) elapsed %s eta %s",
		bar, done, total, frac*100, elapsed.Truncate(time.Second), eta.Truncate(time.Second))
}

// saveIfBetter writes d to outPath if its average exceeds the deck already on disk (or none
// exists). Prints a status line either way.
func saveIfBetter(d *deck.Deck, outPath string) {
	prev, prevAvg := loadExisting(outPath)
	if prev != nil && d.Stats.Avg() <= prevAvg {
		fmt.Printf("\nPrevious best (%.3f) >= current (%.3f), %s unchanged\n", prevAvg, d.Stats.Avg(), outPath)
		return
	}
	if err := writeDeck(d, outPath); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	if prev != nil {
		fmt.Printf("\nNew best (%.3f) beats previous (%.3f), wrote %s\n", d.Stats.Avg(), prevAvg, outPath)
	} else {
		fmt.Printf("\nWrote best deck to %s\n", outPath)
	}
}
