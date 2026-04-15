// Package card defines the Card interface used by the simulator and the
// basic / test card implementations.
package card

// TurnState is the context passed to Card.Play. Cards read it to decide
// what effects to apply; the solver appends each played card to
// CardsPlayed after its Play method returns, so later cards this turn
// can see what was played before them.
type TurnState struct {
	// CardsPlayed is the sequence of cards played (as attacks) this turn,
	// in order. Populated by the solver, not by Play itself.
	CardsPlayed []Card
}

// HasPlayedType reports whether any card played this turn has the given
// type in its Types() set.
func (s *TurnState) HasPlayedType(t string) bool {
	for _, c := range s.CardsPlayed {
		if c.Types()[t] {
			return true
		}
	}
	return false
}

// Card is any Flesh and Blood card that can be in a deck. Methods return
// the card's static profile plus a Play hook for on-play logic.
type Card interface {
	Name() string
	Cost() int
	Pitch() int
	// Attack is the card's base (printed) attack value. Conditional
	// bonuses belong in Play, not here.
	Attack() int
	Defense() int
	// Types is the card's type-line descriptors as a set, e.g.
	// {"Runeblade": true, "Action": true, "Attack": true}. Implementations
	// should return the same map each call (not a fresh literal) — the
	// map is read, never mutated.
	Types() map[string]bool
	// Play is called when the card is played as an attack. It returns
	// the actual damage dealt (which may differ from Attack() after
	// conditional bonuses) and may read state to decide effects.
	Play(s *TurnState) int
}

// --- test cards (used by unit tests and the default sim deck) ---

// TestCardBlue is a generic blue action: pitches 3, defends 3, attacks 1, costs 1.
type TestCardBlue struct{}

func (TestCardBlue) Name() string           { return "TestCardBlue" }
func (TestCardBlue) Cost() int              { return 1 }
func (TestCardBlue) Pitch() int             { return 3 }
func (TestCardBlue) Attack() int            { return 1 }
func (TestCardBlue) Defense() int           { return 3 }
func (TestCardBlue) Types() map[string]bool { return genericAttackTypes }
func (c TestCardBlue) Play(*TurnState) int  { return c.Attack() }

// TestCardRed is a generic red action: pitches 1, defends 1, attacks 3, costs 1.
type TestCardRed struct{}

func (TestCardRed) Name() string           { return "TestCardRed" }
func (TestCardRed) Cost() int              { return 1 }
func (TestCardRed) Pitch() int             { return 1 }
func (TestCardRed) Attack() int            { return 3 }
func (TestCardRed) Defense() int           { return 1 }
func (TestCardRed) Types() map[string]bool { return genericAttackTypes }

var genericAttackTypes = map[string]bool{"Generic": true, "Action": true, "Attack": true}
func (c TestCardRed) Play(*TurnState) int  { return c.Attack() }
