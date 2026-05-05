// Exposed — Generic Attack Reaction. Cost 0. Printed pitch variants: Blue 3.
//
// Text: "If you are **marked**, you can't play this. Target attack gets +1{p}. **Mark** the
// defending hero."

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

var exposedTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAttackReaction)

type ExposedBlue struct{}

func (ExposedBlue) ID() ids.CardID          { return ids.ExposedBlue }
func (ExposedBlue) Name() string            { return "Exposed" }
func (ExposedBlue) Cost(*sim.TurnState) int { return 0 }
func (ExposedBlue) Pitch() int              { return 3 }
func (ExposedBlue) Attack() int             { return 0 }
func (ExposedBlue) Defense() int            { return 0 }
func (ExposedBlue) Types() card.TypeSet     { return exposedTypes }
func (ExposedBlue) GoAgain() bool           { return false }
func (ExposedBlue) ARTargetAllowed(c sim.Card, _ int8) bool {
	return c.Types().IsAttack()
}
func (ExposedBlue) Play(s *sim.TurnState, self *sim.CardState) {
	sim.GrantAttackReactionBuff(s, self, 1)
	s.OpponentMarked = true
}
