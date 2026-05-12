package sim

import (
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/deck"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/turnlogger"
)

// Turn-summary data shapes returned by Best: Role, CardAssignment, TurnSummary, plus the
// CarryState snapshot that captures the winning permutation's end-of-chain state. The
// deck loop adopts CarryState wholesale into the next turn's seed — no per-field
// reconstruction.

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

// CardAssignment is a single card + the role it took this turn. Hand cards produce one per
// card; an arsenal-in card contributes one with FromArsenal set so a turn fits in one slice.
type CardAssignment struct {
	Card        card.Card
	Role        Role
	FromArsenal bool
}

// CarryState is sim's snapshot of the winning permutation's end-of-chain state. The deck
// loop adopts it directly into next turn's seed — no per-field "diff" interpretation.
// Auras / Items are stored as sim's concrete pointer types; the engine boundary converts
// to gameengine.Aura / gameengine.Item when seeding the next turn's Spec.
type CarryState struct {
	Hand           []card.Card
	Deck           *deck.Deck
	Arsenal        card.Card
	Graveyard      []card.Card
	Banish         []card.Card
	Auras          []*Aura
	Items          []*Item
	CardsDrawn     int
	OpponentMarked bool
	Log            []turnlogger.LogEntry
}

// RunechantCount returns the carried Runechant token count, or zero when none are in play.
func (c *CarryState) RunechantCount() int { return auraCountByName(c.Auras, "Runechant") }

// PonderCount returns the carried Ponder token count, or zero when none are in play.
func (c *CarryState) PonderCount() int { return auraCountByName(c.Auras, "Ponder") }

// GoldCount returns the carried Gold token count, or zero when none are in play.
func (c *CarryState) GoldCount() int { return itemCountByName(c.Items, "Gold") }

// SilverCount returns the carried Silver token count, or zero when none are in play.
func (c *CarryState) SilverCount() int { return itemCountByName(c.Items, "Silver") }

// CopperCount returns the carried Copper token count, or zero when none are in play.
func (c *CarryState) CopperCount() int { return itemCountByName(c.Items, "Copper") }

// auraCountByName scans a sim-concrete aura slice for a token aura by display name.
func auraCountByName(auras []*Aura, name string) int {
	for _, a := range auras {
		if a.CardName() == name {
			return a.Count()
		}
	}
	return 0
}

// itemCountByName scans a sim-concrete item slice for a token item by display name.
func itemCountByName(items []*Item, name string) int {
	for _, i := range items {
		if i.CardName() == name {
			return i.Count()
		}
	}
	return 0
}

// TurnSummary is the result of running Best on a hand: the winning card-role assignments
// plus the CarryState snapshot the next turn inherits.
type TurnSummary struct {
	BestLine             []CardAssignment
	SwungWeapons         []string
	Value                int
	State                CarryState
	TriggersFromLastTurn []TriggerContribution
	StartOfTurnAuras     []card.Card
	DealtHand            []card.Card
	IncomingDamage       int
	Cacheable            bool
}

// TriggerContribution is one start-of-turn Aura fire: the aura that fired plus the Damage
// it credited (folded into Value) and the card (if any) the handler revealed onto the hand.
// Text is the card-authored display line — when set, the format layer renders it verbatim
// and skips the inferred "drew X into hand" / "START OF ACTION PHASE (+N)" synthesis.
type TriggerContribution struct {
	Card     card.Card
	Damage   int
	Revealed card.Card
	Text     string
}

// TurnLog is the structured record of a turn's printout, broken into four sections matching
// the natural turn boundaries. Each entry is content-only; the formatter owns indentation,
// section headers, numbering of chain events, and join.
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

// ArsenalIn returns the assignment for the card that started the turn in the arsenal.
func (t TurnSummary) ArsenalIn() (CardAssignment, bool) {
	for _, a := range t.BestLine {
		if a.FromArsenal {
			return a, true
		}
	}
	return CardAssignment{}, false
}

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
