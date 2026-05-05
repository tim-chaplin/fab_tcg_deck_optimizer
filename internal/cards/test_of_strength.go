// Test of Strength — Generic Defense Reaction. Cost 0, Pitch 1, Defense 4. Red only.
// Text: "When this defends, **clash** with the attacking hero. The winner creates a Gold
// token."
//
// Loss is modelled as -1 Value: opponent gets the Gold token, which we approximate as
// one resource handed over.

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
func (TestOfStrengthRed) Types() card.TypeSet     { return DefenseReactionTypes }
func (TestOfStrengthRed) GoAgain() bool           { return false }

func (TestOfStrengthRed) CreatesItem() sim.TokenType { return sim.TokenTypeGold }

func (TestOfStrengthRed) Play(s *sim.TurnState, self *sim.CardState) {
	n := self.DealEffectiveDefense(s)
	s.Log(self, n)
	s.Clash(
		func() {
			s.CreateGold(1)
			s.LogRider(self, 0, "Clash win created a gold token")
		},
		func() {
			// AddValue clamps negatives, so write directly: opponent gains the Gold token,
			// netting us roughly one resource of opposing tempo.
			s.Value--
			s.LogRider(self, -1, "Clash loss conceded gold to opponent")
		},
	)
}
