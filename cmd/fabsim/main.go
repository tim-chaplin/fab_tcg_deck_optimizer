// Command fabsim simulates a Flesh and Blood deck and reports average
// hand value per cycle.
package main

import (
	"flag"
	"fmt"
	"math/rand"
	"time"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

func main() {
	runs := flag.Int("runs", 10000, "number of simulation runs (shuffles)")
	incoming := flag.Int("incoming", 4, "opponent damage per turn")
	seed := flag.Int64("seed", time.Now().UnixNano(), "RNG seed")
	flag.Parse()

	deck := buildDeck()
	rng := rand.New(rand.NewSource(*seed))
	stats := sim.Run(deck, *runs, *incoming, rng)

	fmt.Printf("Deck:           %d cards (hardcoded 20 blue / 20 red)\n", len(deck))
	fmt.Printf("Runs:           %d\n", stats.Runs)
	fmt.Printf("Hands:          %d\n", stats.Hands)
	fmt.Printf("Incoming/turn:  %d\n", *incoming)
	fmt.Printf("Seed:           %d\n", *seed)
	fmt.Println()
	fmt.Printf("Avg hand value (overall):       %.3f\n", stats.Avg())
	fmt.Printf("Avg hand value (cycle 1):       %.3f  (%d hands)\n", stats.FirstCycle.Avg(), stats.FirstCycle.Hands)
	fmt.Printf("Avg hand value (cycle 2):       %.3f  (%d hands)\n", stats.SecondCycle.Avg(), stats.SecondCycle.Hands)
}

func buildDeck() []card.Card {
	deck := make([]card.Card, 0, 40)
	for i := 0; i < 20; i++ {
		deck = append(deck, card.TestCardBlue{})
	}
	for i := 0; i < 20; i++ {
		deck = append(deck, card.TestCardRed{})
	}
	return deck
}
