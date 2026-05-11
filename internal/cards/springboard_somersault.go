// Springboard Somersault — Generic Defense Reaction. Cost 0, Pitch 2, Defense 2. Only printed in
// Yellow.
// Text: "If Springboard Somersault is played from arsenal, it gains +2{d}."
//
// +2{d} when played from arsenal via sim.ArsenalDefenseBonus (docs/dev-standards.md).

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

func (SpringboardSomersaultYellow) Play(s *sim.TurnState, l sim.Logger, self *sim.CardState) {
}
func (SpringboardSomersaultYellow) ArsenalDefenseBonus() int { return 2 }
