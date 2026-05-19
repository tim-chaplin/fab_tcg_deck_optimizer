// Push the Point — Generic Action - Attack. Cost 1. Printed power: Red 4, Yellow 3, Blue 2.
// Printed pitch variants: Red 1, Yellow 2, Blue 3. Defense 2.
//
// Text: "If the last attack on this combat chain hit, Push the Point gains +2{p}."

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

func pushThePointPlay(ge card.GameEngine, l card.Logger, self *card.CardState) {
	if ge.LastAttackHit() {
		self.BonusAttack += 2
	}
}

func (PushThePointRed) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	pushThePointPlay(ge, l, self)
}

func (PushThePointYellow) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	pushThePointPlay(ge, l, self)
}

func (PushThePointBlue) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	pushThePointPlay(ge, l, self)
}
