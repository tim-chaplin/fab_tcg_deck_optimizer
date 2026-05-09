// Exposed — Generic Attack Reaction. Cost 0. Printed pitch variants: Blue 3.
//
// Text: "If you are **marked**, you can't play this. Target attack gets +1{p}. **Mark** the
// defending hero."

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

func (ExposedBlue) ARTargetAllowed(c sim.Card, _ int8) bool {
	return c.Types().IsAttack()
}
func (ExposedBlue) Play(s *sim.TurnState, self *sim.CardState) {
	sim.GrantAttackReactionBuff(s, self, 1)
	s.OpponentMarked = true
}
