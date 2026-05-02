// Nip at the Heels — Generic Attack Reaction. Cost 0. Printed pitch variants: Blue 3.
// Defense 3.
//
// Text: "Target attack with 3 or less base {p} gets +1{p}."
//
// "Target attack" includes weapon attacks; the base-power gate (≤ 3) reads the printed
// Attack() value, not the post-buff total. The predicate accepts both attack action cards
// and weapons whose printed power is ≤ 3.

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

var nipAtTheHeelsTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAttackReaction)

type NipAtTheHeelsBlue struct{}

func (NipAtTheHeelsBlue) ID() ids.CardID          { return ids.NipAtTheHeelsBlue }
func (NipAtTheHeelsBlue) Name() string            { return "Nip at the Heels" }
func (NipAtTheHeelsBlue) Cost(*sim.TurnState) int { return 0 }
func (NipAtTheHeelsBlue) Pitch() int              { return 3 }
func (NipAtTheHeelsBlue) Attack() int             { return 0 }
func (NipAtTheHeelsBlue) Defense() int            { return 3 }
func (NipAtTheHeelsBlue) Types() card.TypeSet     { return nipAtTheHeelsTypes }
func (NipAtTheHeelsBlue) GoAgain() bool           { return false }
func (NipAtTheHeelsBlue) ARTargetAllowed(c sim.Card) bool {
	return c.Types().IsAttack() && c.Attack() <= 3
}
func (NipAtTheHeelsBlue) Play(s *sim.TurnState, self *sim.CardState) {
	sim.GrantAttackReactionBuff(s, NipAtTheHeelsBlue{}.ARTargetAllowed, 1)
	s.Log(self, 0)
}
