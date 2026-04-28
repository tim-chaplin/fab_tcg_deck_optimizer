package deck

// Aggregate statistics accumulated while simulating a deck: total / per-cycle / per-card tallies,
// the single best turn ever seen, and a histogram of hand values that supports Min / Max /
// Median without retaining every individual hand.

import (
	"sort"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/hand"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
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
	// PerCardMarginal carries a coarse correlational view of each card's hand-value impact:
	// for every unique card ID in d.Cards, the mean turn Value across turns where that card
	// was in the dealt hand (or arsenal-in slot) vs turns where it wasn't. The gap between
	// the two means is a smell test — cards whose presence shifts hand value far more than
	// their printed face value would suggest are candidates for buggy or oversimplified
	// implementations.
	//
	// Future problem: cards whose printed effect pays off on a LATER turn (auras, drawn-card
	// payoffs that resolve next turn) often won't surface here — the source card has rotated
	// out of the hand by the time its value lands, so the correlation hits whichever cards
	// happened to share the payoff turn instead. A regression-based estimator over per-hand
	// presence vectors would credit such effects more cleanly.
	PerCardMarginal map[ids.CardID]CardMarginalStats
	// Histogram counts hands seen at each integer Value. Keyed by TurnSummary.Value so Min /
	// Median can be derived without retaining every hand's value. Nil until the first hand is
	// evaluated.
	Histogram map[int]int
}

// BestTurn records a single hand and its optimal turn — the peak draw a deck saw during
// simulation. Summary.BestLine carries the cards and roles in canonical order; Log is the
// structured per-section trace assembled at end of EvaluateWith via hand.BuildTurnLog and
// round-tripped through the JSON layer verbatim. fabsim's print path renders Log via
// hand.FormatTurnLog so saved decks produce the same output as live runs.
type BestTurn struct {
	Summary hand.TurnSummary
	// StartingRunechants is the Runechant count carried in from the previous turn when this hand
	// was played. Only meaningful for Runeblade heroes.
	StartingRunechants int
	// Log is the four-section structured record (StartOfTurn / MyTurn / OpponentTurn /
	// EndOfTurn) of the best turn's printout. Each entry is content-only; the formatter
	// owns indentation, section headers, and chain numbering. EvaluateWith populates it
	// once at end of run via hand.BuildTurnLog.
	Log hand.TurnLog
}

// CardMarginalStats accumulates the with/without sums needed to compute a card's correlational
// marginal hand-value contribution. PresentTotal / PresentHands cover turns where at least one
// copy of the card sat in the dealt hand or arsenal-in slot when hand.Best ran; AbsentTotal /
// AbsentHands cover the rest. PresentHands + AbsentHands always equals the deck's total Hands.
type CardMarginalStats struct {
	PresentTotal float64
	PresentHands int
	AbsentTotal  float64
	AbsentHands  int
}

// PresentMean returns the mean turn Value across turns where this card was present in the
// dealt hand or arsenal-in slot. Zero when the card was never present.
func (m CardMarginalStats) PresentMean() float64 {
	if m.PresentHands == 0 {
		return 0
	}
	return m.PresentTotal / float64(m.PresentHands)
}

// AbsentMean returns the mean turn Value across turns where this card was absent. Zero when
// the card was always present.
func (m CardMarginalStats) AbsentMean() float64 {
	if m.AbsentHands == 0 {
		return 0
	}
	return m.AbsentTotal / float64(m.AbsentHands)
}

// Marginal returns PresentMean - AbsentMean — the correlational hand-value lift associated
// with this card being in the turn's hand. Positive means hands containing the card score
// higher on average; negative means lower. Confounded by co-occurrence with other strong
// cards, so use as a smell test, not a precise per-card valuation. Zero when either bucket
// is empty.
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
