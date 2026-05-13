// Minnowism — Generic Action. Cost 0. Printed pitch variants: Red 1, Yellow 2, Blue 3. Defense 2.
//
// Text: "The next attack action card with 3 or less base {p} you play this turn gains +N{p}. **Go
// again**" (Red N=3, Yellow N=2, Blue N=1.)

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
)

// minnowismIsTarget gates the rider on attack action cards with printed power 3 or less.
func minnowismIsTarget(_ card.GameEngine, pc *card.CardState) bool {
	return pc.Card.Types(nil).IsAttackAction() && pc.Card.Attack() <= 3
}

func (MinnowismRed) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	GrantNextCardBonusAttack(ge, 3, minnowismIsTarget)
}

func (MinnowismYellow) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	GrantNextCardBonusAttack(ge, 2, minnowismIsTarget)
}

func (MinnowismBlue) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	GrantNextCardBonusAttack(ge, 1, minnowismIsTarget)
}
