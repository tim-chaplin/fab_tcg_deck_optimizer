// Package deckstats holds the serializable per-deck evaluation result types — the data
// shape sim produces during Evaluate / EvaluateAdaptive and that fabsim renders and
// deckio persists. Pure data + small derived-stat methods (Mean, Min, Max, …); no
// reference to sim or gameengine, so the package can be imported by either.
//
// sim still owns the runtime working types (TurnSummary with its *GameState pointer,
// per-permutation context, etc.); BestTurn lifts only the fields callers / persisted
// output actually need off TurnSummary so the persisted shape is decoupled from the
// runtime shape.
package deckstats

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/aura"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/gameengine"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/item"
)

// DeckStats is the aggregate hand-value statistics across all simulated runs of a Deck.
type DeckStats struct {
	Runs        int
	Hands       int
	TotalValue  float64
	FirstCycle  CycleStats
	SecondCycle CycleStats
	// Best is the single highest-value hand seen across all runs (ties broken by first
	// occurrence). BestLine is in canonical (post-sort) order. Zero-valued if no hands
	// have been evaluated.
	Best BestTurn
	// PerCardMarginal carries a coarse correlational view of each card's hand-value
	// impact: for every unique card ID in the deck, the mean turn Value across turns
	// where that card was in the dealt hand (or arsenal-in slot) vs turns where it
	// wasn't. The gap between the two means is a smell test — cards whose presence
	// shifts hand value far more than their printed face value would suggest are
	// candidates for buggy or oversimplified implementations.
	//
	// Future problem: cards whose printed effect pays off on a LATER turn (auras,
	// drawn-card payoffs that resolve next turn) often won't surface here — the source
	// card has rotated out of the hand by the time its value lands, so the correlation
	// hits whichever cards happened to share the payoff turn instead. A regression-based
	// estimator over per-hand presence vectors would credit such effects more cleanly.
	PerCardMarginal map[ids.CardID]CardMarginalStats
	// Histogram counts hands seen at each integer Value. Keyed by best-turn Value so
	// Min / Max can be derived without retaining every hand's value. Nil until the first
	// hand is evaluated.
	Histogram map[int]int
}

// Mean returns the overall arithmetic mean hand value.
func (s DeckStats) Mean() float64 {
	if s.Hands == 0 {
		return 0
	}
	return s.TotalValue / float64(s.Hands)
}

// Min returns the lowest Value any simulated hand produced. Zero when no hands have been
// seen.
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

// Max returns the highest Value any simulated hand produced. Zero when no hands have
// been seen.
func (s DeckStats) Max() int {
	m := 0
	for v := range s.Histogram {
		if v > m {
			m = v
		}
	}
	return m
}

// BestTurn records the peak draw a deck saw during simulation: the turn's value, the
// canonical-order card / role assignments, the starting carryover (auras / items / aura-
// fires from last turn / start-of-turn aura snapshot), and the rendered structured Log.
// State is the post-chain *gameengine.GameState the winning chain produced — preserved
// in-memory so callers can inspect Hand / Deck / Graveyard / Arsenal after Evaluate; not
// serialized (the engine has no JSON encoding; the Log is the persisted shape).
type BestTurn struct {
	Value                int
	BestLine             []CardAssignment
	SwungWeapons         []string
	TriggersFromLastTurn []TriggerContribution
	StartOfTurnAuras     []card.Card
	IncomingDamage       int
	// State is the post-chain engine state the winning permutation produced. Lives only
	// in memory; deckio marshals around it.
	State *gameengine.GameState `json:"-"`
	// StartingAuras is the carryover aura set entering this turn — sigils, incantations,
	// and token auras in play when the hand was dealt.
	StartingAuras []*aura.Aura
	// StartingItems is the carryover item set entering this turn — Gold tokens (and
	// future card items) in play when the hand was dealt.
	StartingItems []*item.Item
	// Log is the four-section structured record (StartOfTurn / MyTurn / OpponentTurn /
	// EndOfTurn) of the best turn's printout. Each entry is content-only; the formatter
	// owns indentation, section headers, and chain numbering.
	Log TurnLog
}

// Role is what a card did on a given turn cycle.
type Role uint8

const (
	Pitch Role = iota
	Attack
	Defend
	Held
	// Arsenal marks the card placed into the arsenal slot at end of turn.
	Arsenal
)

// String returns a human-readable role name.
func (r Role) String() string {
	switch r {
	case Pitch:
		return "PITCH"
	case Attack:
		return "ATTACK"
	case Defend:
		return "DEFEND"
	case Held:
		return "HELD"
	case Arsenal:
		return "ARSENAL"
	}
	return "UNKNOWN"
}

// CardAssignment is a single card + the role it took this turn. Hand cards produce one
// per card; an arsenal-in card contributes one with FromArsenal set so a turn fits in
// one slice.
type CardAssignment struct {
	Card        card.Card
	Role        Role
	FromArsenal bool
}

// TriggerContribution is one start-of-turn Aura fire: the aura that fired plus the
// Damage it credited (folded into Value) and the card (if any) the handler revealed onto
// the hand. Text is the card-authored display line — when set, the format layer renders
// it verbatim and skips the inferred "drew X into hand" / "START OF ACTION PHASE (+N)"
// synthesis.
type TriggerContribution struct {
	Card     card.Card
	Damage   int
	Revealed card.Card
	Text     string
}

// CardMarginalStats accumulates the with/without sums needed to compute a card's
// correlational marginal hand-value contribution. PresentTotal / PresentHands cover
// turns where at least one copy of the card sat in the dealt hand or arsenal-in slot
// when Best ran; AbsentTotal / AbsentHands cover the rest. PresentHands + AbsentHands
// always equals the deck's total Hands.
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

// Marginal returns PresentMean - AbsentMean — the correlational hand-value lift
// associated with this card being in the turn's hand. Positive means hands containing
// the card score higher on average; negative means lower. Confounded by co-occurrence
// with other strong cards, so use as a smell test, not a precise per-card valuation.
// Zero when either bucket is empty.
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

// TurnLog is the structured record of a turn's printout, broken into four sections
// matching the natural turn boundaries. Each entry is content-only; the formatter owns
// indentation, section headers, numbering of chain events, and join.
type TurnLog struct {
	StartOfTurn  []string `json:"start_of_turn,omitempty"`
	MyTurn       []string `json:"my_turn,omitempty"`
	OpponentTurn []string `json:"opponent_turn,omitempty"`
	EndOfTurn    []string `json:"end_of_turn,omitempty"`
}

// IsEmpty reports whether all four sections are empty.
func (l TurnLog) IsEmpty() bool {
	return len(l.StartOfTurn) == 0 && len(l.MyTurn) == 0 &&
		len(l.OpponentTurn) == 0 && len(l.EndOfTurn) == 0
}
