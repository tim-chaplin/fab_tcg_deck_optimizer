// Command fabsim generates N random Viserai decks, evaluates each, and reports the best one.
package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/deck"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/deckio"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/hand"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/hero"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/weapon"
)

func main() {
	mode := flag.String("mode", "random", "run mode: random, iterate, or print_only")
	numDecks := flag.Int("decks", 10000, "number of random decks to generate (phase 1)")
	shallowShuffles := flag.Int("shallow-shuffles", 10, "shuffles per deck in phase 1 (wide search)")
	topN := flag.Int("top-n", 100, "number of top decks to advance to phase 2")
	deepShuffles := flag.Int("deep-shuffles", 1000, "shuffles per deck in phase 2 (deep evaluation)")
	incoming := flag.Int("incoming", 4, "opponent damage per turn")
	deckSize := flag.Int("deck-size", 40, "number of cards per deck")
	maxCopies := flag.Int("max-copies", 2, "maximum copies of any single card printing per deck")
	seed := flag.Int64("seed", time.Now().UnixNano(), "RNG seed")
	outPath := flag.String("out", "best_deck.json", "path to write/read the best deck JSON")
	flag.Parse()

	switch *mode {
	case "random":
		runRandom(*numDecks, *shallowShuffles, *topN, *deepShuffles, *incoming, *deckSize, *maxCopies, *seed, *outPath)
	case "iterate":
		runIterate(*deepShuffles, *incoming, *seed, *outPath)
	case "print_only":
		runPrintOnly(*outPath)
	default:
		fmt.Fprintf(os.Stderr, "unknown mode %q (want random, iterate, or print_only)\n", *mode)
		os.Exit(1)
	}
}

func runRandom(numDecks, shallowShuffles, topN, deepShuffles, incoming, deckSize, maxCopies int, seed int64, outPath string) {
	rng := rand.New(rand.NewSource(seed))

	// Phase 1: wide shallow search.
	type candidate struct {
		deck *deck.Deck
		avg  float64
	}
	candidates := make([]candidate, 0, numDecks)

	fmt.Fprintf(os.Stderr, "Phase 1: evaluating %d decks (%d shuffles each)\n", numDecks, shallowShuffles)
	start := time.Now()
	for i := 0; i < numDecks; i++ {
		d := deck.Random(hero.Viserai{}, deckSize, maxCopies, rng)
		stats := d.Evaluate(shallowShuffles, incoming, rng)
		candidates = append(candidates, candidate{deck: d, avg: stats.Avg()})
		printProgress(i+1, numDecks, time.Since(start))
	}
	fmt.Fprintln(os.Stderr)

	avgs := make([]float64, len(candidates))
	for i, c := range candidates {
		avgs[i] = c.avg
	}
	min, median, max := summarize(avgs)
	fmt.Printf("Phase 1: %d decks, %d shuffles each, incoming=%d, seed=%d\n",
		numDecks, shallowShuffles, incoming, seed)
	fmt.Printf("Deck value distribution: min %.3f  median %.3f  max %.3f\n", min, median, max)
	fmt.Println()

	// Select top N by shallow average.
	sort.Slice(candidates, func(i, j int) bool {
		return candidates[i].avg > candidates[j].avg
	})
	if topN > len(candidates) {
		topN = len(candidates)
	}
	finalists := candidates[:topN]

	// Phase 2: deep evaluation of finalists with fresh stats.
	fmt.Fprintf(os.Stderr, "Phase 2: re-evaluating top %d decks (%d shuffles each)\n", topN, deepShuffles)
	var bestDeck *deck.Deck
	bestAvg := -1.0
	start = time.Now()
	for i, c := range finalists {
		d := deck.New(c.deck.Hero, c.deck.Weapons, c.deck.Cards)
		stats := d.Evaluate(deepShuffles, incoming, rng)
		avg := stats.Avg()
		if avg > bestAvg {
			bestAvg = avg
			bestDeck = d
		}
		printProgress(i+1, topN, time.Since(start))
	}
	fmt.Fprintln(os.Stderr)

	fmt.Printf("Phase 2: re-evaluated top %d decks with %d shuffles\n", topN, deepShuffles)
	fmt.Println()
	printBestDeck(bestDeck)

	if outPath != "" {
		prev, prevAvg := loadExisting(outPath)
		if prev == nil || bestDeck.Stats.Avg() > prevAvg {
			data, err := deckio.Marshal(bestDeck)
			if err != nil {
				fmt.Fprintf(os.Stderr, "marshal best deck: %v\n", err)
				os.Exit(1)
			}
			if err := os.WriteFile(outPath, data, 0o644); err != nil {
				fmt.Fprintf(os.Stderr, "write %s: %v\n", outPath, err)
				os.Exit(1)
			}
			if prev != nil {
				fmt.Printf("\nNew best (%.3f) beats previous (%.3f), wrote %s\n", bestDeck.Stats.Avg(), prevAvg, outPath)
			} else {
				fmt.Printf("\nWrote best deck to %s\n", outPath)
			}
		} else {
			fmt.Printf("\nPrevious best (%.3f) >= current (%.3f), %s unchanged\n", prevAvg, bestDeck.Stats.Avg(), outPath)
		}
	}
}

