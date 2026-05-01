// Test of Strength — Generic Block. Cost 0, Pitch 1, Defense 4. Only printed in Red.
//
// Text: "When this defends, **clash** with the attacking hero. The winner creates a Gold token."

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

type TestOfStrengthRed struct{}

func (TestOfStrengthRed) ID() ids.CardID          { return ids.TestOfStrengthRed }
func (TestOfStrengthRed) Name() string            { return "Test of Strength" }
func (TestOfStrengthRed) Cost(*sim.TurnState) int { return 0 }
func (TestOfStrengthRed) Pitch() int              { return 1 }
func (TestOfStrengthRed) Attack() int             { return 0 }
func (TestOfStrengthRed) Defense() int            { return 4 }
func (TestOfStrengthRed) Types() card.TypeSet     { return defenseReactionTypes }
func (TestOfStrengthRed) GoAgain() bool           { return false }

// not implemented: gold tokens
func (TestOfStrengthRed) NotImplemented() {}
func (TestOfStrengthRed) Play(s *sim.TurnState, self *sim.CardState) {
	n := self.DealEffectiveDefense(s)
	s.Log(self, n)
	s.ClashValue(sim.GoldTokenValue)
}
