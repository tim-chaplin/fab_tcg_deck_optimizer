// Test of Strength — Generic Defense Reaction. Cost 0, Pitch 1, Defense 4. Red only.
// Text: "When this defends, **clash** with the attacking hero. The winner creates a Gold
// token."
//
// Loss is modelled as -1 Value: opponent gets the Gold token, which we approximate as
// one resource handed over.

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
)

func (TestOfStrengthRed) Play(g card.GameEngine, l card.Logger, self *card.CardState) {
	g.Clash(
		func() {
			g.CreateGold(1)
			l.AppendPostTrigger(self.Card.DisplayName(), "Clash win created a gold token", 0)
		},
		func() {
			g.AddValue(-1)
			l.AppendPostTrigger(self.Card.DisplayName(), "Clash loss conceded gold to opponent", -1)
		},
	)
}
