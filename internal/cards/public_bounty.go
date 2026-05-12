// Public Bounty — Generic Action. Cost 1. Printed pitch variants: Red 1, Yellow 2, Blue 3. Defense
// 2.
//
// Text: "**Mark** target opposing hero. The next time you attack a **marked** hero this turn, the
// attack gets +N{p}. **Go again**" (Red N=3, Yellow N=2, Blue N=1.)

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
)

func publicBountyPlay(g card.GameEngine, l card.Logger, self *card.CardState, n int) {
	g.MarkOpponent()
	GrantNextCardBonusAttack(g, n, IsAttack)
}

func (PublicBountyRed) Play(g card.GameEngine, l card.Logger, self *card.CardState) {
	publicBountyPlay(g, l, self, 3)
}

func (PublicBountyYellow) Play(g card.GameEngine, l card.Logger, self *card.CardState) {
	publicBountyPlay(g, l, self, 2)
}

func (PublicBountyBlue) Play(g card.GameEngine, l card.Logger, self *card.CardState) {
	publicBountyPlay(g, l, self, 1)
}
