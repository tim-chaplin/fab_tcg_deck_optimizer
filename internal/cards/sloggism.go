// Sloggism — Generic Action. Cost 3. Printed pitch variants: Red 1, Yellow 2, Blue 3. Defense 2.
//
// Text: "The next attack action card with cost 2 or greater you play this turn gains +N{p}. **Go
// again**" (Red N=6, Yellow N=5, Blue N=4.)

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

// sloggismIsTarget gates the rider on attack action cards whose cost is 2 or more.
func sloggismIsTarget(s *sim.TurnState, pc *sim.CardState) bool {
	return pc.Card.Types().IsAttackAction() && pc.Card.Cost(s) >= 2
}

func (SloggismRed) Play(s *sim.TurnState, l sim.Logger, self *sim.CardState) {
	GrantNextCardBonusAttack(s, 6, sloggismIsTarget)
}

func (SloggismYellow) Play(s *sim.TurnState, l sim.Logger, self *sim.CardState) {
	GrantNextCardBonusAttack(s, 5, sloggismIsTarget)
}

func (SloggismBlue) Play(s *sim.TurnState, l sim.Logger, self *sim.CardState) {
	GrantNextCardBonusAttack(s, 4, sloggismIsTarget)
}
