package deck

// Aggregate statistics accumulated while simulating a deck: total / per-cycle / per-card tallies,
// the single best turn ever seen, and a histogram of hand values that supports Min / Max /
// Median without retaining every individual hand.

import (
	"sort"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/hand"
)

// Stats holds aggregate hand-value statistics across all simulated runs.
type Stats struct {
	Runs        int
	Hands       int
	TotalValue  float64
	FirstCycle  CycleStats
	SecondCycle CycleStats
	// Best is the single highest-value hand seen across all runs (ties broken by first
	// occurrence). Summary.BestLine is in canonical (post-sort) order. Zero-valued if no hands
	// have been evaluated.
	Best BestTurn
	// PerCard attributes hand-level outcomes back to the cards that appeared in those hands.
	// Populated once per hand after hand.Best picks the winner — attribution cost is negligible
	// next to the underlying search.
	PerCard map[card.ID]CardPlayStats
	// Histogram counts hands seen at each integer Value. Keyed by TurnSummary.Value so Min /
	// Median can be derived without retaining every hand's value. Nil until the first hand is
	// evaluated.
	Histogram map[int]int
}

// BestTurn records a single hand and its optimal turn — the peak draw a deck saw during
// simulation. Summary.BestLine carries the cards and roles in canonical order.
type BestTurn struct {
	Summary hand.TurnSummary
	// StartingRunechants is the Runechant count carried in from the previous turn when this hand
	// was played. Only meaningful for Runeblade heroes.
	StartingRunechants int
}

// CardPlayStats captures how a single card contributed across hands it appeared in. Plays counts
// hands where it attacked or defended; Pitches counts hands where it was spent for resources.
// TotalContribution sums role-specific credit from the winning-line replay:
//
//   - Pitch   → Card.Pitch() (1/2/3 resource value, damage-equivalent by convention).
//   - Attack  → Card.Play() return plus the hero's OnCardPlayed trigger chained off it, at the
//     moment the card resolved in the winning permutation.
//   - Defend  → proportional share of min(sumDefense, incomingDamage), plus the card's own
//     Play return if it's a defense reaction.
//
// Useful as a directional per-card signal. The Defense share is proportional not causal: a
// defender soaking the whole block looks equal to a weaker one padding the same partition.
type CardPlayStats struct {
	Plays             int
	Pitches           int
	TotalContribution float64
}

// Avg returns mean per-card contribution across every hand where this card appeared (Plays +
// Pitches). Returns 0 when the card was never seen.
func (c CardPlayStats) Avg() float64 {
	n := c.Plays + c.Pitches
	if n == 0 {
		return 0
	}
	return c.TotalContribution / float64(n)
}

// CycleStats tracks total value and hand count for a single deck cycle.
type CycleStats struct {
	Hands int
	Total float64
}

// Mean returns the arithmetic mean hand value for this cycle.
func (c CycleStats) Mean() float64 {
	if c.Hands == 0 {
		return 0
	}
	return c.Total / float64(c.Hands)
}

// Mean returns the overall arithmetic mean hand value.
func (s Stats) Mean() float64 {
	if s.Hands == 0 {
		return 0
	}
	return s.TotalValue / float64(s.Hands)
}

// Min returns the lowest Value any simulated hand produced. Zero when no hands have been seen.
func (s Stats) Min() int {
	if len(s.Histogram) == 0 {
		return 0
	}
	first := true
	m := 0
	for v := range s.Histogram {
		if first || v < m {
			m = v
			first = false
		}
	}
	return m
}

// Max returns the highest Value any simulated hand produced. Zero when no hands have been seen.
func (s Stats) Max() int {
	m := 0
	for v := range s.Histogram {
		if v > m {
			m = v
		}
	}
	return m
}

// Median returns the median hand value. With an even number of hands it's the mean of the two
// middle values (so it can be fractional). Zero when no hands have been seen.
func (s Stats) Median() float64 {
	if s.Hands == 0 || len(s.Histogram) == 0 {
		return 0
	}
	keys := make([]int, 0, len(s.Histogram))
	for v := range s.Histogram {
		keys = append(keys, v)
	}
	sort.Ints(keys)
	// Walk the sorted values in order, counting cumulative hands until we pass the median
	// rank(s). lower = rank s.Hands/2 (0-indexed); upper = rank (s.Hands-1)/2 for even Hands.
	lowerRank := (s.Hands - 1) / 2
	upperRank := s.Hands / 2
	var lower, upper int
	cum := 0
	foundLower := false
	for _, v := range keys {
		cum += s.Histogram[v]
		if !foundLower && cum > lowerRank {
			lower = v
			foundLower = true
		}
		if cum > upperRank {
			upper = v
			break
		}
	}
	return float64(lower+upper) / 2
}
