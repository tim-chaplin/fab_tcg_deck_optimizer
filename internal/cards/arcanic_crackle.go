// Arcanic Crackle — Runeblade Action - Attack. Cost 0, Defense 3, Arcane 1.
// Printed power: Red 3, Yellow 2, Blue 1.
// Text: "Deal 1 arcane damage to target hero."
//
// The printed 1 arcane is added to combat damage (both hit the same target). Play also sets
// ArcaneDamageDealt so same-turn triggers reading "if you've dealt arcane damage this turn"
// fire.

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

func (ArcanicCrackleRed) Play(s sim.GameEngine, l sim.Logger, self *sim.CardState) {
	s.DealArcaneDamage(l, self, 1)
}

func (ArcanicCrackleYellow) Play(s sim.GameEngine, l sim.Logger, self *sim.CardState) {
	s.DealArcaneDamage(l, self, 1)
}

func (ArcanicCrackleBlue) Play(s sim.GameEngine, l sim.Logger, self *sim.CardState) {
	s.DealArcaneDamage(l, self, 1)
}
