// Lunging Press — Generic Attack Reaction. Cost 0. Printed pitch variants: Blue 3. Defense 2.
//
// Text: "Target attack action card gains +1{p}."
//
// Predicate is "attack action card" (not "attack"), so weapon swings are excluded.

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

func (LungingPressBlue) ARTargetAllowed(_ card.GameEngine, target *card.CardState, _ int8) bool {
	return target.Card.Types(nil).IsAttackAction()
}
func (LungingPressBlue) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	self.GrantAttackReactionBuff(ge, l, 1)
}
