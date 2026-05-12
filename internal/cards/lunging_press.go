// Lunging Press — Generic Attack Reaction. Cost 0. Printed pitch variants: Blue 3. Defense 2.
//
// Text: "Target attack action card gains +1{p}."
//
// Predicate is "attack action card" (not "attack"), so weapon swings are excluded.

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
)

func (LungingPressBlue) ARTargetAllowed(_ card.GameEngine, c card.Card, _ int8) bool {
	return c.Types(nil).IsAttackAction()
}
func (LungingPressBlue) Play(g card.GameEngine, l card.Logger, self *card.CardState) {
	self.GrantAttackReactionBuff(g, l, 1)
}
