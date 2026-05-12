// Starting Stake — Generic Action. Cost 0, Pitch 2, Defense 3. Yellow only.
// Text: "If you control no Gold tokens, create a Gold token."

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
)

func (StartingStakeYellow) Play(g card.GameEngine, l card.Logger, self *card.CardState) {
	if g.GoldCount() == 0 {
		g.CreateGold(1)
		l.AppendPostTrigger(self.Card.DisplayName(), "Created a gold token", 0)
	}
}
