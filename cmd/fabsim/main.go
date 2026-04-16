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
	numDecks := flag.Int("decks", 1000, "number of random decks to generate and evaluate")
	shuffles := flag.Int("shuffles", 100, "number of shuffles to simulate per deck")
	incoming := flag.Int("incoming", 4, "opponent damage per turn")
	deckSize := flag.Int("deck-size", 40, "number of cards per deck")
	maxCopies := flag.Int("max-copies", 2, "maximum copies of any single card printing per deck")
	seed := flag.Int64("seed", time.Now().UnixNano(), "RNG seed")
	outPath := flag.String("out", "best_deck.json", "path to write the best deck as JSON")
	flag.Parse()

	rng := rand.New(rand.NewSource(*seed))

	var bestDeck *deck.Deck
	bestAvg := -1.0
	avgs := make([]float64, 0, *numDecks)

	start := time.Now()
	for i := 0; i < *numDecks; i++ {
		d := deck.Random(hero.Viserai{}, *deckSize, *maxCopies, rng)
		stats := d.Evaluate(*shuffles, *incoming, rng)
		avg := stats.Avg()
		avgs = append(avgs, avg)
		if avg > bestAvg {
			bestAvg = avg
			bestDeck = d
		}
		printProgress(i+1, *numDecks, time.Since(start))
	}
	fmt.Fprintln(os.Stderr)

	fmt.Printf("Generated %d decks, %d shuffles each, incoming=%d, seed=%d\n",
		*numDecks, *shuffles, *incoming, *seed)
	min, median, max := summarize(avgs)
	fmt.Printf("Deck value distribution: min %.3f  median %.3f  max %.3f\n", min, median, max)
	fmt.Println()
	printBestDeck(bestDeck)

	if *outPath != "" {
		data, err := deckio.Marshal(bestDeck)
		if err != nil {
			fmt.Fprintf(os.Stderr, "marshal best deck: %v\n", err)
			os.Exit(1)
		}
		if err := os.WriteFile(*outPath, data, 0o644); err != nil {
			fmt.Fprintf(os.Stderr, "write %s: %v\n", *outPath, err)
			os.Exit(1)
		}
		fmt.Printf("\nWrote best deck to %s\n", *outPath)
	}
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
