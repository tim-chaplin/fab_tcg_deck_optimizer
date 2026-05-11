// Scout the Periphery — Generic Action. Cost 0. Printed pitch variants: Red 1, Yellow 2, Blue 3.
// Defense 2.
//
// Text: "Look at the top card of target hero's deck. The next attack action card you play from
// arsenal this turn gains +N{p}. **Go again**" (Red N=3, Yellow N=2, Blue N=1.)
//
// Modelling: deck-peek rider isn't modelled. The +N{p} grant targets an attack action card
// that itself was played from arsenal — scan CardsRemaining for the first attack action with
// CardState.FromArsenal set. Since the arsenal holds at most one card, the grant only fires
// when the arsenal-in card is an attack action queued later in the chain.

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

// scoutThePeripheryIsTarget gates the rider on attack action cards played from arsenal.
func scoutThePeripheryIsTarget(_ *sim.TurnState, pc *sim.CardState) bool {
	return pc.FromArsenal && pc.Card.Types().IsAttackAction()
}

func (ScoutThePeripheryRed) Play(s *sim.TurnState, l sim.Logger, self *sim.CardState) {
	GrantNextCardBonusAttack(s, 3, scoutThePeripheryIsTarget)
}

func (ScoutThePeripheryYellow) Play(s *sim.TurnState, l sim.Logger, self *sim.CardState) {
	GrantNextCardBonusAttack(s, 2, scoutThePeripheryIsTarget)
}

func (ScoutThePeripheryBlue) Play(s *sim.TurnState, l sim.Logger, self *sim.CardState) {
	GrantNextCardBonusAttack(s, 1, scoutThePeripheryIsTarget)
}
