// Package card defines the Card interface used by the simulator and the basic / test card
// implementations.
package card

// TurnState is the context passed to Card.Play. Cards read it to decide what effects to apply;
// the solver appends each played card to CardsPlayed after its Play method returns, so later cards
// this turn can see what was played before them.
type TurnState struct {
	// CardsPlayed is the sequence of cards played (as attacks) this turn, in order. Populated by the
	// solver, not by Play itself.
	CardsPlayed []Card
	// AuraCreated is set when a card or ability creates an aura this turn (e.g. Runechant tokens,
	// which are auras). Effects that check "if you've played or created an aura this turn" should
	// OR this with CardsPlayed containing an Aura-typed card.
	AuraCreated bool
	// CardsRemaining is the cards that will be played after the current one in the turn's ordering.
	// Populated by the solver before each Play so an effect can peek forward (e.g. Condemn to
	// Slaughter buffing the "next Runeblade attack"). Read-only from Play.
	CardsRemaining []Card
	// Overpower is set when an attack with the Overpower keyword is being played. Not yet consumed by
	// the solver — blocked damage should eventually be forwarded to the hero when Overpower is true.
	Overpower bool
}

// HasPlayedType reports whether any card played this turn has the given type in its Types() set.
func (s *TurnState) HasPlayedType(t string) bool {
	for _, c := range s.CardsPlayed {
		if c.Types()[t] {
			return true
		}
	}
	return false
}

// Card is any Flesh and Blood card that can be in a deck. Methods return the card's static profile
// plus a Play hook for on-play logic.
type Card interface {
	Name() string
	Cost() int
	Pitch() int
	// Attack is the card's base (printed) attack value. Conditional bonuses belong in Play, not here.
	Attack() int
	Defense() int
	// Types is the card's type-line descriptors as a set, e.g. {"Runeblade": true, "Action": true,
	// "Attack": true}. Implementations should return the same map each call (not a fresh literal) —
	// the map is read, never mutated.
	Types() map[string]bool
	// GoAgain reports whether playing this card grants an additional action point this turn. Cards
	// printed with "Go again" return true.
	GoAgain() bool
	// Play is called when the card is played as an attack. It returns the actual damage dealt (which
	// may differ from Attack() after conditional bonuses) and may read state to decide effects.
	Play(s *TurnState) int
}

