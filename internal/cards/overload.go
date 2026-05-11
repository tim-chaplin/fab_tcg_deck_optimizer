// Overload — Generic Action - Attack. Cost 0. Printed power: Red 3, Yellow 2, Blue 1. Printed
// pitch variants: Red 1, Yellow 2, Blue 3. Defense 2.
//
// Text: "**Dominate** If Overload hits, it gains **go again**."

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
)

func overloadPlay(s card.GameEngine, l card.Logger, self *card.CardState) {
	if s.LikelyToHit(self) {
		self.GrantedGoAgain = true
	}
}

func (OverloadRed) Dominate() {}
func (OverloadRed) Play(s card.GameEngine, l card.Logger, self *card.CardState) {
	overloadPlay(s, l, self)
}

func (OverloadYellow) Dominate() {}
func (OverloadYellow) Play(s card.GameEngine, l card.Logger, self *card.CardState) {
	overloadPlay(s, l, self)
}

func (OverloadBlue) Dominate() {}
func (OverloadBlue) Play(s card.GameEngine, l card.Logger, self *card.CardState) {
	overloadPlay(s, l, self)
}
