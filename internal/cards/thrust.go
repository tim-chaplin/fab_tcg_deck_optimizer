// Thrust — Generic Attack Reaction. Cost 1. Printed pitch variants: Red 1. Defense 2.
//
// Text: "Target sword attack gains +3{p}."
//
// Predicate is "sword attack" (no "action card" qualifier), so Sword weapons qualify too.

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

func (ThrustRed) ARTargetAllowed(c sim.Card, _ int8) bool {
	t := c.Types()
	return t.Has(card.TypeSword) && t.IsAttack()
}
func (ThrustRed) Play(s sim.GameEngine, l sim.Logger, self *sim.CardState) {
	sim.GrantAttackReactionBuff(s, l, self, 3)
}
