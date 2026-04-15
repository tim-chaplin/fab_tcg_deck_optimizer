// Command fabsim simulates a Flesh and Blood deck and reports average
// hand value per cycle.
package main

import (
	"flag"
	"fmt"
	"math/rand"
	"time"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card/runeblade"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/deck"
)

func main() {
	shuffles := flag.Int("shuffles", 10000, "number of shuffles to simulate for this deck")
	incoming := flag.Int("incoming", 4, "opponent damage per turn")
	seed := flag.Int64("seed", time.Now().UnixNano(), "RNG seed")
	flag.Parse()

	d := deck.New(buildDeck())
	rng := rand.New(rand.NewSource(*seed))
	stats := d.Evaluate(*shuffles, *incoming, rng)

	fmt.Printf("Deck:           %d cards (Shrill of Skullform + Malefic Incantation, all colors)\n", len(d.Cards))
	fmt.Printf("Shuffles:       %d\n", stats.Runs)
	fmt.Printf("Hands:          %d\n", stats.Hands)
	fmt.Printf("Incoming/turn:  %d\n", *incoming)
	fmt.Printf("Seed:           %d\n", *seed)
	fmt.Println()
	fmt.Printf("Avg hand value (overall):       %.3f\n", stats.Avg())
	fmt.Printf("Avg hand value (cycle 1):       %.3f  (%d hands)\n", stats.FirstCycle.Avg(), stats.FirstCycle.Hands)
	fmt.Printf("Avg hand value (cycle 2):       %.3f  (%d hands)\n", stats.SecondCycle.Avg(), stats.SecondCycle.Hands)
}

// buildDeck assembles the demo deck from every card we've implemented so
// far. Each unique card variant gets the FaB per-name maximum of 3
// copies. This currently produces fewer than 40 cards — more variants
// will be added as they're implemented.
func buildDeck() []card.Card {
	variants := []card.Card{
		runeblade.ShrillOfSkullformRed{},
		runeblade.ShrillOfSkullformYellow{},
		runeblade.ShrillOfSkullformBlue{},
		runeblade.MaleficIncantationRed{},
		runeblade.MaleficIncantationYellow{},
		runeblade.MaleficIncantationBlue{},
	}
	deck := make([]card.Card, 0, len(variants)*3)
	for _, v := range variants {
		for i := 0; i < 3; i++ {
			deck = append(deck, v)
		}
	}
	return deck
}
