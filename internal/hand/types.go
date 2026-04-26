package hand

// Turn-summary data shapes returned by Best: Role, CardAssignment, TurnSummary, plus the
// CarryState snapshot that captures the winning permutation's end-of-chain TurnState. The
// deck loop adopts CarryState wholesale into the next turn's seed — no per-field
// reconstruction.

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

// Role is what a card did on a given turn cycle.
type Role uint8

const (
	Pitch Role = iota
	Attack
	Defend
	Held
	// Arsenal marks the card placed into the arsenal slot at end of turn. Contributes
	// nothing to this turn's Value (it scores on the future turn it's played) and carries
	// across via TurnSummary.State.Arsenal.
	Arsenal
)

// CardAssignment is a single card + the role it took this turn. Hand cards produce one per
// card; an arsenal-in card contributes one with FromArsenal set so a turn fits in one slice.
type CardAssignment struct {
	Card        card.Card
	Role        Role
	FromArsenal bool
}

// CarryState is the slice of TurnState that carries across the turn boundary — the winning
// permutation's end-of-chain snapshot. The deck loop adopts CarryState directly into next
// turn's TurnState seed; no per-field "diff" interpretation. Cards that mutate state during
// Play write through to whichever slice ends up here.
type CarryState struct {
	// Hand is the cards in hand at end of chain — partition Held cards plus anything
	// tutored or drawn that didn't get played. Becomes next turn's Held prefix.
	Hand []card.Card
	// Deck is the deck at end of chain (top-to-bottom). Reflects every mid-chain mutation
	// (DrawOne pops, tutor removals, alt-cost prepends). Becomes next turn's draw pile.
	Deck []card.Card
	// Arsenal is the arsenal slot at end of chain. Set by the partition (arsenal-in stayed),
	// or filled post-hoc by promoting a Hand card when the slot is empty.
	Arsenal card.Card
	// Graveyard is every card that landed in the graveyard this turn — played hand cards,
	// tutored-and-played cards, AuraTriggers that destroyed themselves.
	Graveyard []card.Card
	// Banish is cards moved into the banished zone this turn.
	Banish []card.Card
	// Runechants is the live token count at end of chain. Carries across.
	Runechants int
	// AuraTriggers is the surviving AuraTrigger set at end of chain. Carries across.
	AuraTriggers []card.AuraTrigger
	// Log is the per-event chain trace of the winning permutation — one entry per Play, hero
	// trigger, aura trigger, ephemeral trigger, weapon swing. Stored as raw card.LogEntry
	// structs to defer fmt.Sprintf cost until BuildTurnLog runs at end of EvaluateWith,
	// keeping the snapshot path allocation-light (only the winning permutation's log gets
	// formatted, and only when the deck-level Best actually changes). Doesn't carry across
	// turns semantically but rides on CarryState's snapshot mechanism for free.
	Log []card.LogEntry
}

// TurnSummary is the result of running Best on a hand: the winning card-role assignments
// plus the CarryState snapshot the next turn inherits.
type TurnSummary struct {
	// BestLine is the winning partition. Hand cards come first in canonical (post-sort)
	// order; the previous-turn arsenal card, if any, is the last entry with FromArsenal=true.
	BestLine []CardAssignment
	// SwungWeapons names the weapons swung this turn in the winning permutation. Weapons
	// resolve through the dispatcher and log "WEAPON ATTACK" lines into State.Log, so the
	// numbered printout reads weapon swings from there. SwungWeapons stays on the summary
	// for the deckio JSON round-trip — Marshal serialises it under "weapons" so a reloaded
	// best turn still names the swung weapons even when State.Log is absent.
	SwungWeapons []string
	// Value is the turn's total score (damage dealt + damage prevented).
	Value int
	// State is the winning permutation's end-of-chain CarryState. The deck loop copies
	// every field into next turn's seed.
	State CarryState
	// TriggersFromLastTurn records the AuraTriggers whose start-of-turn handlers fired at
	// the top of this turn, each with the damage-equivalent it credited.
	TriggersFromLastTurn []TriggerContribution
	// StartOfTurnAuras lists the aura cards that were in play at the top of this turn — one
	// entry per AuraTrigger carried in from the previous turn.
	StartOfTurnAuras []card.Card
}

// TriggerContribution is one start-of-turn AuraTrigger fire: the aura that fired plus the
// Damage it credited (folded into Value) and the card (if any) the handler revealed onto
// the hand.
type TriggerContribution struct {
	Card     card.Card
	Damage   int
	Revealed card.Card
}

// TurnLog is the structured record of a turn's printout, broken into four sections matching
// the natural turn boundaries. Each entry is content-only — "Hocus Pocus [B]: PITCH" —
// so the formatter owns indentation, section headers, numbering of chain events, and
// join. JSON serializes the struct directly so the on-disk shape is browsable / diffable
// per section.
type TurnLog struct {
	// StartOfTurn captures the turn's starting state and any start-of-turn trigger fires:
	// dealt hand, arsenal-in card, auras / runechants in play, then the carryover
	// AuraTrigger handler effects (Sigil reveals, damage credits). Mixed format — informational
	// lines like "Hand: A, B, C, D" sit alongside event lines like "Sigil of the Arknight
	// [B]: drew X into hand"; the formatter renders both unnumbered.
	StartOfTurn []string `json:"start_of_turn,omitempty"`
	// MyTurn is the numbered entries for the "My turn:" section: attack-phase pitches
	// followed by the chain (Play / hero trigger / aura trigger / ephemeral trigger / weapon
	// swing lines, in resolution order).
	MyTurn []string `json:"my_turn,omitempty"`
	// OpponentTurn is the numbered entries for the "Opponent's turn:" section: defense-phase
	// pitches, plain blocks, and Defense Reactions, in that order.
	OpponentTurn []string `json:"opponent_turn,omitempty"`
	// EndOfTurn captures the turn's ending state: the cards in hand, the arsenal slot's
	// contents, and the auras / runechants surviving into the next turn. Mirrors
	// StartOfTurn's mixed-informational format ("Hand: A, B", "Arsenal: X (stayed)",
	// "Auras: Y, 1 Runechant"); the formatter renders unnumbered.
	EndOfTurn []string `json:"end_of_turn,omitempty"`
}

// IsEmpty reports whether all four sections are empty — true for an unscored deck or a
// hand where Best returned without a winning line. Marshal / Unmarshal / printBestTurn use
// this to short-circuit the best-turn block.
func (l TurnLog) IsEmpty() bool {
	return len(l.StartOfTurn) == 0 && len(l.MyTurn) == 0 &&
		len(l.OpponentTurn) == 0 && len(l.EndOfTurn) == 0
}

// ArsenalIn returns the assignment for the card that started the turn in the arsenal, if
// any.
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
