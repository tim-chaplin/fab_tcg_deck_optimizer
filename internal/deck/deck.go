// Package deck represents a candidate FaB deck and the hand-value stats
// accumulated from simulating it. Future deck-search code will create
// many Decks, evaluate each, and compare their Stats.
package deck

import (
	"math/rand"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/hand"
)

// Deck is a single candidate deck and the hand-value stats accumulated
// from simulating it.
type Deck struct {
	Cards []card.Card
	Stats Stats
}

// New constructs a Deck with the given cards and zeroed stats.
func New(cards []card.Card) *Deck {
	return &Deck{Cards: cards}
}

// Stats holds aggregate hand-value statistics across all simulated runs.
type Stats struct {
	Runs        int
	Hands       int
	TotalValue  float64
	FirstCycle  CycleStats
	SecondCycle CycleStats
}

// CycleStats tracks total value and hand count for a single deck cycle.
type CycleStats struct {
	Hands int
	Total float64
}

// Avg returns the average hand value for this cycle.
func (c CycleStats) Avg() float64 {
	if c.Hands == 0 {
		return 0
	}
	return c.Total / float64(c.Hands)
}

// Avg returns the overall average hand value.
func (s Stats) Avg() float64 {
	if s.Hands == 0 {
		return 0
	}
	return s.TotalValue / float64(s.Hands)
}

// Evaluate simulates `runs` shuffles of the deck. For each run it draws
// successive hands of hand.HandSize from the top, computes the optimal
// play against an opponent attacking for incomingDamage, and returns
// Pitched cards to the bottom of the deck (in hand order). Played and
// defended cards are spent. Each run ends when fewer than hand.HandSize
// cards remain.
//
// A "cycle" is one pass through the original deck size: cumulative hands
// 0..(deckSize/HandSize - 1) are cycle 1, the next deckSize/HandSize
// hands are cycle 2.
//
// Results accumulate into d.Stats and are also returned for convenience.
func (d *Deck) Evaluate(runs int, incomingDamage int, rng *rand.Rand) Stats {
	d.Stats.Runs += runs
	deckSize := len(d.Cards)
	if deckSize < hand.HandSize {
		return d.Stats
	}
	handsPerCycle := deckSize / hand.HandSize

	working := make([]card.Card, 0, deckSize)
	for r := 0; r < runs; r++ {
		working = append(working[:0], d.Cards...)
		rng.Shuffle(len(working), func(i, j int) {
			working[i], working[j] = working[j], working[i]
		})

		handIdx := 0
		for len(working) >= hand.HandSize {
			h := working[:hand.HandSize]
			play := hand.Best(h, incomingDamage)
			v := float64(play.Value())

			d.Stats.TotalValue += v
			d.Stats.Hands++
			switch handIdx / handsPerCycle {
			case 0:
				d.Stats.FirstCycle.Hands++
				d.Stats.FirstCycle.Total += v
			case 1:
				d.Stats.SecondCycle.Hands++
				d.Stats.SecondCycle.Total += v
			}

			// Recycle: pitched cards go to the bottom (in hand order);
			// attacked and defended cards are spent. If nothing was
			// pitched, the deck shrinks by hand.HandSize this turn.
			pitched := make([]card.Card, 0, hand.HandSize)
			for i, c := range h {
				if play.Roles[i] == hand.Pitch {
					pitched = append(pitched, c)
				}
			}
			working = append(working[hand.HandSize:], pitched...)
			handIdx++
		}
	}
	return d.Stats
}
