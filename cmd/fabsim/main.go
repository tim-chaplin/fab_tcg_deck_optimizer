// Command fabsim generates N random Viserai decks, evaluates each, and reports the best one.
package main

import (
	"flag"
	"fmt"
	"math/rand"
	"sort"
	"time"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/deck"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/hand"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/hero"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/weapon"
)

func main() {
	numDecks := flag.Int("decks", 100, "number of random decks to generate and evaluate")
	shuffles := flag.Int("shuffles", 1000, "number of shuffles to simulate per deck")
	incoming := flag.Int("incoming", 4, "opponent damage per turn")
	deckSize := flag.Int("deck-size", 40, "number of cards per deck")
	maxCopies := flag.Int("max-copies", 2, "maximum copies of any single card printing per deck")
	seed := flag.Int64("seed", time.Now().UnixNano(), "RNG seed")
	flag.Parse()

	rng := rand.New(rand.NewSource(*seed))

	var bestDeck *deck.Deck
	bestAvg := -1.0

	for i := 0; i < *numDecks; i++ {
		d := deck.Random(hero.Viserai{}, *deckSize, *maxCopies, rng)
		stats := d.Evaluate(*shuffles, *incoming, rng)
		if stats.Avg() > bestAvg {
			bestAvg = stats.Avg()
			bestDeck = d
		}
	}

	fmt.Printf("Generated %d decks, %d shuffles each, incoming=%d, seed=%d\n",
		*numDecks, *shuffles, *incoming, *seed)
	fmt.Println()
	printBestDeck(bestDeck)
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
		fmt.Printf("  Best hand seen (value %d): %s\n", b.Play.Value, hand.FormatRoles(b.Hand, b.Play.Roles))
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

func weaponNames(ws []weapon.Weapon) string {
	names := make([]string, len(ws))
	for i, w := range ws {
		names[i] = w.Name()
	}
	return fmt.Sprintf("%v", names)
}