func runIterate(shuffles, incoming int, seed int64, outPath string) {
	best, bestAvg := loadExisting(outPath)
	if best == nil {
		fmt.Fprintf(os.Stderr, "no existing deck at %s — run with --mode=random first\n", outPath)
		os.Exit(1)
	}
	fmt.Printf("Loaded best deck (avg %.3f) from %s\n", bestAvg, outPath)
	fmt.Println("Press Enter to stop.")

	rng := rand.New(rand.NewSource(seed))

	// Signal channel: background goroutine reads stdin and sends on stop.
	stop := make(chan struct{}, 1)
	go func() {
		buf := make([]byte, 1)
		os.Stdin.Read(buf)
		stop <- struct{}{}
	}()

	iter := 0
	improvements := 0
	start := time.Now()
	for {
		select {
		case <-stop:
			fmt.Fprintf(os.Stderr, "\nStopping after %d iterations (%d improvements) in %s\n",
				iter, improvements, time.Since(start).Truncate(time.Second))
			fmt.Println()
			printBestDeck(best)
			return
		default:
		}

		iter++
		candidate := deck.Mutate(best, rng)
		d := deck.New(candidate.Hero, candidate.Weapons, candidate.Cards)
		stats := d.Evaluate(shuffles, incoming, rng)
		avg := stats.Avg()

		if avg > bestAvg {
			improvements++
			bestAvg = avg
			best = d
			data, err := deckio.Marshal(best)
			if err == nil {
				os.WriteFile(outPath, data, 0o644)
			}
			fmt.Fprintf(os.Stderr, "\r[iter %d] new best: %.3f (+%d improvements)        \n",
				iter, bestAvg, improvements)
		}
		if iter%100 == 0 {
			fmt.Fprintf(os.Stderr, "\r[iter %d] best: %.3f (%d improvements, %s elapsed)        ",
				iter, bestAvg, improvements, time.Since(start).Truncate(time.Second))
		}
	}
}

// loadExisting reads and deserializes the deck at path. Returns (nil, 0) if the file doesn't
// exist or can't be parsed — the caller treats that as "no previous best".
func loadExisting(path string) (*deck.Deck, float64) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, 0
	}
	d, err := deckio.Unmarshal(data)
	if err != nil {
		return nil, 0
	}
	return d, d.Stats.Avg()
}

func runPrintOnly(path string) {
	data, err := os.ReadFile(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "read %s: %v\n", path, err)
		os.Exit(1)
	}
	d, err := deckio.Unmarshal(data)
	if err != nil {
		fmt.Fprintf(os.Stderr, "unmarshal %s: %v\n", path, err)
		os.Exit(1)
	}
	printBestDeck(d)
}

// summarize returns (min, median, max) of vs. Panics if vs is empty. Median of an even-length
// slice is the mean of the two middle elements.
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

func printBestDeck(d *deck.Deck) {
	s := d.Stats
	fmt.Printf("Best deck (avg %.3f over %d hands)\n", s.Avg(), s.Hands)
	fmt.Printf("  Hero:    %s\n", d.Hero.Name())
	fmt.Printf("  Weapons: %s\n", weaponNames(d.Weapons))
	fmt.Printf("  Cycle 1 avg: %.3f  (%d hands)\n", s.FirstCycle.Avg(), s.FirstCycle.Hands)
	fmt.Printf("  Cycle 2 avg: %.3f  (%d hands)\n", s.SecondCycle.Avg(), s.SecondCycle.Hands)
	var red, yellow, blue int
	for _, c := range d.Cards {
		switch c.Pitch() {
		case 1:
			red++
		case 2:
			yellow++
		case 3:
			blue++
		}
	}
	fmt.Printf("  Pitch:   %d red / %d yellow / %d blue\n", red, yellow, blue)
	if b := s.Best; b.Hand != nil {
		line := hand.FormatRoles(b.Hand, b.Play.Roles)
		for _, w := range b.Play.Weapons {
			line += ", " + w + ": ATTACK"
		}
		fmt.Printf("  Best hand seen (value %d): %s\n", b.Play.Value, line)
	}
	fmt.Println()
	fmt.Println("Card list:")
	counts := map[string]int{}
	for _, c := range d.Cards {
		counts[c.Name()]++
	}
	names := make([]string, 0, len(counts))
	for n := range counts {
		names = append(names, n)
	}
	sort.Strings(names)
	for _, n := range names {
		fmt.Printf("  %dx %s\n", counts[n], n)
	}
}

// printProgress renders a single-line progress bar to stderr, overwriting itself with \r.
// Shows deck count, percent, elapsed time, and ETA based on the average per-deck time so far.
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

func weaponNames(ws []weapon.Weapon) string {
	names := make([]string, len(ws))
	for i, w := range ws {
		names[i] = w.Name()
	}
	return fmt.Sprintf("%v", names)
}
