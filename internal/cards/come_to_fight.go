// Come to Fight — Generic Action. Cost 1. Printed pitch variants: Red 1, Yellow 2, Blue 3. Defense
// 3.
//
// Text: "The next attack action card you play this turn gains +N{p}. **Go again**" (Red N=3,
// Yellow N=2, Blue N=1.)

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

func (ComeToFightRed) Play(s *sim.TurnState, l sim.Logger, self *sim.CardState) {
	GrantNextCardBonusAttack(s, 3, IsAttackAction)
	n := self.DealEffectiveAttack(s)
	l.Log(self, n)
}

func (ComeToFightYellow) Play(s *sim.TurnState, l sim.Logger, self *sim.CardState) {
	GrantNextCardBonusAttack(s, 2, IsAttackAction)
	n := self.DealEffectiveAttack(s)
	l.Log(self, n)
}

func (ComeToFightBlue) Play(s *sim.TurnState, l sim.Logger, self *sim.CardState) {
	GrantNextCardBonusAttack(s, 1, IsAttackAction)
	n := self.DealEffectiveAttack(s)
	l.Log(self, n)
}
