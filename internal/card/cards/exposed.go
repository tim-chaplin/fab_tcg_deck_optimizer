// Exposed — Generic Attack Reaction. Cost 0. Printed pitch variants: Blue 3.
//
// Text: "If you are **marked**, you can't play this. Target attack gets +1{p}. **Mark** the
// defending hero."

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

func (ExposedBlue) ARTargetAllowed(_ card.GameEngine, target *card.CardState, _ int8) bool {
	return target.Card.Types(nil).IsAttack()
}
func (ExposedBlue) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	self.GrantAttackReactionBuff(ge, l, 1)
	ge.MarkOpponent()
}
