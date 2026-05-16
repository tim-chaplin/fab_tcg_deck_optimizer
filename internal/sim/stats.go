package sim

// Aggregate statistics accumulated while simulating a deck: total / per-cycle / per-card tallies,
// the single best turn ever seen, and a histogram of hand values that supports Min / Max without
// retaining every individual hand.

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/aura"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/item"
)

// DeckStats holds aggregate hand-value statistics across all simulated runs of a Deck.
type DeckStats struct {
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
	// for every unique card ID in the deck, the mean turn Value across turns where that card
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
	// Max can be derived without retaining every hand's value. Nil until the first hand is
	// evaluated.
	Histogram map[int]int
}

// BestTurn records a single hand and its optimal turn — the peak draw a deck saw during
// simulation. Summary.BestLine carries the cards and roles in canonical order; Log is the
// structured per-section trace assembled at end of Evaluate via BuildTurnLog and
// round-tripped through the JSON layer verbatim. fabsim's print path renders Log via
// FormatTurnLog so saved decks produce the same output as live runs.
type BestTurn struct {
	Summary TurnSummary
	// StartingAuras is the carryover aura set entering this turn — sigils, incantations,
	// and token auras in play when the hand was dealt.
	StartingAuras []*aura.Aura
	// StartingItems is the carryover item set entering this turn — Gold tokens (and
	// future card items) in play when the hand was dealt.
	StartingItems []*item.Item
	// Log is the four-section structured record (StartOfTurn / MyTurn / OpponentTurn /
	// EndOfTurn) of the best turn's printout. Each entry is content-only; the formatter
	// owns indentation, section headers, and chain numbering. Evaluate populates it
	// once at end of run via BuildTurnLog.
	Log TurnLog
}

// CardMarginalStats accumulates the with/without sums needed to compute a card's correlational
// marginal hand-value contribution. PresentTotal / PresentHands cover turns where at least one
// copy of the card sat in the dealt hand or arsenal-in slot when Best ran; AbsentTotal /
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
func (s DeckStats) Mean() float64 {
	if s.Hands == 0 {
		return 0
	}
	return s.TotalValue / float64(s.Hands)
}

// Min returns the lowest Value any simulated hand produced. Zero when no hands have been seen.
func (s DeckStats) Min() int {
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
func (s DeckStats) Max() int {
	m := 0
	for v := range s.Histogram {
		if v > m {
			m = v
		}
	}
	return m
}
