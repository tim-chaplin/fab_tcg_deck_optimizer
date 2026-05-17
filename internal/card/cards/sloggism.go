// Sloggism — Generic Action. Cost 3. Printed pitch variants: Red 1, Yellow 2, Blue 3. Defense 2.
//
// Text: "The next attack action card with cost 2 or greater you play this turn gains +N{p}. **Go
// again**" (Red N=6, Yellow N=5, Blue N=4.)

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

// sloggismIsTarget gates the rider on attack action cards whose cost is 2 or more.
func sloggismIsTarget(ge card.GameEngine, pc *card.CardState) bool {
	return pc.Card.Types(nil).IsAttackAction() && pc.Card.Cost(ge) >= 2
}

func (SloggismRed) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	GrantNextCardBonusAttack(ge, 6, sloggismIsTarget)
}

func (SloggismYellow) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	GrantNextCardBonusAttack(ge, 5, sloggismIsTarget)
}

func (SloggismBlue) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	GrantNextCardBonusAttack(ge, 4, sloggismIsTarget)
}
