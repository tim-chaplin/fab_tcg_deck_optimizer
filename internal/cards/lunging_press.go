// Lunging Press — Generic Attack Reaction. Cost 0. Printed pitch variants: Blue 3. Defense 2.
//
// Text: "Target attack action card gains +1{p}."
//
// Predicate is "attack action card" (not "attack"), so weapon swings are excluded.

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

func (LungingPressBlue) ARTargetAllowed(c sim.Card, _ int8) bool {
	return c.Types().IsAttackAction()
}
func (LungingPressBlue) Play(s *sim.TurnState, l sim.Logger, self *sim.CardState) {
	sim.GrantAttackReactionBuff(s, l, self, 1)
}
