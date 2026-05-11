// Vexing Malice — Runeblade Action - Attack. Cost 1, Defense 3, Arcane 2.
// Printed power: Red 3, Yellow 2, Blue 1.
// Text: "Deal 2 arcane damage to target hero."
//
// The printed 2 arcane is added to combat damage (both hit the same target). Play also sets
// ArcaneDamageDealt so same-turn triggers keyed on that flag fire.

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

func (VexingMaliceRed) Play(s sim.GameEngine, l sim.Logger, self *sim.CardState) {
	s.DealArcaneDamage(l, self, 2)
}

func (VexingMaliceYellow) Play(s sim.GameEngine, l sim.Logger, self *sim.CardState) {
	s.DealArcaneDamage(l, self, 2)
}

func (VexingMaliceBlue) Play(s sim.GameEngine, l sim.Logger, self *sim.CardState) {
	s.DealArcaneDamage(l, self, 2)
}
