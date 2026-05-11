// Test of Strength — Generic Defense Reaction. Cost 0, Pitch 1, Defense 4. Red only.
// Text: "When this defends, **clash** with the attacking hero. The winner creates a Gold
// token."
//
// Loss is modelled as -1 Value: opponent gets the Gold token, which we approximate as
// one resource handed over.

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

func (TestOfStrengthRed) Play(s *sim.TurnState, l sim.Logger, self *sim.CardState) {
	n := self.DealEffectiveDefense(s)
	self.Log(l, n)
	s.Clash(
		func() {
			s.CreateGold(1)
			self.LogRider(l, 0, "Clash win created a gold token")
		},
		func() {
			// AddValue clamps negatives, so write directly: opponent gains the Gold token,
			// netting us roughly one resource of opposing tempo.
			s.Value--
			self.LogRider(l, -1, "Clash loss conceded gold to opponent")
		},
	)
}
