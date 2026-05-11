// Nip at the Heels — Generic Attack Reaction. Cost 0. Printed pitch variants: Blue 3.
// Defense 3.
//
// Text: "Target attack with 3 or less base {p} gets +1{p}."
//
// Predicate accepts attack action cards and weapons (just "attack"); the ≤ 3 gate reads
// printed Attack(), not the post-buff total.

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

func (NipAtTheHeelsBlue) ARTargetAllowed(c sim.Card, _ int8) bool {
	return c.Types().IsAttack() && c.Attack() <= 3
}
func (NipAtTheHeelsBlue) Play(s *sim.TurnState, l sim.Logger, self *sim.CardState) {
	sim.GrantAttackReactionBuff(s, l, self, 1)
}
