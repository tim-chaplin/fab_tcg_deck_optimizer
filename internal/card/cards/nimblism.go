// Nimblism — Generic Action. Cost 0. Printed pitch variants: Red 1, Yellow 2, Blue 3. Defense 2.
//
// Text: "The next attack action card with cost 1 or less you play this turn gains +N{p}. **Go
// again**" (Red N=3, Yellow N=2, Blue N=1.)

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

// nimblismIsTarget gates the rider on attack action cards whose cost is 1 or less.
func nimblismIsTarget(ge card.GameEngine, pc *card.CardState) bool {
	return pc.Card.Types(nil).IsAttackAction() && pc.Card.Cost(ge) <= 1
}

func (NimblismRed) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	GrantNextCardBonusAttack(ge, 3, nimblismIsTarget)
}

func (NimblismYellow) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	GrantNextCardBonusAttack(ge, 2, nimblismIsTarget)
}

func (NimblismBlue) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	GrantNextCardBonusAttack(ge, 1, nimblismIsTarget)
}
