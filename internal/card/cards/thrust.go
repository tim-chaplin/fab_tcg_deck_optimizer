// Thrust — Generic Attack Reaction. Cost 1. Printed pitch variants: Red 1. Defense 2.
//
// Text: "Target sword attack gains +3{p}."
//
// Predicate is "sword attack" (no "action card" qualifier), so Sword weapons qualify too.

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

func (ThrustRed) ARTargetAllowed(_ card.GameEngine, target *card.CardState, _ int8) bool {
	t := target.Card.Types(nil)
	return t.Has(card.TypeSword) && t.IsAttack()
}
func (ThrustRed) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	self.GrantAttackReactionBuff(ge, l, 3)
}
