// Spring Load — Generic Action - Attack. Cost 1. Printed power: Red 2, Yellow 2, Blue 2. Printed
// pitch variants: Red 1, Yellow 2, Blue 3. Defense 2.
//
// Text: "When this attacks, if you have no cards in hand, it gains +3{p}."

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
)

// springLoadPlay applies the +3{p} 'no cards in hand' rider.
func springLoadPlay(ge card.GameEngine, l card.Logger, self *card.CardState) {
	if len(ge.Hand()) == 0 {
		self.BonusAttack += 3
	}
}

func (SpringLoadRed) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	springLoadPlay(ge, l, self)
}

func (SpringLoadYellow) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	springLoadPlay(ge, l, self)
}

func (SpringLoadBlue) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	springLoadPlay(ge, l, self)
}
