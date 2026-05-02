// Lunging Press — Generic Attack Reaction. Cost 0. Printed pitch variants: Blue 3. Defense 2.
//
// Text: "Target attack action card gains +1{p}."
//
// Targets exclude weapon attacks: the printed text says "attack action card", not "attack",
// so a chain whose only attack candidate is a weapon swing can't legally play this AR.

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

var lungingPressTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAttackReaction)

type LungingPressBlue struct{}

func (LungingPressBlue) ID() ids.CardID          { return ids.LungingPressBlue }
func (LungingPressBlue) Name() string            { return "Lunging Press" }
func (LungingPressBlue) Cost(*sim.TurnState) int { return 0 }
func (LungingPressBlue) Pitch() int              { return 3 }
func (LungingPressBlue) Attack() int             { return 0 }
func (LungingPressBlue) Defense() int            { return 2 }
func (LungingPressBlue) Types() card.TypeSet     { return lungingPressTypes }
func (LungingPressBlue) GoAgain() bool           { return false }
func (LungingPressBlue) ARTargetAllowed(c sim.Card) bool {
	return c.Types().IsAttackAction()
}
func (LungingPressBlue) Play(s *sim.TurnState, self *sim.CardState) {
	sim.GrantAttackReactionBuff(s, LungingPressBlue{}.ARTargetAllowed, 1)
	s.Log(self, 0)
}
