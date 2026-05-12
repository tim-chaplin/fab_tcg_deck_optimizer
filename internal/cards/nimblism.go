// Nimblism — Generic Action. Cost 0. Printed pitch variants: Red 1, Yellow 2, Blue 3. Defense 2.
//
// Text: "The next attack action card with cost 1 or less you play this turn gains +N{p}. **Go
// again**" (Red N=3, Yellow N=2, Blue N=1.)

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
)

// nimblismIsTarget gates the rider on attack action cards whose cost is 1 or less.
func nimblismIsTarget(g card.GameEngine, pc *card.CardState) bool {
	return pc.Card.Types(nil).IsAttackAction() && pc.Card.Cost(g) <= 1
}

func (NimblismRed) Play(g card.GameEngine, l card.Logger, self *card.CardState) {
	GrantNextCardBonusAttack(g, 3, nimblismIsTarget)
}

func (NimblismYellow) Play(g card.GameEngine, l card.Logger, self *card.CardState) {
	GrantNextCardBonusAttack(g, 2, nimblismIsTarget)
}

func (NimblismBlue) Play(g card.GameEngine, l card.Logger, self *card.CardState) {
	GrantNextCardBonusAttack(g, 1, nimblismIsTarget)
}
