// Captain's Call — Generic Action. Cost 0. Printed pitch variants: Red 1, Yellow 2, Blue 3.
// Defense 2. Go again.
//
// Text: "Choose 1; The next attack action card with cost N or less you play this turn gains
// +2{p}. The next attack action card with cost N or less you play this turn gains **go
// again**. **Go again**" (Red N=2, Yellow N=1, Blue N=0.)
//
// Mode 0 buffs the next cost-≤N attack action card +2{p}; mode 1 grants it go again.

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
)

// captainsCallPlay applies the modal grant to the next cost-≤maxCost attack action card in
// CardsRemaining. Fizzles silently if no follow-up attack action matches.
func captainsCallPlay(ge card.GameEngine, l card.Logger, self *card.CardState, maxCost int) {
	for _, pc := range ge.CardsRemaining() {
		if !pc.Card.Types(nil).IsAttackAction() {
			continue
		}
		if pc.Card.Cost(ge) > maxCost {
			continue
		}
		switch self.Mode {
		case 0:
			pc.BonusAttack += 2
		case 1:
			pc.GrantedGoAgain = true
		}
		break
	}
}

func (CaptainsCallRed) Modes() int { return 2 }
func (CaptainsCallRed) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	captainsCallPlay(ge, l, self, 2)
}

func (CaptainsCallYellow) Modes() int { return 2 }
func (CaptainsCallYellow) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	captainsCallPlay(ge, l, self, 1)
}

func (CaptainsCallBlue) Modes() int { return 2 }
func (CaptainsCallBlue) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	captainsCallPlay(ge, l, self, 0)
}
