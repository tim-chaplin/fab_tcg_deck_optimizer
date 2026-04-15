// Package card defines the basic Card type used by the simulator.
package card

// Card is a generic Flesh and Blood action card. The simulator currently
// only models the four numbers needed to value a hand in isolation.
type Card struct {
	Name   string
	Cost   int
	Pitch  int
	Attack int
	Defend int
}

// TestCardBlue is a generic blue action: pitches 3, defends 3, attacks 1, costs 1.
var TestCardBlue = Card{Name: "Blue", Cost: 1, Pitch: 3, Attack: 1, Defend: 3}

// TestCardRed is a generic red action: pitches 1, defends 1, attacks 3, costs 1.
var TestCardRed = Card{Name: "Red", Cost: 1, Pitch: 1, Attack: 3, Defend: 1}
