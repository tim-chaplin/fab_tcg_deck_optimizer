// Nimblism — Generic Action. Cost 0. Printed pitch variants: Red 1, Yellow 2, Blue 3. Defense 2.
//
// Text: "The next attack action card with cost 1 or less you play this turn gains +N{p}. **Go
// again**" (Red N=3, Yellow N=2, Blue N=1.)

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

// nimblismIsTarget gates the rider on attack action cards whose cost is 1 or less.
func nimblismIsTarget(s sim.GameEngine, pc *sim.CardState) bool {
	return pc.Card.Types().IsAttackAction() && pc.Card.Cost(s) <= 1
}

func (NimblismRed) Play(s sim.GameEngine, l sim.Logger, self *sim.CardState) {
	GrantNextCardBonusAttack(s, 3, nimblismIsTarget)
}

func (NimblismYellow) Play(s sim.GameEngine, l sim.Logger, self *sim.CardState) {
	GrantNextCardBonusAttack(s, 2, nimblismIsTarget)
}

func (NimblismBlue) Play(s sim.GameEngine, l sim.Logger, self *sim.CardState) {
	GrantNextCardBonusAttack(s, 1, nimblismIsTarget)
}
