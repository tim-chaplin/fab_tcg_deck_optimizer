// Nip at the Heels — Generic Attack Reaction. Cost 0. Printed pitch variants: Blue 3.
// Defense 3.
//
// Text: "Target attack with 3 or less base {p} gets +1{p}."
//
// Predicate accepts attack action cards and weapons (just "attack"); the ≤ 3 gate reads
// printed Attack(), not the post-buff total.

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

func (NipAtTheHeelsBlue) ARTargetAllowed(_ card.GameEngine, target *card.CardState, _ int8) bool {
	return target.Card.Types(nil).IsAttack() && target.Card.Attack() <= 3
}
func (NipAtTheHeelsBlue) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	self.GrantAttackReactionBuff(ge, l, 1)
}
