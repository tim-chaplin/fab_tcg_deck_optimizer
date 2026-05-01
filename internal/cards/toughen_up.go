// Toughen Up — Generic Defense Reaction. Cost 2, Pitch 3, Defense 4. Only printed in Blue.

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

type ToughenUpBlue struct{}

func (ToughenUpBlue) ID() ids.CardID          { return ids.ToughenUpBlue }
func (ToughenUpBlue) Name() string            { return "Toughen Up" }
func (ToughenUpBlue) Cost(*sim.TurnState) int { return 2 }
func (ToughenUpBlue) Pitch() int              { return 3 }
func (ToughenUpBlue) Attack() int             { return 0 }
func (ToughenUpBlue) Defense() int            { return 4 }
func (ToughenUpBlue) Types() card.TypeSet     { return defenseReactionTypes }
func (ToughenUpBlue) GoAgain() bool           { return false }
func (ToughenUpBlue) Play(s *sim.TurnState, self *sim.CardState) {
	n := self.DealEffectiveDefense(s)
	s.Log(self, n)
}
