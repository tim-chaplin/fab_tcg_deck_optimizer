package deck

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
)

// Per-deck evaluation result types. sim produces these during Evaluate / EvaluateAdaptive;
// cmd/fabsim renders them; internal/deckio persists them. Pure data + small derived-stat
// methods — no reference to sim or gameengine, so this file can live in v2/deck (which
// gameengine imports) without introducing a cycle.
//
// sim still owns its runtime working type (TurnSummary, holding *GameState and chain
// scratch state); BestTurn carries only the durable outcome callers actually read.

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

// BestTurn records the peak draw a deck saw during simulation. Pruned to just the fields
// production code reads:
//   - Value: the headline number cmd/fabsim prints and deckio persists.
//   - BestLine: the canonical-order assignments sim's "have-we-seen-a-best-yet" sentinel
//     gates on (len > 0). FormatBestLine reads it for one-line renders.
//   - Log: the four-section structured printout deckio persists and cmd/fabsim renders.
//
// Fields that previously lived here but had no production reader (State,
// StartingAuras/Items, SwungWeapons, IncomingDamage, StartOfTurnAuras,
// TriggersFromLastTurn, DealtHand) were removed when this package moved out of sim.
// sim's runtime TurnSummary still carries them where needed during evaluation.
type BestTurn struct {
	Value    int
	BestLine []CardAssignment
	Log      TurnLog
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
// one slice. Card is the rich v2/card.Card — sim's formatter / chain runner reads Pitch
// / Attack / Cost / Types off it directly.
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
