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
}

// TurnSummary is the result of running Best on a hand: the winning card-role assignments
// plus the CarryState snapshot the next turn inherits.
type TurnSummary struct {
	// BestLine is the winning partition. Hand cards come first in canonical (post-sort)
	// order; the previous-turn arsenal card, if any, is the last entry with FromArsenal=true.
	BestLine []CardAssignment
	// SwungWeapons names the weapons swung this turn in the winning permutation. Weapons have
	// no BestLine entry (they're not hand cards), so the printout reads them from here.
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
