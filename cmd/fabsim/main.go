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

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
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

	cfg := config{
		numDecks:        *numDecks,
		shallowShuffles: *shallowShuffles,
		topN:            *topN,
		deepShuffles:    *deepShuffles,
		incoming:        *incoming,
		deckSize:        *deckSize,
		maxCopies:       *maxCopies,
		seed:            *seed,
		outPath:         *outPath,
	}

	switch *mode {
	case "random":
		runRandom(cfg)
	case "iterate":
		runIterate(cfg)
	case "print_only":
		runPrintOnly(cfg.outPath)
	default:
		fmt.Fprintf(os.Stderr, "unknown mode %q (want random, iterate, or print_only)\n", *mode)
		os.Exit(1)
	}
}

type config struct {
	numDecks        int
	shallowShuffles int
	topN            int
	deepShuffles    int
	incoming        int
	deckSize        int
	maxCopies       int
	seed            int64
	outPath         string
}

func runRandom(cfg config) {
	rng := rand.New(rand.NewSource(cfg.seed))

	// Phase 1: wide shallow search.
	type candidate struct {
		deck *deck.Deck
		avg  float64
	}
	candidates := make([]candidate, 0, cfg.numDecks)

	fmt.Fprintf(os.Stderr, "Phase 1: evaluating %d decks (%d shuffles each)\n", cfg.numDecks, cfg.shallowShuffles)
	start := time.Now()
	for i := 0; i < cfg.numDecks; i++ {
		d := deck.Random(hero.Viserai{}, cfg.deckSize, cfg.maxCopies, rng)
		stats := d.Evaluate(cfg.shallowShuffles, cfg.incoming, rng)
		candidates = append(candidates, candidate{deck: d, avg: stats.Avg()})
		printProgress(i+1, cfg.numDecks, time.Since(start))
	}
	fmt.Fprintln(os.Stderr)

	avgs := make([]float64, len(candidates))
	for i, c := range candidates {
		avgs[i] = c.avg
	}
	min, median, max := summarize(avgs)
	fmt.Printf("Phase 1: %d decks, %d shuffles each, incoming=%d, seed=%d\n",
		cfg.numDecks, cfg.shallowShuffles, cfg.incoming, cfg.seed)
	fmt.Printf("Deck value distribution: min %.3f  median %.3f  max %.3f\n", min, median, max)
	fmt.Println()

	// Select top N by shallow average.
	sort.Slice(candidates, func(i, j int) bool {
		return candidates[i].avg > candidates[j].avg
	})
	topN := cfg.topN
	if topN > len(candidates) {
		topN = len(candidates)
	}
	finalists := candidates[:topN]

	// Phase 2: deep evaluation of finalists with fresh stats.
	fmt.Fprintf(os.Stderr, "Phase 2: re-evaluating top %d decks (%d shuffles each)\n", topN, cfg.deepShuffles)
	var bestDeck *deck.Deck
	bestAvg := -1.0
	start = time.Now()
	for i, c := range finalists {
		d := deck.New(c.deck.Hero, c.deck.Weapons, c.deck.Cards)
		stats := d.Evaluate(cfg.deepShuffles, cfg.incoming, rng)
		avg := stats.Avg()
		if avg > bestAvg {
			bestAvg = avg
			bestDeck = d
		}
		printProgress(i+1, topN, time.Since(start))
	}
	fmt.Fprintln(os.Stderr)

	fmt.Printf("Phase 2: re-evaluated top %d decks with %d shuffles\n", topN, cfg.deepShuffles)
	fmt.Println()
	printBestDeck(bestDeck)

	if cfg.outPath != "" {
		saveIfBetter(bestDeck, cfg.outPath)
	}
}

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
		if data, err := deckio.Marshal(best); err == nil {
			os.WriteFile(cfg.outPath, data, 0o644)
		}
		fmt.Printf("Starting deck avg %.3f, saved to %s\n", bestAvg, cfg.outPath)
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
		mutations := deck.AllMutations(best)
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

			d := deck.New(mut.Deck.Hero, mut.Deck.Weapons, mut.Deck.Cards)
			stats := d.Evaluate(cfg.deepShuffles, cfg.incoming, rng)
			avg := stats.Avg()
			if avg > bestAvg {
				improvements++
				fmt.Fprintf(os.Stderr, "\r[round %d] improvement at %d/%d: %.3f → %.3f (%s), restarting        \n",
					round, i+1, len(mutations), bestAvg, avg, mut.Description)
				bestAvg = avg
				best = d
				if data, err := deckio.Marshal(best); err == nil {
					os.WriteFile(cfg.outPath, data, 0o644)
				}
				improved = true
				break
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

// saveIfBetter writes d to outPath if its average exceeds the previously saved deck (or if no
// previous deck exists). Prints a status line either way.
func saveIfBetter(d *deck.Deck, outPath string) {
	prev, prevAvg := loadExisting(outPath)
	if prev != nil && d.Stats.Avg() <= prevAvg {
		fmt.Printf("\nPrevious best (%.3f) >= current (%.3f), %s unchanged\n", prevAvg, d.Stats.Avg(), outPath)
		return
	}
	data, err := deckio.Marshal(d)
	if err != nil {
		fmt.Fprintf(os.Stderr, "marshal best deck: %v\n", err)
		os.Exit(1)
	}
	if err := os.WriteFile(outPath, data, 0o644); err != nil {
		fmt.Fprintf(os.Stderr, "write %s: %v\n", outPath, err)
		os.Exit(1)
	}
	if prev != nil {
		fmt.Printf("\nNew best (%.3f) beats previous (%.3f), wrote %s\n", d.Stats.Avg(), prevAvg, outPath)
	} else {
		fmt.Printf("\nWrote best deck to %s\n", outPath)
	}
}

func runPrintOnly(path string) {
	d, _ := loadExisting(path)
	if d == nil {
		fmt.Fprintf(os.Stderr, "could not load deck from %s\n", path)
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
		prefix := fmt.Sprintf("  Best hand seen (value %d", b.Play.Value)
		if d.Hero.Types().Has(card.TypeRuneblade) {
			prefix += fmt.Sprintf(", %d carryover runechants", b.StartingRunechants)
		}
		fmt.Printf("%s): %s\n", prefix, line)
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
