package deck

import (
	"io"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

// Per-deck evaluation result types. Pure data + small derived-stat methods, with no
// reference to sim or gameengine so this file can sit in internal/deck (which gameengine
// imports) without introducing a cycle. BestTurn carries only the durable outcome —
// runtime attack-turn scratch lives on sim's TurnSummary.

// Stats is the aggregate hand-value statistics across all simulated runs of a Deck.
type Stats struct {
	Runs        int
	Hands       int
	TotalValue  float64
	FirstCycle  CycleStats
	SecondCycle CycleStats
	// Best is the single highest-value hand seen across all runs (ties broken by first
	// occurrence). Zero-valued if no hands have been evaluated.
	Best BestTurn
	// Histogram counts hands seen at each integer Value. Keyed by best-turn Value so
	// Min / Max can be derived without retaining every hand's value. Nil until the first
	// hand is evaluated.
	Histogram map[int]int
	// IncomingPhysicalDamage and IncomingArcaneDamage record the matchup that scored every
	// hand, so a saved deck carries the assumptions it was tuned against. Zero when unevaluated
	// (Runs == 0).
	IncomingPhysicalDamage int
	IncomingArcaneDamage   int
	// PrintBest streams the peak turn's printout to w. Nil when no callable replay is
	// attached — loaded-from-JSON stats and evals with no recorded best — so callers
	// nil-check and skip the section.
	PrintBest func(w io.Writer) `json:"-"`
}

// Mean returns the overall arithmetic mean hand value.
func (s Stats) Mean() float64 {
	if s.Hands == 0 {
		return 0
	}
	return s.TotalValue / float64(s.Hands)
}

// Min returns the lowest Value any simulated hand produced. Zero when no hands have been
// seen.
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

// Max returns the highest Value any simulated hand produced. Zero when no hands have
// been seen.
func (s Stats) Max() int {
	m := 0
	for v := range s.Histogram {
		if v > m {
			m = v
		}
	}
	return m
}

// BestTurn records the peak draw a deck saw during simulation. The attack turn's printout is
// produced lazily via Stats.PrintBest at render time, not eval time.
type BestTurn struct {
	Value    int
	BestLine []card.CardAssignment
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
