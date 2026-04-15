// Package sim runs the deck simulation: shuffle, draw hands of HandSize,
// pitch back to the bottom of the deck per FaB rules, and aggregate
// hand-value stats over many runs.
package sim

import (
	"math/rand"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/hand"
)

// HandSize is the number of cards drawn per hand.
const HandSize = 4

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

// Run simulates `runs` shuffles of `deck`. For each run it draws successive
// hands of HandSize from the top of the deck, computes the optimal play
// against an opponent attacking for incomingDamage, and returns Pitched
// cards to the bottom of the deck (in hand order). Played and defended
// cards are spent. The run ends when fewer than HandSize cards remain.
//
// A "cycle" is one pass through the original deck size: cumulative hands
// 0..(deckSize/HandSize - 1) are cycle 1, the next deckSize/HandSize
// hands are cycle 2.
func Run(deck []card.Card, runs int, incomingDamage int, rng *rand.Rand) Stats {
	stats := Stats{Runs: runs}
	deckSize := len(deck)
	handsPerCycle := deckSize / HandSize

	working := make([]card.Card, 0, deckSize)
	for r := 0; r < runs; r++ {
		working = append(working[:0], deck...)
		rng.Shuffle(len(working), func(i, j int) {
			working[i], working[j] = working[j], working[i]
		})

		handIdx := 0
		for len(working) >= HandSize {
			h := working[:HandSize]
			play := hand.Best(h, incomingDamage)
			v := float64(play.Value())

			stats.TotalValue += v
			stats.Hands++
			cycle := handIdx / handsPerCycle
			switch cycle {
			case 0:
				stats.FirstCycle.Hands++
				stats.FirstCycle.Total += v
			case 1:
				stats.SecondCycle.Hands++
				stats.SecondCycle.Total += v
			}

			// Recycle: pitched cards go to the bottom (in hand order);
			// attacked and defended cards are spent. If nothing was
			// pitched, the deck shrinks by HandSize this turn.
			pitched := make([]card.Card, 0, HandSize)
			for i, c := range h {
				if play.Roles[i] == hand.Pitch {
					pitched = append(pitched, c)
				}
			}
			working = append(working[HandSize:], pitched...)
			handIdx++
		}
	}
	return stats
}
