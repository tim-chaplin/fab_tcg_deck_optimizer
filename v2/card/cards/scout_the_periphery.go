// Scout the Periphery — Generic Action. Cost 0. Printed pitch variants: Red 1, Yellow 2, Blue 3.
// Defense 2.
//
// Text: "Look at the top card of target hero's deck. The next attack action card you play from
// arsenal this turn gains +N{p}. **Go again**" (Red N=3, Yellow N=2, Blue N=1.)
//
// Modelling: deck-peek rider isn't modelled. The +N{p} grant only fires when the arsenal-in
// card is an attack action queued later in the chain.

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
)

// scoutThePeripheryIsTarget gates the rider on attack action cards played from arsenal.
func scoutThePeripheryIsTarget(_ card.GameEngine, pc *card.CardState) bool {
	return pc.FromArsenal && pc.Card.Types(nil).IsAttackAction()
}

func (ScoutThePeripheryRed) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	GrantNextCardBonusAttack(ge, 3, scoutThePeripheryIsTarget)
}

func (ScoutThePeripheryYellow) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	GrantNextCardBonusAttack(ge, 2, scoutThePeripheryIsTarget)
}

func (ScoutThePeripheryBlue) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	GrantNextCardBonusAttack(ge, 1, scoutThePeripheryIsTarget)
}
