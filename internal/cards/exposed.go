// Exposed — Generic Attack Reaction. Cost 0. Printed pitch variants: Blue 3.
//
// Text: "If you are **marked**, you can't play this. Target attack gets +1{p}. **Mark** the
// defending hero."

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
)

func (ExposedBlue) ARTargetAllowed(c card.Card, _ int8) bool {
	return c.Types().IsAttack()
}
func (ExposedBlue) Play(s card.GameEngine, l card.Logger, self *card.CardState) {
	card.GrantAttackReactionBuff(s, l, self, 1)
	s.SetOpponentMarked(true)
}
