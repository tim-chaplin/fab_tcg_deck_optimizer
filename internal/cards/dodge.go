// Dodge — Generic Defense Reaction. Cost 0, Pitch 3, Defense 2. Only printed in Blue.

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

var DefenseReactionTypes = card.NewTypeSet(card.TypeGeneric, card.TypeDefenseReaction)

type DodgeBlue struct{}

func (DodgeBlue) ID() ids.CardID          { return ids.DodgeBlue }
func (DodgeBlue) Name() string            { return "Dodge" }
func (DodgeBlue) Cost(*sim.TurnState) int { return 0 }
func (DodgeBlue) Pitch() int              { return 3 }
func (DodgeBlue) Attack() int             { return 0 }
func (DodgeBlue) Defense() int            { return 2 }
func (DodgeBlue) Types() card.TypeSet     { return DefenseReactionTypes }
func (DodgeBlue) GoAgain() bool           { return false }
func (DodgeBlue) Play(s *sim.TurnState, self *sim.CardState) {
	n := self.DealEffectiveDefense(s)
	s.Log(self, n)
}
