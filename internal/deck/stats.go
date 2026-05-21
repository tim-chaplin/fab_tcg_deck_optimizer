package deck

import (
	"io"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/ids"
)

// Per-deck evaluation result types. Pure data + small derived-stat methods, with no
// reference to sim or gameengine so this file can sit in internal/deck (which gameengine
// imports) without introducing a cycle. BestTurn carries only the durable outcome —
// runtime chain scratch lives on sim's TurnSummary.

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
	// PerCardMarginal is a coarse correlational view of each card's hand-value impact:
	// per unique card ID, the mean turn Value across turns where it was in the dealt hand
	// or arsenal-in slot vs turns where it wasn't. A smell test for buggy or oversimplified
	// implementations.
	//
	// Caveat: cards whose effect pays off on a LATER turn (auras, drawn-card payoffs) often
	// won't surface here — the source has rotated out of the hand by the time its value
	// lands, so the correlation hits whatever else sat in the payoff turn's hand.
	PerCardMarginal map[ids.CardID]CardMarginalStats
	// Histogram counts hands seen at each integer Value. Keyed by best-turn Value so
	// Min / Max can be derived without retaining every hand's value. Nil until the first
	// hand is evaluated.
	Histogram map[int]int
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

// BestTurn records the peak draw a deck saw during simulation. The chain's printout is
// produced lazily via Stats.PrintBest at render time, not eval time.
type BestTurn struct {
	Value    int
	BestLine []card.CardAssignment
}

// CardMarginalStats accumulates the with/without sums for a card's correlational marginal
// hand-value contribution. Present* covers turns where at least one copy sat in the dealt
// hand or arsenal-in slot when Best ran; Absent* covers the rest. PresentHands +
// AbsentHands equals Stats.Hands.
type CardMarginalStats struct {
	PresentTotal float64
	PresentHands int
	AbsentTotal  float64
	AbsentHands  int
}

// PresentMean returns the mean turn Value across turns where this card was present in
// the dealt hand or arsenal-in slot. Zero when the card was never present.
func (m CardMarginalStats) PresentMean() float64 {
	if m.PresentHands == 0 {
		return 0
	}
	return m.PresentTotal / float64(m.PresentHands)
}

// AbsentMean returns the mean turn Value across turns where this card was absent. Zero
// when the card was always present.
func (m CardMarginalStats) AbsentMean() float64 {
	if m.AbsentHands == 0 {
		return 0
	}
	return m.AbsentTotal / float64(m.AbsentHands)
}

// Marginal returns PresentMean - AbsentMean — the correlational hand-value lift associated
// with this card being in the turn's hand. Confounded by co-occurrence with other strong
// cards, so use as a smell test. Zero when either bucket is empty.
func (m CardMarginalStats) Marginal() float64 {
	if m.PresentHands == 0 || m.AbsentHands == 0 {
		return 0
	}
	return m.PresentMean() - m.AbsentMean()
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
