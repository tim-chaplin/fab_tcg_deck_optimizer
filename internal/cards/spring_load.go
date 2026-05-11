// Spring Load — Generic Action - Attack. Cost 1. Printed power: Red 2, Yellow 2, Blue 2. Printed
// pitch variants: Red 1, Yellow 2, Blue 3. Defense 2.
//
// Text: "When this attacks, if you have no cards in hand, it gains +3{p}."

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
)

// springLoadPlay applies the +3{p} 'no cards in hand' rider, then emits the chain step.
func springLoadPlay(s card.GameEngine, l card.Logger, self *card.CardState) {
	if len(s.Hand()) == 0 {
		self.BonusAttack += 3
	}
}

func (SpringLoadRed) Play(s card.GameEngine, l card.Logger, self *card.CardState) {
	springLoadPlay(s, l, self)
}

func (SpringLoadYellow) Play(s card.GameEngine, l card.Logger, self *card.CardState) {
	springLoadPlay(s, l, self)
}

func (SpringLoadBlue) Play(s card.GameEngine, l card.Logger, self *card.CardState) {
	springLoadPlay(s, l, self)
}
